package main

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"math/big"
	"time"
)

var (
	// RFC 3279, 2.3 Public Key Algorithms
	//
	//	pkcs-1 OBJECT IDENTIFIER ::== { iso(1) member-body(2) us(840)
	//		rsadsi(113549) pkcs(1) 1 }
	//
	// rsaEncryption OBJECT IDENTIFIER ::== { pkcs1-1 1 }
	//
	//	id-dsa OBJECT IDENTIFIER ::== { iso(1) member-body(2) us(840)
	//		x9-57(10040) x9cm(4) 1 }
	oidPublicKeyRSA = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	oidPublicKeyDSA = asn1.ObjectIdentifier{1, 2, 840, 10040, 4, 1}
	// RFC 5480, 2.1.1 Unrestricted Algorithm Identifier and Parameters
	//
	//	id-ecPublicKey OBJECT IDENTIFIER ::= {
	//		iso(1) member-body(2) us(840) ansi-X9-62(10045) keyType(2) 1 }
	oidPublicKeyECDSA = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
	// RFC 8410, Section 3
	//
	//	id-X25519    OBJECT IDENTIFIER ::= { 1 3 101 110 }
	//	id-Ed25519   OBJECT IDENTIFIER ::= { 1 3 101 112 }
	oidPublicKeyX25519  = asn1.ObjectIdentifier{1, 3, 101, 110}
	oidPublicKeyEd25519 = asn1.ObjectIdentifier{1, 3, 101, 112}
)

// RFC 5480, 2.1.1.1. Named Curve
//
//	secp224r1 OBJECT IDENTIFIER ::= {
//	  iso(1) identified-organization(3) certicom(132) curve(0) 33 }
//
//	secp256r1 OBJECT IDENTIFIER ::= {
//	  iso(1) member-body(2) us(840) ansi-X9-62(10045) curves(3)
//	  prime(1) 7 }
//
//	secp384r1 OBJECT IDENTIFIER ::= {
//	  iso(1) identified-organization(3) certicom(132) curve(0) 34 }
//
//	secp521r1 OBJECT IDENTIFIER ::= {
//	  iso(1) identified-organization(3) certicom(132) curve(0) 35 }
//
// NB: secp256r1 is equivalent to prime256v1
var (
	oidNamedCurveP224 = asn1.ObjectIdentifier{1, 3, 132, 0, 33}
	oidNamedCurveP256 = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}
	oidNamedCurveP384 = asn1.ObjectIdentifier{1, 3, 132, 0, 34}
	oidNamedCurveP521 = asn1.ObjectIdentifier{1, 3, 132, 0, 35}
)

var emptyRawValue = asn1.RawValue{}

var signatureAlgorithmDetails = []struct {
	algo       x509.SignatureAlgorithm
	name       string
	oid        asn1.ObjectIdentifier
	params     asn1.RawValue
	pubKeyAlgo x509.PublicKeyAlgorithm
	hash       crypto.Hash
	isRSAPSS   bool
}{
	{x509.MD5WithRSA, "MD5-RSA", oidSignatureMD5WithRSA, asn1.NullRawValue, x509.RSA, crypto.MD5, false},
	{x509.SHA1WithRSA, "SHA1-RSA", oidSignatureSHA1WithRSA, asn1.NullRawValue, x509.RSA, crypto.SHA1, false},
	{x509.SHA1WithRSA, "SHA1-RSA", oidISOSignatureSHA1WithRSA, asn1.NullRawValue, x509.RSA, crypto.SHA1, false},
	{x509.SHA256WithRSA, "SHA256-RSA", oidSignatureSHA256WithRSA, asn1.NullRawValue, x509.RSA, crypto.SHA256, false},
	{x509.SHA384WithRSA, "SHA384-RSA", oidSignatureSHA384WithRSA, asn1.NullRawValue, x509.RSA, crypto.SHA384, false},
	{x509.SHA512WithRSA, "SHA512-RSA", oidSignatureSHA512WithRSA, asn1.NullRawValue, x509.RSA, crypto.SHA512, false},
	{x509.SHA256WithRSAPSS, "SHA256-RSAPSS", oidSignatureRSAPSS, pssParametersSHA256, x509.RSA, crypto.SHA256, true},
	{x509.SHA384WithRSAPSS, "SHA384-RSAPSS", oidSignatureRSAPSS, pssParametersSHA384, x509.RSA, crypto.SHA384, true},
	{x509.SHA512WithRSAPSS, "SHA512-RSAPSS", oidSignatureRSAPSS, pssParametersSHA512, x509.RSA, crypto.SHA512, true},
	{x509.DSAWithSHA1, "DSA-SHA1", oidSignatureDSAWithSHA1, emptyRawValue, x509.DSA, crypto.SHA1, false},
	{x509.DSAWithSHA256, "DSA-SHA256", oidSignatureDSAWithSHA256, emptyRawValue, x509.DSA, crypto.SHA256, false},
	{x509.ECDSAWithSHA1, "ECDSA-SHA1", oidSignatureECDSAWithSHA1, emptyRawValue, x509.ECDSA, crypto.SHA1, false},
	{x509.ECDSAWithSHA256, "ECDSA-SHA256", oidSignatureECDSAWithSHA256, emptyRawValue, x509.ECDSA, crypto.SHA256, false},
	{x509.ECDSAWithSHA384, "ECDSA-SHA384", oidSignatureECDSAWithSHA384, emptyRawValue, x509.ECDSA, crypto.SHA384, false},
	{x509.ECDSAWithSHA512, "ECDSA-SHA512", oidSignatureECDSAWithSHA512, emptyRawValue, x509.ECDSA, crypto.SHA512, false},
	{x509.PureEd25519, "Ed25519", oidSignatureEd25519, emptyRawValue, x509.Ed25519, crypto.Hash(0) /* no pre-hashing */, false},
}

// extKeyUsageOIDs contains the mapping between an ExtKeyUsage and its OID.
var extKeyUsageOIDs = []struct {
	extKeyUsage x509.ExtKeyUsage
	oid         asn1.ObjectIdentifier
}{
	{x509.ExtKeyUsageAny, oidExtKeyUsageAny},
	{x509.ExtKeyUsageServerAuth, oidExtKeyUsageServerAuth},
	{x509.ExtKeyUsageClientAuth, oidExtKeyUsageClientAuth},
	{x509.ExtKeyUsageCodeSigning, oidExtKeyUsageCodeSigning},
	{x509.ExtKeyUsageEmailProtection, oidExtKeyUsageEmailProtection},
	{x509.ExtKeyUsageIPSECEndSystem, oidExtKeyUsageIPSECEndSystem},
	{x509.ExtKeyUsageIPSECTunnel, oidExtKeyUsageIPSECTunnel},
	{x509.ExtKeyUsageIPSECUser, oidExtKeyUsageIPSECUser},
	{x509.ExtKeyUsageTimeStamping, oidExtKeyUsageTimeStamping},
	{x509.ExtKeyUsageOCSPSigning, oidExtKeyUsageOCSPSigning},
	{x509.ExtKeyUsageMicrosoftServerGatedCrypto, oidExtKeyUsageMicrosoftServerGatedCrypto},
	{x509.ExtKeyUsageNetscapeServerGatedCrypto, oidExtKeyUsageNetscapeServerGatedCrypto},
	{x509.ExtKeyUsageMicrosoftCommercialCodeSigning, oidExtKeyUsageMicrosoftCommercialCodeSigning},
	{x509.ExtKeyUsageMicrosoftKernelCodeSigning, oidExtKeyUsageMicrosoftKernelCodeSigning},
}

// emptyASN1Subject is the ASN.1 DER encoding of an empty Subject, which is
// just an empty SEQUENCE.
var emptyASN1Subject = []byte{0x30, 0}

type basicConstraints struct {
	IsCA       bool `asn1:"optional"`
	MaxPathLen int  `asn1:"optional,default:-1"`
}

// RFC 5280,  4.2.1.1
type authKeyId struct {
	Id []byte `asn1:"optional,tag:0"`
}

// RFC 5280, 4.2.2.1
type authorityInfoAccess struct {
	Method   asn1.ObjectIdentifier
	Location asn1.RawValue
}

type distributionPointName struct {
	FullName     []asn1.RawValue  `asn1:"optional,tag:0"`
	RelativeName pkix.RDNSequence `asn1:"optional,tag:1"`
}

// RFC 5280, 4.2.1.14
type distributionPoint struct {
	DistributionPoint distributionPointName `asn1:"optional,tag:0"`
	Reason            asn1.BitString        `asn1:"optional,tag:1"`
	CRLIssuer         asn1.RawValue         `asn1:"optional,tag:2"`
}

type validity struct {
	NotBefore, NotAfter time.Time
}

// These structures reflect the ASN.1 structure of X.509 certificates.:

type certificate struct {
	TBSCertificate     tbsCertificate
	SignatureAlgorithm pkix.AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

type tbsCertificate struct {
	Raw                asn1.RawContent
	Version            int `asn1:"optional,explicit,default:0,tag:0"`
	SerialNumber       *big.Int
	SignatureAlgorithm pkix.AlgorithmIdentifier
	Issuer             asn1.RawValue
	Validity           validity
	Subject            asn1.RawValue
	PublicKey          publicKeyInfo
	UniqueId           asn1.BitString   `asn1:"optional,tag:1"`
	SubjectUniqueId    asn1.BitString   `asn1:"optional,tag:2"`
	Extensions         []pkix.Extension `asn1:"omitempty,optional,explicit,tag:3"`
}

const (
	KEYWORD_COUNTRY          = "country"
	KEYWORD_PROVINCE         = "province"
	KEYWORD_LOCALITY         = "locality"
	KEYWORD_STREETADDRESS    = "streetaddress"
	KEYWORD_POSTALCODE       = "postalcode"
	KEYWORD_ORGANIZATION     = "organization"
	KEYWORD_ORGANIZATIONUNIT = "organizationalunit"
	KEYWORD_COMMONNAME       = "commonname"
	KEYWORD_SERIALNUMBER     = "serialnumber"
	KEYWROD_EXTRANAMES       = "extranames"
	KEYWORD_TYPE             = "type"
	KEYWORD_VALUE            = "value"
	KEYWORD_TAG              = "tag"
	KEYWORD_CONTENT          = "content"
	KEYWORD_ID               = "id"
	KEYWORD_CRITICAL         = "critical"
	KEYWORD_VERSION          = "version"
	KEYWORD_ATTRIBUTES       = "attributes"
	KEYWORD_EXTENSIONS       = "extensions"
	KEYWORD_EXTRA_EXTENSIONS = "extraextensions"
)

const (
	KEYWORD_EXT_KEYUSAGE_ANY                               = "any"
	KEYWORD_EXT_KEYUSAGE_SERVER_AUTH                       = "serverauth"
	KEYWORD_EXT_KEYUSAGE_CLIENT_AUTH                       = "clientauth"
	KEYWORD_EXT_KEYUSAGE_CODE_SIGNING                      = "codesigning"
	KEYWORD_EXT_KEYUSAGE_EMAIL_PROTECTION                  = "emailprotection"
	KEYWORD_EXT_KEYUSAGE_IPSEC_END_SYSTEM                  = "ipsecencsystem"
	KEYWORD_EXT_KEYUSAGE_IPSEC_TUNNEL                      = "ipsectunnel"
	KEYWORD_EXT_KEYUSAGE_IPSEC_USER                        = "ipsecuser"
	KEYWORD_EXT_KEYUSAGE_TIME_STAMPING                     = "timestamping"
	KEYWORD_EXT_KEYUSAGE_OCSP_SIGNING                      = "ocspsigning"
	KEYWORD_EXT_KEYUSAGE_MICROSOFT_SERVER_GATED_CRYPTO     = "microsoftservergatedcrypto"
	KEYWORD_EXT_KEYUSAGE_NETSCAPE_SERVER_GATED_CRYPTO      = "netscapeserergatedcrypto"
	KEYWORD_EXT_KEYUSAGE_MICROSOFT_COMMERCIAL_CODE_SIGNING = "microsoftcommercialcodesigning"
	KEYWORD_EXT_KEYUSAGE_MICROSOFT_KERNEL_CODE_SIGNING     = "microsoftkernelcodesigning"
)

const (
	KEYWORD_MD2_WITH_RSA       = "md2withrsa"
	KEYWORD_MD5_WITH_RSA       = "md5withrsa"
	KEYWORD_SHA1_WITH_RSA      = "sha1withrsa"
	KEYWORD_SHA256_WITH_RSA    = "sha256withrsa"
	KEYWORD_SHA384_WITH_RSA    = "sha384withrsa"
	KEYWORD_SHA512_WITH_RSA    = "sha512withrsa"
	KEYWORD_DSA_WITH_SHA1      = "dsawithsha1"
	KEYWORD_DSA_WITH_SHA256    = "dsawithsha256"
	KEYWORD_ECDSA_WITH_SHA1    = "ecdsawithsha1"
	KEYWORD_ECDSA_WITH_SHA256  = "ecdsawithsha256"
	KEYWORD_ECDSA_WITH_SHA384  = "ecdsawithsha384"
	KEYWORD_ECDSA_WITH_SHA512  = "ecdsawithsha512"
	KEYWORD_SHA256_WITH_RSAPSS = "sha256withrsapss"
	KEYWORD_SHA384_WITH_RSAPSS = "sha384withrsapss"
	KEYWORD_SHA512_WITH_RSAPSS = "sha512withrsapss"
	KEYWORD_PURE_ED25519       = "puered25519"
)

const (
	KEYWORD_DIGITAL_SIGNATURE  = "digitalsignature"
	KEYWORD_CONTENT_COMMITMENT = "contentcommitment"
	KEYWORD_KEY_ENCIPHERMENT   = "keyencipherment"
	KEYWORD_DATA_ENCIPHERMENT  = "dataencipherment"
	KEYWORD_KEY_AGREEMENT      = "keyagreement"
	KEYWORD_CERT_SIGN          = "certsign"
	KEYWORD_CRL_SIGN           = "crlsign"
	KEYWORD_ENCIPHER_ONLY      = "encipheronly"
	KEYWORD_DECIPHER_ONLY      = "decipheronly"
)

const (
	KEYWORD_EXT_KEY_USAGE                  = "extkeyusage"
	KEYWORD_AUTHORITY_KEY_ID               = "authoritykeyid"
	KEYWORD_BASIC_CONSTRAINTS_VALID        = "basicconstraintsvalid"
	KEYWORD_CRL_DISTRIBUTION_POINTS        = "crldistributionpoints"
	KEYWORD_DNS_NAMES                      = "dnsnames"
	KEYWORD_EMAIL_ADDRESSES                = "emailaddresses"
	KEYWORD_EXCLUDED_DNS_DOMAINS           = "excludeddnsdomains"
	KEYWORD_EXCLUDED_EMAIL_ADDRESSES       = "excludedemailaddresses"
	KEYWORD_EXCLUDED_IP_RANGES             = "excludedipranges"
	KEYWORD_EXCLUDED_URI_DOMAINS           = "excludeduridomains"
	KEYWORD_IP_ADDRESSES                   = "ipaddresses"
	KEYWORD_IS_CA                          = "isca"
	KEYWORD_ISSUING_CERTIFICATE_URL        = "issuingcertificateurl"
	KEYWORD_KEY_USAGE                      = "keyusage"
	KEYWORD_MAX_PATH_LEN                   = "maxpathlen"
	KEYWORD_MAX_PATH_LEN_ZERO              = "maxpathlenzero"
	KEYWORD_NOT_AFTER                      = "notafter"
	KEYWORD_NOT_BEFORE                     = "notbefore"
	KEYWORD_O_C_S_P_SERVER                 = "ocspserver"
	KEYWORD_PERMITTED_DNS_DOMAINS          = "permitteddnsdomains"
	KEYWORD_PERMITTED_DNS_DOMAINS_CRITICAL = "permitteddnsdomainscritical"
	KEYWORD_PERMITTED_EMAIL_ADDRESSES      = "permittedemailaddresses"
	KEYWORD_PERMITTED_IP_RANGES            = "permittedipranges"
	KEYWORD_PERMITTED_URI_DOMAINS          = "permitteduridomains"
	KEYWORD_POLICY_IDENTIFIERS             = "policyidentifiers"
	KEYWORD_POLICIES                       = "policies"
	KEYWORD_SERIAL_NUMBER                  = "serialnumber"
	KEYWORD_SIGNATURE_ALGORITHM            = "signaturealgorithm"
	KEYWORD_SUBJECT                        = "subject"
	KEYWORD_SUBJECT_KEY_ID                 = "subjectkeyid"
	KEYWORD_URIS                           = "uris"
	KEYWORD_UNKNOWN_EXT_KEY_USAGE          = "unknownextkeyusage"
)

var extKeyUsageValue = []struct {
	extKeyUsage x509.ExtKeyUsage
	key         string
}{
	{x509.ExtKeyUsageAny, KEYWORD_EXT_KEYUSAGE_ANY},
	{x509.ExtKeyUsageServerAuth, KEYWORD_EXT_KEYUSAGE_SERVER_AUTH},
	{x509.ExtKeyUsageClientAuth, KEYWORD_EXT_KEYUSAGE_CLIENT_AUTH},
	{x509.ExtKeyUsageCodeSigning, KEYWORD_EXT_KEYUSAGE_CODE_SIGNING},
	{x509.ExtKeyUsageEmailProtection, KEYWORD_EXT_KEYUSAGE_EMAIL_PROTECTION},
	{x509.ExtKeyUsageIPSECEndSystem, KEYWORD_EXT_KEYUSAGE_IPSEC_END_SYSTEM},
	{x509.ExtKeyUsageIPSECTunnel, KEYWORD_EXT_KEYUSAGE_IPSEC_TUNNEL},
	{x509.ExtKeyUsageIPSECUser, KEYWORD_EXT_KEYUSAGE_IPSEC_USER},
	{x509.ExtKeyUsageTimeStamping, KEYWORD_EXT_KEYUSAGE_TIME_STAMPING},
	{x509.ExtKeyUsageOCSPSigning, KEYWORD_EXT_KEYUSAGE_OCSP_SIGNING},
	{x509.ExtKeyUsageMicrosoftServerGatedCrypto, KEYWORD_EXT_KEYUSAGE_MICROSOFT_SERVER_GATED_CRYPTO},
	{x509.ExtKeyUsageNetscapeServerGatedCrypto, KEYWORD_EXT_KEYUSAGE_NETSCAPE_SERVER_GATED_CRYPTO},
	{x509.ExtKeyUsageMicrosoftCommercialCodeSigning, KEYWORD_EXT_KEYUSAGE_MICROSOFT_COMMERCIAL_CODE_SIGNING},
	{x509.ExtKeyUsageMicrosoftKernelCodeSigning, KEYWORD_EXT_KEYUSAGE_MICROSOFT_KERNEL_CODE_SIGNING},
}

var keyUsageValue = []struct {
	value x509.KeyUsage
	key   string
}{
	{x509.KeyUsageDigitalSignature, KEYWORD_DIGITAL_SIGNATURE},
	{x509.KeyUsageContentCommitment, KEYWORD_CONTENT_COMMITMENT},
	{x509.KeyUsageKeyEncipherment, KEYWORD_KEY_ENCIPHERMENT},
	{x509.KeyUsageDataEncipherment, KEYWORD_DATA_ENCIPHERMENT},
	{x509.KeyUsageKeyAgreement, KEYWORD_KEY_AGREEMENT},
	{x509.KeyUsageCertSign, KEYWORD_CERT_SIGN},
	{x509.KeyUsageCRLSign, KEYWORD_CRL_SIGN},
	{x509.KeyUsageEncipherOnly, KEYWORD_ENCIPHER_ONLY},
	{x509.KeyUsageDecipherOnly, KEYWORD_DECIPHER_ONLY},
}

const (
	KEYWORD_DNSNAME                      = "dnsname"
	KEYWORD_INTERMEDIATES                = "intermediates"
	KEYWORD_ROOTS                        = "roots"
	KEYWORD_CURRENTTIME                  = "currenttime"
	KEYWORD_MAX_CONSTRAINTS_COMPARISIONS = "maxconstraintscomparisions"
)

// errNotParsed is returned when a certificate without ASN.1 contents is
// verified. Platform-specific verification needs the ASN.1 contents.
var errNotParsed = errors.New("x509: missing ASN.1 contents; use ParseCertificate")

// errNotParsed is returned when a certificate without ASN.1 contents is
// verified. Platform-specific verification needs the ASN.1 contents.
var errNoRootError = errors.New("x509: no root set")

const (
	leafCertificate = iota
	intermediateCertificate
	rootCertificate
)

var (
	pssParametersSHA256 = asn1.RawValue{FullBytes: []byte{48, 52, 160, 15, 48, 13, 6, 9, 96, 134, 72, 1, 101, 3, 4, 2, 1, 5, 0, 161, 28, 48, 26, 6, 9, 42, 134, 72, 134, 247, 13, 1, 1, 8, 48, 13, 6, 9, 96, 134, 72, 1, 101, 3, 4, 2, 1, 5, 0, 162, 3, 2, 1, 32}}
	pssParametersSHA384 = asn1.RawValue{FullBytes: []byte{48, 52, 160, 15, 48, 13, 6, 9, 96, 134, 72, 1, 101, 3, 4, 2, 2, 5, 0, 161, 28, 48, 26, 6, 9, 42, 134, 72, 134, 247, 13, 1, 1, 8, 48, 13, 6, 9, 96, 134, 72, 1, 101, 3, 4, 2, 2, 5, 0, 162, 3, 2, 1, 48}}
	pssParametersSHA512 = asn1.RawValue{FullBytes: []byte{48, 52, 160, 15, 48, 13, 6, 9, 96, 134, 72, 1, 101, 3, 4, 2, 3, 5, 0, 161, 28, 48, 26, 6, 9, 42, 134, 72, 134, 247, 13, 1, 1, 8, 48, 13, 6, 9, 96, 134, 72, 1, 101, 3, 4, 2, 3, 5, 0, 162, 3, 2, 1, 64}}
)

const (
	nameTypeEmail = 1
	nameTypeDNS   = 2
	nameTypeURI   = 6
	nameTypeIP    = 7
)

var (
	oidAuthorityInfoAccessOcsp    = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 48, 1}
	oidAuthorityInfoAccessIssuers = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 48, 2}
)

var (
	oidExtensionSubjectKeyId          = []int{2, 5, 29, 14}
	oidExtensionKeyUsage              = []int{2, 5, 29, 15}
	oidExtensionExtendedKeyUsage      = []int{2, 5, 29, 37}
	oidExtensionAuthorityKeyId        = []int{2, 5, 29, 35}
	oidExtensionBasicConstraints      = []int{2, 5, 29, 19}
	oidExtensionSubjectAltName        = []int{2, 5, 29, 17}
	oidExtensionCertificatePolicies   = []int{2, 5, 29, 32}
	oidExtensionNameConstraints       = []int{2, 5, 29, 30}
	oidExtensionCRLDistributionPoints = []int{2, 5, 29, 31}
	oidExtensionAuthorityInfoAccess   = []int{1, 3, 6, 1, 5, 5, 7, 1, 1}
	oidExtensionCRLNumber             = []int{2, 5, 29, 20}
	oidExtensionReasonCode            = []int{2, 5, 29, 21}
)

var (
	oidSignatureMD5WithRSA      = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 4}
	oidSignatureSHA1WithRSA     = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 5}
	oidSignatureSHA256WithRSA   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 11}
	oidSignatureSHA384WithRSA   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 12}
	oidSignatureSHA512WithRSA   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 13}
	oidSignatureRSAPSS          = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 10}
	oidSignatureDSAWithSHA1     = asn1.ObjectIdentifier{1, 2, 840, 10040, 4, 3}
	oidSignatureDSAWithSHA256   = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 3, 2}
	oidSignatureECDSAWithSHA1   = asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 1}
	oidSignatureECDSAWithSHA256 = asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 2}
	oidSignatureECDSAWithSHA384 = asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 3}
	oidSignatureECDSAWithSHA512 = asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 4}
	oidSignatureEd25519         = asn1.ObjectIdentifier{1, 3, 101, 112}

	oidSHA256 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}
	oidSHA384 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 2}
	oidSHA512 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 3}

	oidMGF1 = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 8}

	// oidISOSignatureSHA1WithRSA means the same as oidSignatureSHA1WithRSA
	// but it's specified by ISO. Microsoft's makecert.exe has been known
	// to produce certificates with this OID.
	oidISOSignatureSHA1WithRSA = asn1.ObjectIdentifier{1, 3, 14, 3, 2, 29}
)

// RFC 5280, 4.2.1.12  Extended Key Usage
//
//	anyExtendedKeyUsage OBJECT IDENTIFIER ::= { id-ce-extKeyUsage 0 }
//
//	id-kp OBJECT IDENTIFIER ::= { id-pkix 3 }
//
//	id-kp-serverAuth             OBJECT IDENTIFIER ::= { id-kp 1 }
//	id-kp-clientAuth             OBJECT IDENTIFIER ::= { id-kp 2 }
//	id-kp-codeSigning            OBJECT IDENTIFIER ::= { id-kp 3 }
//	id-kp-emailProtection        OBJECT IDENTIFIER ::= { id-kp 4 }
//	id-kp-timeStamping           OBJECT IDENTIFIER ::= { id-kp 8 }
//	id-kp-OCSPSigning            OBJECT IDENTIFIER ::= { id-kp 9 }
var (
	oidExtKeyUsageAny                            = asn1.ObjectIdentifier{2, 5, 29, 37, 0}
	oidExtKeyUsageServerAuth                     = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 1}
	oidExtKeyUsageClientAuth                     = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 2}
	oidExtKeyUsageCodeSigning                    = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 3}
	oidExtKeyUsageEmailProtection                = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 4}
	oidExtKeyUsageIPSECEndSystem                 = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 5}
	oidExtKeyUsageIPSECTunnel                    = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 6}
	oidExtKeyUsageIPSECUser                      = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 7}
	oidExtKeyUsageTimeStamping                   = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 8}
	oidExtKeyUsageOCSPSigning                    = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 9}
	oidExtKeyUsageMicrosoftServerGatedCrypto     = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 10, 3, 3}
	oidExtKeyUsageNetscapeServerGatedCrypto      = asn1.ObjectIdentifier{2, 16, 840, 1, 113730, 4, 1}
	oidExtKeyUsageMicrosoftCommercialCodeSigning = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 22}
	oidExtKeyUsageMicrosoftKernelCodeSigning     = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 61, 1, 1}
)

// oidExtensionRequest is a PKCS #9 OBJECT IDENTIFIER that indicates requested
// extensions in a CSR.
var oidExtensionRequest = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 14}
