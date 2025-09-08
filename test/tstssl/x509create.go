package main

import (
	"bytes"
	"crypto"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"golang.org/x/crypto/cryptobyte"
	cryptobyte_asn1 "golang.org/x/crypto/cryptobyte/asn1"
	"io"
	"logutil"
	"math/big"
	"net"
	"net/url"
)

func oidFromECDHCurve(curve ecdh.Curve) (asn1.ObjectIdentifier, bool) {
	switch curve {
	case ecdh.X25519():
		return oidPublicKeyX25519, true
	case ecdh.P256():
		return oidNamedCurveP256, true
	case ecdh.P384():
		return oidNamedCurveP384, true
	case ecdh.P521():
		return oidNamedCurveP521, true
	}

	return nil, false
}

func oidFromNamedCurve(curve elliptic.Curve) (asn1.ObjectIdentifier, bool) {
	switch curve {
	case elliptic.P224():
		return oidNamedCurveP224, true
	case elliptic.P256():
		return oidNamedCurveP256, true
	case elliptic.P384():
		return oidNamedCurveP384, true
	case elliptic.P521():
		return oidNamedCurveP521, true
	}

	return nil, false
}

func marshalPublicKey(pub any) (publicKeyBytes []byte, publicKeyAlgorithm pkix.AlgorithmIdentifier, err error) {
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		publicKeyBytes, err = asn1.Marshal(pkcs1PublicKey{
			N: pub.N,
			E: pub.E,
		})
		if err != nil {
			return nil, pkix.AlgorithmIdentifier{}, err
		}
		publicKeyAlgorithm.Algorithm = oidPublicKeyRSA
		// This is a NULL parameters value which is required by
		// RFC 3279, Section 2.3.1.
		publicKeyAlgorithm.Parameters = asn1.NullRawValue
	case *ecdsa.PublicKey:
		oid, ok := oidFromNamedCurve(pub.Curve)
		if !ok {
			return nil, pkix.AlgorithmIdentifier{}, errors.New("x509: unsupported elliptic curve")
		}
		if !pub.Curve.IsOnCurve(pub.X, pub.Y) {
			return nil, pkix.AlgorithmIdentifier{}, errors.New("x509: invalid elliptic curve public key")
		}
		publicKeyBytes = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
		publicKeyAlgorithm.Algorithm = oidPublicKeyECDSA
		var paramBytes []byte
		paramBytes, err = asn1.Marshal(oid)
		if err != nil {
			return
		}
		publicKeyAlgorithm.Parameters.FullBytes = paramBytes
	case ed25519.PublicKey:
		publicKeyBytes = pub
		publicKeyAlgorithm.Algorithm = oidPublicKeyEd25519
	case *ecdh.PublicKey:
		publicKeyBytes = pub.Bytes()
		if pub.Curve() == ecdh.X25519() {
			publicKeyAlgorithm.Algorithm = oidPublicKeyX25519
		} else {
			oid, ok := oidFromECDHCurve(pub.Curve())
			if !ok {
				return nil, pkix.AlgorithmIdentifier{}, errors.New("x509: unsupported elliptic curve")
			}
			publicKeyAlgorithm.Algorithm = oidPublicKeyECDSA
			var paramBytes []byte
			paramBytes, err = asn1.Marshal(oid)
			if err != nil {
				return
			}
			publicKeyAlgorithm.Parameters.FullBytes = paramBytes
		}
	default:
		return nil, pkix.AlgorithmIdentifier{}, fmt.Errorf("x509: unsupported public key type: %T", pub)
	}

	return publicKeyBytes, publicKeyAlgorithm, nil
}

// signingParamsForKey returns the signature algorithm and its Algorithm
// Identifier to use for signing, based on the key type. If sigAlgo is not zero
// then it overrides the default.
func signingParamsForKey(key crypto.Signer, sigAlgo x509.SignatureAlgorithm) (x509.SignatureAlgorithm, pkix.AlgorithmIdentifier, error) {
	var ai pkix.AlgorithmIdentifier
	var pubType x509.PublicKeyAlgorithm
	var defaultAlgo x509.SignatureAlgorithm

	switch pub := key.Public().(type) {
	case *rsa.PublicKey:
		pubType = x509.RSA
		defaultAlgo = x509.SHA256WithRSA

	case *ecdsa.PublicKey:
		pubType = x509.ECDSA
		switch pub.Curve {
		case elliptic.P224(), elliptic.P256():
			defaultAlgo = x509.ECDSAWithSHA256
		case elliptic.P384():
			defaultAlgo = x509.ECDSAWithSHA384
		case elliptic.P521():
			defaultAlgo = x509.ECDSAWithSHA512
		default:
			return 0, ai, errors.New("x509: unsupported elliptic curve")
		}

	case ed25519.PublicKey:
		pubType = x509.Ed25519
		defaultAlgo = x509.PureEd25519

	default:
		return 0, ai, errors.New("x509: only RSA, ECDSA and Ed25519 keys supported")
	}

	if sigAlgo == 0 {
		sigAlgo = defaultAlgo
	}

	for _, details := range signatureAlgorithmDetails {
		if details.algo == sigAlgo {
			if details.pubKeyAlgo != pubType {
				return 0, ai, errors.New("x509: requested SignatureAlgorithm does not match private key type")
			}
			if details.hash == crypto.MD5 {
				return 0, ai, errors.New("x509: signing with MD5 is not supported")
			}

			return sigAlgo, pkix.AlgorithmIdentifier{
				Algorithm:  details.oid,
				Parameters: details.params,
			}, nil
		}
	}

	return 0, ai, errors.New("x509: unknown SignatureAlgorithm")
}

func subjectBytes(cert *x509.Certificate) ([]byte, error) {
	if len(cert.RawSubject) > 0 {
		return cert.RawSubject, nil
	}

	return asn1.Marshal(cert.Subject.ToRDNSequence())
}

// oidInExtensions reports whether an extension with the given oid exists in
// extensions.
func oidInExtensions(oid asn1.ObjectIdentifier, extensions []pkix.Extension) bool {
	for _, e := range extensions {
		if e.Id.Equal(oid) {
			return true
		}
	}
	return false
}

func reverseBitsInAByte(in byte) byte {
	b1 := in>>4 | in<<4
	b2 := b1>>2&0x33 | b1<<2&0xcc
	b3 := b2>>1&0x55 | b2<<1&0xaa
	return b3
}

func marshalKeyUsage(ku x509.KeyUsage) (pkix.Extension, error) {
	ext := pkix.Extension{Id: oidExtensionKeyUsage, Critical: true}

	var a [2]byte
	a[0] = reverseBitsInAByte(byte(ku))
	a[1] = reverseBitsInAByte(byte(ku >> 8))

	l := 1
	if a[1] != 0 {
		l = 2
	}

	bitString := a[:l]
	var err error
	ext.Value, err = asn1.Marshal(asn1.BitString{Bytes: bitString, BitLength: asn1BitLength(bitString)})
	return ext, err
}

// asn1BitLength returns the bit-length of bitString by considering the
// most-significant bit in a byte to be the "first" bit. This convention
// matches ASN.1, but differs from almost everything else.
func asn1BitLength(bitString []byte) int {
	bitLen := len(bitString) * 8

	for i := range bitString {
		b := bitString[len(bitString)-i-1]

		for bit := uint(0); bit < 8; bit++ {
			if (b>>bit)&1 == 1 {
				return bitLen
			}
			bitLen--
		}
	}

	return 0
}

func oidFromExtKeyUsage(eku x509.ExtKeyUsage) (oid asn1.ObjectIdentifier, ok bool) {
	for _, pair := range extKeyUsageOIDs {
		if eku == pair.extKeyUsage {
			return pair.oid, true
		}
	}
	return
}

func marshalExtKeyUsage(extUsages []x509.ExtKeyUsage, unknownUsages []asn1.ObjectIdentifier) (pkix.Extension, error) {
	ext := pkix.Extension{Id: oidExtensionExtendedKeyUsage}

	oids := make([]asn1.ObjectIdentifier, len(extUsages)+len(unknownUsages))
	for i, u := range extUsages {
		if oid, ok := oidFromExtKeyUsage(u); ok {
			oids[i] = oid
		} else {
			return ext, errors.New("x509: unknown extended key usage")
		}
	}

	copy(oids[len(extUsages):], unknownUsages)

	var err error
	ext.Value, err = asn1.Marshal(oids)
	return ext, err
}

func marshalBasicConstraints(isCA bool, maxPathLen int, maxPathLenZero bool) (pkix.Extension, error) {
	ext := pkix.Extension{Id: oidExtensionBasicConstraints, Critical: true}
	// Leaving MaxPathLen as zero indicates that no maximum path
	// length is desired, unless MaxPathLenZero is set. A value of
	// -1 causes encoding/asn1 to omit the value as desired.
	if maxPathLen == 0 && !maxPathLenZero {
		maxPathLen = -1
	}
	var err error
	ext.Value, err = asn1.Marshal(basicConstraints{isCA, maxPathLen})
	logutil.Debug("isCA %v maxPathLen %d maxPathLenZero %v ext.Value %v", isCA, maxPathLen, maxPathLenZero, ext.Value)
	return ext, err
}

// marshalSANs marshals a list of addresses into a the contents of an X.509
// SubjectAlternativeName extension.
func marshalSANs(dnsNames, emailAddresses []string, ipAddresses []net.IP, uris []*url.URL) (derBytes []byte, err error) {
	var rawValues []asn1.RawValue
	for _, name := range dnsNames {
		if err := isIA5String(name); err != nil {
			return nil, err
		}
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeDNS, Class: 2, Bytes: []byte(name)})
	}
	for _, email := range emailAddresses {
		if err := isIA5String(email); err != nil {
			return nil, err
		}
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeEmail, Class: 2, Bytes: []byte(email)})
	}
	for _, rawIP := range ipAddresses {
		// If possible, we always want to encode IPv4 addresses in 4 bytes.
		ip := rawIP.To4()
		if ip == nil {
			ip = rawIP
		}
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeIP, Class: 2, Bytes: ip})
	}
	for _, uri := range uris {
		uriStr := uri.String()
		if err := isIA5String(uriStr); err != nil {
			return nil, err
		}
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeURI, Class: 2, Bytes: []byte(uriStr)})
	}
	return asn1.Marshal(rawValues)
}

func marshalCertificatePolicies(policies []x509.OID, policyIdentifiers []asn1.ObjectIdentifier) (pkix.Extension, error) {
	ext := pkix.Extension{Id: oidExtensionCertificatePolicies}

	b := cryptobyte.NewBuilder(make([]byte, 0, 128))
	b.AddASN1(cryptobyte_asn1.SEQUENCE, func(child *cryptobyte.Builder) {

		for _, v := range policies {
			child.AddASN1(cryptobyte_asn1.SEQUENCE, func(child *cryptobyte.Builder) {
				child.AddASN1(cryptobyte_asn1.OBJECT_IDENTIFIER, func(child *cryptobyte.Builder) {
					var vb []byte
					vb, _ = v.MarshalBinary()
					if len(vb) == 0 {
						child.SetError(errors.New("invalid policy object identifier"))
						return
					}
					child.AddBytes(vb)
				})
			})
		}

		for _, v := range policyIdentifiers {
			child.AddASN1(cryptobyte_asn1.SEQUENCE, func(child *cryptobyte.Builder) {
				child.AddASN1ObjectIdentifier(v)
			})
		}
	})

	var err error
	ext.Value, err = b.Bytes()
	return ext, err
}

func buildCertExtensions(template *x509.Certificate, subjectIsEmpty bool, authorityKeyId []byte, subjectKeyId []byte) (ret []pkix.Extension, err error) {
	ret = make([]pkix.Extension, 10 /* maximum number of elements. */)
	n := 0

	if template.KeyUsage != 0 &&
		!oidInExtensions(oidExtensionKeyUsage, template.ExtraExtensions) {
		ret[n], err = marshalKeyUsage(template.KeyUsage)
		if err != nil {
			return nil, err
		}
		n++
	}

	if (len(template.ExtKeyUsage) > 0 || len(template.UnknownExtKeyUsage) > 0) &&
		!oidInExtensions(oidExtensionExtendedKeyUsage, template.ExtraExtensions) {
		ret[n], err = marshalExtKeyUsage(template.ExtKeyUsage, template.UnknownExtKeyUsage)
		if err != nil {
			return nil, err
		}
		n++
	}

	if template.BasicConstraintsValid && !oidInExtensions(oidExtensionBasicConstraints, template.ExtraExtensions) {
		ret[n], err = marshalBasicConstraints(template.IsCA, template.MaxPathLen, template.MaxPathLenZero)
		if err != nil {
			return nil, err
		}
		n++
	}

	if len(subjectKeyId) > 0 && !oidInExtensions(oidExtensionSubjectKeyId, template.ExtraExtensions) {
		ret[n].Id = oidExtensionSubjectKeyId
		ret[n].Value, err = asn1.Marshal(subjectKeyId)
		if err != nil {
			return
		}
		n++
	}

	if len(authorityKeyId) > 0 && !oidInExtensions(oidExtensionAuthorityKeyId, template.ExtraExtensions) {
		ret[n].Id = oidExtensionAuthorityKeyId
		ret[n].Value, err = asn1.Marshal(authKeyId{authorityKeyId})
		if err != nil {
			return
		}
		n++
	}

	if (len(template.OCSPServer) > 0 || len(template.IssuingCertificateURL) > 0) &&
		!oidInExtensions(oidExtensionAuthorityInfoAccess, template.ExtraExtensions) {
		ret[n].Id = oidExtensionAuthorityInfoAccess
		var aiaValues []authorityInfoAccess
		for _, name := range template.OCSPServer {
			aiaValues = append(aiaValues, authorityInfoAccess{
				Method:   oidAuthorityInfoAccessOcsp,
				Location: asn1.RawValue{Tag: 6, Class: 2, Bytes: []byte(name)},
			})
		}
		for _, name := range template.IssuingCertificateURL {
			aiaValues = append(aiaValues, authorityInfoAccess{
				Method:   oidAuthorityInfoAccessIssuers,
				Location: asn1.RawValue{Tag: 6, Class: 2, Bytes: []byte(name)},
			})
		}
		ret[n].Value, err = asn1.Marshal(aiaValues)
		if err != nil {
			return
		}
		n++
	}

	if (len(template.DNSNames) > 0 || len(template.EmailAddresses) > 0 || len(template.IPAddresses) > 0 || len(template.URIs) > 0) &&
		!oidInExtensions(oidExtensionSubjectAltName, template.ExtraExtensions) {
		ret[n].Id = oidExtensionSubjectAltName
		// From RFC 5280, Section 4.2.1.6:
		// “If the subject field contains an empty sequence ... then
		// subjectAltName extension ... is marked as critical”
		ret[n].Critical = subjectIsEmpty
		ret[n].Value, err = marshalSANs(template.DNSNames, template.EmailAddresses, template.IPAddresses, template.URIs)
		if err != nil {
			return
		}
		n++
	}

	//usePolicies := false
	if ((len(template.PolicyIdentifiers) > 0) || (len(template.Policies) > 0)) &&
		!oidInExtensions(oidExtensionCertificatePolicies, template.ExtraExtensions) {
		ret[n], err = marshalCertificatePolicies(template.Policies, template.PolicyIdentifiers)
		if err != nil {
			return nil, err
		}
		n++
	}

	if (len(template.PermittedDNSDomains) > 0 || len(template.ExcludedDNSDomains) > 0 ||
		len(template.PermittedIPRanges) > 0 || len(template.ExcludedIPRanges) > 0 ||
		len(template.PermittedEmailAddresses) > 0 || len(template.ExcludedEmailAddresses) > 0 ||
		len(template.PermittedURIDomains) > 0 || len(template.ExcludedURIDomains) > 0) &&
		!oidInExtensions(oidExtensionNameConstraints, template.ExtraExtensions) {
		ret[n].Id = oidExtensionNameConstraints
		ret[n].Critical = template.PermittedDNSDomainsCritical

		ipAndMask := func(ipNet *net.IPNet) []byte {
			maskedIP := ipNet.IP.Mask(ipNet.Mask)
			ipAndMask := make([]byte, 0, len(maskedIP)+len(ipNet.Mask))
			ipAndMask = append(ipAndMask, maskedIP...)
			ipAndMask = append(ipAndMask, ipNet.Mask...)
			return ipAndMask
		}

		serialiseConstraints := func(dns []string, ips []*net.IPNet, emails []string, uriDomains []string) (der []byte, err error) {
			var b cryptobyte.Builder

			for _, name := range dns {
				if err = isIA5String(name); err != nil {
					return nil, err
				}

				b.AddASN1(cryptobyte_asn1.SEQUENCE, func(b *cryptobyte.Builder) {
					b.AddASN1(cryptobyte_asn1.Tag(2).ContextSpecific(), func(b *cryptobyte.Builder) {
						b.AddBytes([]byte(name))
					})
				})
			}

			for _, ipNet := range ips {
				b.AddASN1(cryptobyte_asn1.SEQUENCE, func(b *cryptobyte.Builder) {
					b.AddASN1(cryptobyte_asn1.Tag(7).ContextSpecific(), func(b *cryptobyte.Builder) {
						b.AddBytes(ipAndMask(ipNet))
					})
				})
			}

			for _, email := range emails {
				if err = isIA5String(email); err != nil {
					return nil, err
				}

				b.AddASN1(cryptobyte_asn1.SEQUENCE, func(b *cryptobyte.Builder) {
					b.AddASN1(cryptobyte_asn1.Tag(1).ContextSpecific(), func(b *cryptobyte.Builder) {
						b.AddBytes([]byte(email))
					})
				})
			}

			for _, uriDomain := range uriDomains {
				if err = isIA5String(uriDomain); err != nil {
					return nil, err
				}

				b.AddASN1(cryptobyte_asn1.SEQUENCE, func(b *cryptobyte.Builder) {
					b.AddASN1(cryptobyte_asn1.Tag(6).ContextSpecific(), func(b *cryptobyte.Builder) {
						b.AddBytes([]byte(uriDomain))
					})
				})
			}

			return b.Bytes()
		}

		permitted, err := serialiseConstraints(template.PermittedDNSDomains, template.PermittedIPRanges, template.PermittedEmailAddresses, template.PermittedURIDomains)
		if err != nil {
			return nil, err
		}

		excluded, err := serialiseConstraints(template.ExcludedDNSDomains, template.ExcludedIPRanges, template.ExcludedEmailAddresses, template.ExcludedURIDomains)
		if err != nil {
			return nil, err
		}

		var b cryptobyte.Builder
		b.AddASN1(cryptobyte_asn1.SEQUENCE, func(b *cryptobyte.Builder) {
			if len(permitted) > 0 {
				b.AddASN1(cryptobyte_asn1.Tag(0).ContextSpecific().Constructed(), func(b *cryptobyte.Builder) {
					b.AddBytes(permitted)
				})
			}

			if len(excluded) > 0 {
				b.AddASN1(cryptobyte_asn1.Tag(1).ContextSpecific().Constructed(), func(b *cryptobyte.Builder) {
					b.AddBytes(excluded)
				})
			}
		})

		ret[n].Value, err = b.Bytes()
		if err != nil {
			return nil, err
		}
		n++
	}

	if len(template.CRLDistributionPoints) > 0 &&
		!oidInExtensions(oidExtensionCRLDistributionPoints, template.ExtraExtensions) {
		ret[n].Id = oidExtensionCRLDistributionPoints

		var crlDp []distributionPoint
		for _, name := range template.CRLDistributionPoints {
			dp := distributionPoint{
				DistributionPoint: distributionPointName{
					FullName: []asn1.RawValue{
						{Tag: 6, Class: 2, Bytes: []byte(name)},
					},
				},
			}
			crlDp = append(crlDp, dp)
		}

		ret[n].Value, err = asn1.Marshal(crlDp)
		if err != nil {
			return
		}
		n++
	}

	// Adding another extension here? Remember to update the maximum number
	// of elements in the make() at the top of the function and the list of
	// template fields used in CreateCertificate documentation.

	return append(ret[:n], template.ExtraExtensions...), nil
}

func hash_func(algo x509.SignatureAlgorithm) crypto.Hash {
	for _, details := range signatureAlgorithmDetails {
		if details.algo == algo {
			return details.hash
		}
	}
	return crypto.Hash(0)
}

func is_RSAPSS(algo x509.SignatureAlgorithm) bool {
	for _, details := range signatureAlgorithmDetails {
		if details.algo == algo {
			return details.isRSAPSS
		}
	}
	return false
}

func signaturePublicKeyAlgoMismatchError(expectedPubKeyAlgo x509.PublicKeyAlgorithm, pubKey any) error {
	return fmt.Errorf("x509: signature algorithm specifies an %s public key, but have public key of type %T", expectedPubKeyAlgo.String(), pubKey)
}

// checkSignature verifies that signature is a valid signature over signed from
// a crypto.PublicKey.
func checkSignature(algo x509.SignatureAlgorithm, signed, signature []byte, publicKey crypto.PublicKey, allowSHA1 bool) (err error) {
	var hashType crypto.Hash
	var pubKeyAlgo x509.PublicKeyAlgorithm

	for _, details := range signatureAlgorithmDetails {
		if details.algo == algo {
			hashType = details.hash
			pubKeyAlgo = details.pubKeyAlgo
			break
		}
	}

	switch hashType {
	case crypto.Hash(0):
		if pubKeyAlgo != x509.Ed25519 {
			logutil.Error(" ")
			return x509.ErrUnsupportedAlgorithm
		}
	case crypto.MD5:
		return x509.InsecureAlgorithmError(algo)
	case crypto.SHA1:
		// SHA-1 signatures are only allowed for CRLs and CSRs.
		if !allowSHA1 {
			logutil.Error(" ")
			return x509.InsecureAlgorithmError(algo)
		}
		fallthrough
	default:
		if !hashType.Available() {
			logutil.Error(" ")
			return x509.ErrUnsupportedAlgorithm
		}
		h := hashType.New()
		h.Write(signed)
		signed = h.Sum(nil)
	}

	switch pub := publicKey.(type) {
	case *rsa.PublicKey:
		if pubKeyAlgo != x509.RSA {
			logutil.Error(" ")
			return signaturePublicKeyAlgoMismatchError(pubKeyAlgo, pub)
		}
		if is_RSAPSS(algo) {
			return rsa.VerifyPSS(pub, hashType, signed, signature, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
		} else {
			return rsa.VerifyPKCS1v15(pub, hashType, signed, signature)
		}
	case *ecdsa.PublicKey:
		if pubKeyAlgo != x509.ECDSA {
			return signaturePublicKeyAlgoMismatchError(pubKeyAlgo, pub)
		}
		if !ecdsa.VerifyASN1(pub, signed, signature) {
			return errors.New("x509: ECDSA verification failure")
		}
		return
	case ed25519.PublicKey:
		if pubKeyAlgo != x509.Ed25519 {
			return signaturePublicKeyAlgoMismatchError(pubKeyAlgo, pub)
		}
		if !ed25519.Verify(pub, signed, signature) {
			return errors.New("x509: Ed25519 verification failure")
		}
		return
	}
	logutil.Error(" ")
	return x509.ErrUnsupportedAlgorithm
}

func signTBS(tbs []byte, key crypto.Signer, sigAlg x509.SignatureAlgorithm, rand io.Reader) ([]byte, error) {
	signed := tbs
	hashFunc := hash_func(sigAlg)
	if hashFunc != 0 {
		h := hashFunc.New()
		h.Write(signed)
		signed = h.Sum(nil)
	}

	var signerOpts crypto.SignerOpts = hashFunc
	if is_RSAPSS(sigAlg) {
		signerOpts = &rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthEqualsHash,
			Hash:       hashFunc,
		}
	}

	signature, err := key.Sign(rand, signed, signerOpts)
	if err != nil {
		return nil, err
	}

	// Check the signature to ensure the crypto.Signer behaved correctly.
	if err := checkSignature(sigAlg, tbs, signature, key.Public(), true); err != nil {
		return nil, fmt.Errorf("x509: signature returned by signer is invalid: %w", err)
	}

	return signature, nil
}

func createCertificate(rand io.Reader, template, parent *x509.Certificate, pub, priv any) ([]byte, error) {

	key, ok := priv.(crypto.Signer)
	if !ok {
		return nil, errors.New("x509: certificate private key does not implement crypto.Signer")
	}

	serialNumber := template.SerialNumber
	if serialNumber == nil {
		// Generate a serial number following RFC 5280 Section 4.1.2.2 if one is not provided.
		// Requirements:
		//   - serial number must be positive
		//   - at most 20 octets when encoded
		maxSerial := big.NewInt(1).Lsh(big.NewInt(1), 20*8)
		for {
			var err error
			serialNumber, err = cryptorand.Int(rand, maxSerial)
			if err != nil {
				return nil, err
			}
			// If the serial is exactly 20 octets, check if the high bit of the first byte is set.
			// If so, generate a new serial, since it will be padded with a leading 0 byte during
			// encoding so that the serial is not interpreted as a negative integer, making it
			// 21 octets.
			if serialBytes := serialNumber.Bytes(); len(serialBytes) > 0 && (len(serialBytes) < 20 || serialBytes[0]&0x80 == 0) {
				break
			}
		}
	}

	// RFC 5280 Section 4.1.2.2: serial number must be positive
	//
	// We _should_ also restrict serials to <= 20 octets, but it turns out a lot of people
	// get this wrong, in part because the encoding can itself alter the length of the
	// serial. For now we accept these non-conformant serials.
	if serialNumber.Sign() == -1 {
		return nil, errors.New("x509: serial number must be positive")
	}

	if template.BasicConstraintsValid && !template.IsCA && template.MaxPathLen != -1 && (template.MaxPathLen != 0 || template.MaxPathLenZero) {
		return nil, errors.New("x509: only CAs are allowed to specify MaxPathLen")
	}

	signatureAlgorithm, algorithmIdentifier, err := signingParamsForKey(key, template.SignatureAlgorithm)
	if err != nil {
		return nil, err
	}

	publicKeyBytes, publicKeyAlgorithm, err := marshalPublicKey(pub)
	if err != nil {
		return nil, err
	}
	if getPublicKeyAlgorithmFromOID(publicKeyAlgorithm.Algorithm) == x509.UnknownPublicKeyAlgorithm {
		return nil, fmt.Errorf("x509: unsupported public key type: %T", pub)
	}

	asn1Issuer, err := subjectBytes(parent)
	if err != nil {
		return nil, err
	}

	asn1Subject, err := subjectBytes(template)
	if err != nil {
		return nil, err
	}

	authorityKeyId := template.AuthorityKeyId
	if !bytes.Equal(asn1Issuer, asn1Subject) && len(parent.SubjectKeyId) > 0 {
		authorityKeyId = parent.SubjectKeyId
	}

	subjectKeyId := template.SubjectKeyId
	if len(subjectKeyId) == 0 && template.IsCA {
		// SubjectKeyId generated using method 1 in RFC 5280, Section 4.2.1.2:
		//   (1) The keyIdentifier is composed of the 160-bit SHA-1 hash of the
		//   value of the BIT STRING subjectPublicKey (excluding the tag,
		//   length, and number of unused bits).
		h := sha1.Sum(publicKeyBytes)
		subjectKeyId = h[:]
		logutil.DebugBuffer(subjectKeyId, "SubjectKeyId")
		logutil.DebugBuffer(publicKeyBytes, "publicKeyBytes")
	}

	// Check that the signer's public key matches the private key, if available.
	type privateKey interface {
		Equal(crypto.PublicKey) bool
	}
	if privPub, ok := key.Public().(privateKey); !ok {
		return nil, errors.New("x509: internal error: supported public key does not implement Equal")
	} else if parent.PublicKey != nil && !privPub.Equal(parent.PublicKey) {
		return nil, errors.New("x509: provided PrivateKey doesn't match parent's PublicKey")
	}

	extensions, err := buildCertExtensions(template, bytes.Equal(asn1Subject, emptyASN1Subject), authorityKeyId, subjectKeyId)
	if err != nil {
		return nil, err
	}

	encodedPublicKey := asn1.BitString{BitLength: len(publicKeyBytes) * 8, Bytes: publicKeyBytes}
	c := tbsCertificate{
		Version:            2,
		SerialNumber:       serialNumber,
		SignatureAlgorithm: algorithmIdentifier,
		Issuer:             asn1.RawValue{FullBytes: asn1Issuer},
		Validity:           validity{template.NotBefore.UTC(), template.NotAfter.UTC()},
		Subject:            asn1.RawValue{FullBytes: asn1Subject},
		PublicKey:          publicKeyInfo{nil, publicKeyAlgorithm, encodedPublicKey},
		Extensions:         extensions,
	}

	tbsCertContents, err := asn1.Marshal(c)
	if err != nil {
		return nil, err
	}
	c.Raw = tbsCertContents

	signature, err := signTBS(tbsCertContents, key, signatureAlgorithm, rand)
	if err != nil {
		return nil, err
	}

	return asn1.Marshal(certificate{
		TBSCertificate:     c,
		SignatureAlgorithm: algorithmIdentifier,
		SignatureValue:     asn1.BitString{Bytes: signature, BitLength: len(signature) * 8},
	})

}
