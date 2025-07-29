package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"dbgutil"
	"encoding/asn1"
	"fileop"
	"fmt"
	"jsonext"
	"logutil"
	"math/big"
	"net"
	"net/url"
	"strings"
	"time"
)

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
)

func get_extra_names(valarr []interface{}) (retv []pkix.AttributeTypeAndValue, err error) {
	var curmap map[string]interface{}
	var idx, jdx int
	var ok bool
	var carr []interface{}
	var curoid asn1.ObjectIdentifier
	var curattr pkix.AttributeTypeAndValue
	var curi int
	var vali int
	var iarr asn1.RawValue
	var valmap map[string]interface{}
	var intval interface{}
	retv = []pkix.AttributeTypeAndValue{}
	for idx = 0; idx < len(valarr); idx += 1 {
		curmap, ok = valarr[idx].(map[string]interface{})
		if !ok {
			err = dbgutil.FormatError("[%d] not valid map[string]interface{}", idx)
			return
		}

		carr, ok = curmap[KEYWORD_TYPE].([]interface{})
		if !ok {
			err = dbgutil.FormatError("[%d].[%s] not array", idx, KEYWORD_TYPE)
			return
		}

		curattr = pkix.AttributeTypeAndValue{}
		curoid = asn1.ObjectIdentifier{}
		for jdx = 0; jdx < len(carr); jdx += 1 {
			curi, err = get_int_value(carr[jdx], fmt.Sprintf("[%s].[%d]", KEYWORD_TYPE, jdx))
			if err != nil {
				return
			}
			logutil.Debug("curi %d", curi)
			curoid = append(curoid, curi)
		}

		curattr.Type = curoid
		valmap, ok = curmap[KEYWORD_VALUE].(map[string]interface{})
		if !ok {
			curattr.Value = nil
		} else {
			iarr = asn1.RawValue{}
			iarr.Class = 0
			iarr.Tag = 0
			iarr.IsCompound = false
			iarr.Bytes = []byte{}
			iarr.FullBytes = []byte{}

			intval, ok = valmap[KEYWORD_TAG]
			if !ok {
				err = dbgutil.FormatError("[%s][%d] no %s", KEYWORD_VALUE, idx, KEYWORD_TAG)
				return
			}
			vali, err = get_int_value(intval, fmt.Sprintf("[%s].[%d]", KEYWROD_EXTRANAMES, idx))
			if err != nil {
				return
			}

			iarr.Tag = (vali & 0x1f)
			if (vali & 0x20) != 0 {
				iarr.IsCompound = true
			}
			iarr.Class = ((vali >> 6) & 0x3)

			carr, ok = valmap[KEYWORD_CONTENT].([]interface{})
			if !ok {
				err = dbgutil.FormatError("[%s][%d] no %s", KEYWORD_VALUE, idx, KEYWORD_CONTENT)
				return
			}
			iarr.Bytes, err = get_bytes_value(carr, fmt.Sprintf("[%s][%d] [%s]", KEYWORD_VALUE, idx, KEYWORD_CONTENT))
			if err != nil {
				return
			}

			curattr.Value = iarr

		}

		retv = append(retv, curattr)
	}

	err = nil
	return
}

func get_array_string(mapv map[string]interface{}, key string) (retv []string) {
	//var ok bool
	var err error
	var idx int
	var valarr []interface{}
	retv = []string{}

	valarr, err = jsonext.GetJsonValueArray(key, mapv)
	if err != nil {
		logutil.Debug("[%s] not ok %s", key, err.Error())
		err = nil
		return
	}
	for idx = 0; idx < len(valarr); idx += 1 {
		retv = append(retv, valarr[idx].(string))
	}
	return
}

func get_pkixname_value(mapv map[string]interface{}, note string) (name pkix.Name, err error) {

	var valarr []interface{}
	var valinter interface{}
	var ok bool
	name = pkix.Name{}

	name.Country = get_array_string(mapv, KEYWORD_COUNTRY)
	name.Province = get_array_string(mapv, KEYWORD_PROVINCE)
	name.Locality = get_array_string(mapv, KEYWORD_LOCALITY)
	name.StreetAddress = get_array_string(mapv, KEYWORD_STREETADDRESS)
	name.PostalCode = get_array_string(mapv, KEYWORD_POSTALCODE)
	name.Organization = get_array_string(mapv, KEYWORD_ORGANIZATION)
	name.OrganizationalUnit = get_array_string(mapv, KEYWORD_ORGANIZATIONUNIT)

	valinter, ok = mapv[KEYWORD_COMMONNAME]
	if ok {
		name.CommonName = valinter.(string)
	} else {
		name.CommonName = ""
	}

	valinter, ok = mapv[KEYWORD_SERIALNUMBER]
	if ok {
		name.SerialNumber = valinter.(string)
	} else {
		name.SerialNumber = ""
	}

	valarr, ok = mapv[KEYWROD_EXTRANAMES].([]interface{})
	if ok {
		name.ExtraNames, err = get_extra_names(valarr)
		if err != nil {
			return
		}
	} else {
		name.ExtraNames = []pkix.AttributeTypeAndValue{}
	}

	err = nil
	return

}

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

func get_key_ext_usage(arrs []string, note string) (retv []x509.ExtKeyUsage, err error) {
	retv = []x509.ExtKeyUsage{}
	var idx, jdx int
	var matched bool
	for idx = 0; idx < len(arrs); idx += 1 {
		matched = false
		for jdx = 0; jdx < len(extKeyUsageValue); jdx += 1 {
			if extKeyUsageValue[jdx].key == arrs[idx] {
				matched = true
				retv = append(retv, extKeyUsageValue[jdx].extKeyUsage)
				break
			}
		}

		if !matched {
			err = dbgutil.FormatError("[%s].[%d][%s] not supported", note, idx, arrs[idx])
			return
		}
	}
	err = nil
	return
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

func get_keyusage_value(arrs []string, note string) (retv x509.KeyUsage, err error) {
	retv = 0
	var jdx int
	var idx int
	var matched bool
	var vali int = 0
	for idx = 0; idx < len(arrs); idx += 1 {
		matched = false
		for jdx = 0; jdx < len(keyUsageValue); jdx += 1 {
			if keyUsageValue[jdx].key == arrs[idx] {
				matched = true
				vali |= int(keyUsageValue[jdx].value)
				break
			}
		}

		if !matched {
			err = dbgutil.FormatError("[%s].[%d] [%s] not valid key usage", note, idx, arrs[idx])
			return
		}

	}

	err = nil
	retv = x509.KeyUsage(vali)
	return
}

func get_bytes_value(valarr []interface{}, note string) (retv []byte, err error) {
	var idx int
	var valf float64
	var vali int
	var ok bool
	retv = []byte{}

	for idx = 0; idx < len(valarr); idx += 1 {
		vali, ok = valarr[idx].(int)
		if !ok {
			valf, ok = valarr[idx].(float64)
			if !ok {
				err = dbgutil.FormatError("[%s].[%d] not valid", note, idx)
				return
			}
			vali = int(valf)
		}

		if vali < 0 || vali > 255 {
			err = dbgutil.FormatError("[%s].[%d] %d not byte value", note, idx, vali)
			return
		}

		retv = append(retv, byte(vali))
	}

	err = nil
	return
}

func get_objoid_array(valarr []interface{}, note string) (retv asn1.ObjectIdentifier, err error) {
	var idx int
	var valf float64
	var vali int
	var ok bool
	retv = asn1.ObjectIdentifier{}

	for idx = 0; idx < len(valarr); idx += 1 {
		vali, ok = valarr[idx].(int)
		if !ok {
			valf, ok = valarr[idx].(float64)
			if !ok {
				err = dbgutil.FormatError("[%s].[%d] not valid", note, idx)
				return
			}
			vali = int(valf)
		}

		retv = append(retv, vali)
	}

	err = nil
	return
}

func get_int_value(val interface{}, note string) (ival int, err error) {
	var vali int
	var valf float64
	var ok bool
	vali, ok = val.(int)
	if ok {
		ival = vali
	} else {
		valf, ok = val.(float64)
		if ok {
			ival = int(valf)
		} else {
			err = dbgutil.FormatError("[%s] not int value", note)
			return
		}
	}
	err = nil
	return
}

func get_time_value(times string, note string) (retv time.Time, err error) {
	retv, err = time.Parse("2020-02-02 13:20:50", times)
	if err != nil {
		err = dbgutil.FormatError("[%s] [%s] parse error %s please use [2020-02-02 13:20:50] format", note, times, err.Error())
		return
	}
	return
}

func get_netip_value(val []interface{}, note string) (retv []*net.IPNet, err error) {
	var idx int
	var ok bool
	var s string
	var curnet *net.IPNet
	retv = []*net.IPNet{}
	for idx = 0; idx < len(val); idx += 1 {
		s, ok = val[idx].(string)
		if !ok {
			err = dbgutil.FormatError("[%s].[%d] not valid string", note, idx)
			return
		}
		_, curnet, err = net.ParseCIDR(s)
		if err != nil {
			err = dbgutil.FormatError("[%s].[%d] parse [%s] error %s", note, idx, s, err.Error())
			return
		}
		retv = append(retv, curnet)
	}
	err = nil
	return
}

func get_ip_value(val []interface{}, note string) (retv []net.IP, err error) {
	var idx int
	var ok bool
	var s string
	var curnet net.IP
	retv = []net.IP{}
	for idx = 0; idx < len(val); idx += 1 {
		s, ok = val[idx].(string)
		if !ok {
			err = dbgutil.FormatError("[%s].[%d] not valid string", note, idx)
			return
		}
		curnet = net.ParseIP(s)
		retv = append(retv, curnet)
	}
	err = nil
	return
}

func get_objoids_value(val []interface{}, note string) (retv []asn1.ObjectIdentifier, err error) {
	var idx int
	var curoid asn1.ObjectIdentifier
	var curinter []interface{}
	var ok bool
	retv = []asn1.ObjectIdentifier{}
	for idx = 0; idx < len(val); idx += 1 {
		curinter, ok = val[idx].([]interface{})
		if !ok {
			err = dbgutil.FormatError("[%s].[%d] not array type", note, idx)
			return
		}
		curoid, err = get_objoid_array(curinter, fmt.Sprintf("[%s].[%d]", note, idx))
		if err != nil {
			return
		}
		retv = append(retv, curoid)
	}
	err = nil
	return
}

func get_bigint_value(val string, note string) (retv *big.Int, err error) {
	var base int = 10
	var inputs string = val
	var ok bool
	if strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
		inputs = val[2:]
		base = 16
	} else if strings.HasPrefix(val, "x") || strings.HasPrefix(val, "X") {
		inputs = val[1:]
		base = 16
	}
	retv = big.NewInt(0)
	_, ok = retv.SetString(inputs, base)
	if !ok {
		err = dbgutil.FormatError("[%s] [%s] not valid big.Int ", note, val)
		return
	}
	err = nil
	return
}

var matchAlgorithm = []struct {
	key string
	val x509.SignatureAlgorithm
}{
	{KEYWORD_MD2_WITH_RSA, x509.MD2WithRSA},
	{KEYWORD_MD5_WITH_RSA, x509.MD5WithRSA},
	{KEYWORD_SHA1_WITH_RSA, x509.SHA1WithRSA},
	{KEYWORD_SHA256_WITH_RSA, x509.SHA256WithRSA},
	{KEYWORD_SHA384_WITH_RSA, x509.SHA384WithRSA},
	{KEYWORD_SHA512_WITH_RSA, x509.SHA512WithRSA},
	{KEYWORD_DSA_WITH_SHA1, x509.DSAWithSHA1},
	{KEYWORD_DSA_WITH_SHA256, x509.DSAWithSHA256},
	{KEYWORD_ECDSA_WITH_SHA1, x509.ECDSAWithSHA1},
	{KEYWORD_ECDSA_WITH_SHA256, x509.ECDSAWithSHA256},
	{KEYWORD_ECDSA_WITH_SHA384, x509.ECDSAWithSHA384},
	{KEYWORD_ECDSA_WITH_SHA512, x509.ECDSAWithSHA512},
	{KEYWORD_SHA256_WITH_RSAPSS, x509.SHA256WithRSAPSS},
	{KEYWORD_SHA384_WITH_RSAPSS, x509.SHA384WithRSAPSS},
	{KEYWORD_SHA512_WITH_RSAPSS, x509.SHA512WithRSAPSS},
	{KEYWORD_PURE_ED25519, x509.PureEd25519},
}

func get_algorithm_value(val string, note string) (retv x509.SignatureAlgorithm, err error) {
	var idx int
	retv = x509.UnknownSignatureAlgorithm
	for idx = 0; idx < len(matchAlgorithm); idx += 1 {
		if val == matchAlgorithm[idx].key {
			retv = matchAlgorithm[idx].val
			err = nil
			return
		}
	}

	err = dbgutil.FormatError("[%s] val [%s] not support", note, val)
	return
}

func get_oids_value(val []string, note string) (retv []x509.OID, err error) {
	var idx int
	var curoid x509.OID
	retv = []x509.OID{}

	for idx = 0; idx < len(val); idx += 1 {
		curoid, err = x509.ParseOID(val[idx])
		if err != nil {
			err = dbgutil.FormatError("[%s].[%d] [%s] error %s", note, idx, val[idx], err.Error())
			return
		}
		retv = append(retv, curoid)
	}
	err = nil
	return
}

func get_urls_value(val []string, note string) (retv []*url.URL, err error) {
	var idx int
	var cururl *url.URL
	retv = []*url.URL{}

	for idx = 0; idx < len(val); idx += 1 {
		cururl, err = url.Parse(val[idx])
		if err != nil {
			err = dbgutil.FormatError("[%s].[%d] [%s] error %s", note, idx, val[idx], err.Error())
			return
		}
		retv = append(retv, cururl)
	}
	err = nil
	return
}

func trans_inter_to_string(valarr []interface{}, note string) (arrs []string, err error) {
	arrs = []string{}
	var idx int
	var s string
	var ok bool
	for idx = 0; idx < len(valarr); idx += 1 {
		s, ok = valarr[idx].(string)
		if !ok {
			err = dbgutil.FormatError("[%s].[%s] not string", note, idx)
			return
		}
		arrs = append(arrs, s)
	}
	err = nil
	return
}

func get_certificate_file(f string) (x509temp x509.Certificate, err error) {
	var s string
	var ok bool
	var mapv map[string]interface{}
	var arrs []string
	var valarr []interface{}
	var intval interface{}
	var valb bool
	var vals string
	var valmap map[string]interface{}

	x509temp = x509.Certificate{}
	s, err = fileop.ReadFile(f)
	if err != nil {
		return
	}
	mapv, err = jsonext.GetJsonMap(s)
	if err != nil {
		return
	}

	x509temp.ExtKeyUsage = []x509.ExtKeyUsage{}
	valarr, ok = mapv[KEYWORD_EXT_KEY_USAGE].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_EXT_KEY_USAGE)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_EXT_KEY_USAGE)
		if err != nil {
			return
		}
		x509temp.ExtKeyUsage, err = get_key_ext_usage(arrs, KEYWORD_EXT_KEY_USAGE)
		if err != nil {
			return
		}
	}

	x509temp.AuthorityKeyId = []byte{}
	valarr, ok = mapv[KEYWORD_AUTHORITY_KEY_ID].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_AUTHORITY_KEY_ID)
		x509temp.AuthorityKeyId, err = get_bytes_value(valarr, KEYWORD_AUTHORITY_KEY_ID)
		if err != nil {
			return
		}
	}

	x509temp.BasicConstraintsValid = false
	valb, ok = mapv[KEYWORD_BASIC_CONSTRAINTS_VALID].(bool)
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_BASIC_CONSTRAINTS_VALID)
		x509temp.BasicConstraintsValid = valb
	}

	x509temp.CRLDistributionPoints = []string{}
	valarr, ok = mapv[KEYWORD_CRL_DISTRIBUTION_POINTS].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_CRL_DISTRIBUTION_POINTS)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_CRL_DISTRIBUTION_POINTS)
		if err != nil {
			return
		}
		x509temp.CRLDistributionPoints = arrs
	}

	x509temp.DNSNames = []string{}
	valarr, ok = mapv[KEYWORD_DNS_NAMES].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_DNS_NAMES)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_DNS_NAMES)
		if err != nil {
			return
		}
		x509temp.DNSNames = arrs
	}

	x509temp.EmailAddresses = []string{}
	valarr, ok = mapv[KEYWORD_EMAIL_ADDRESSES].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_EMAIL_ADDRESSES)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_EMAIL_ADDRESSES)
		if err != nil {
			return
		}
		x509temp.EmailAddresses = arrs
	}

	x509temp.ExcludedDNSDomains = []string{}
	valarr, ok = mapv[KEYWORD_EXCLUDED_DNS_DOMAINS].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_EXCLUDED_DNS_DOMAINS)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_EXCLUDED_DNS_DOMAINS)
		if err != nil {
			return
		}
		x509temp.ExcludedDNSDomains = arrs
	}

	x509temp.ExcludedEmailAddresses = []string{}
	valarr, ok = mapv[KEYWORD_EXCLUDED_EMAIL_ADDRESSES].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_EXCLUDED_EMAIL_ADDRESSES)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_EXCLUDED_EMAIL_ADDRESSES)
		if err != nil {
			return
		}
		x509temp.ExcludedEmailAddresses = arrs
	}

	x509temp.ExcludedIPRanges = []*net.IPNet{}
	valarr, ok = mapv[KEYWORD_EXCLUDED_IP_RANGES].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_EXCLUDED_IP_RANGES)
		x509temp.ExcludedIPRanges, err = get_netip_value(valarr, KEYWORD_EXCLUDED_IP_RANGES)
		if err != nil {
			return
		}
	}

	x509temp.ExcludedURIDomains = []string{}
	valarr, ok = mapv[KEYWORD_EXCLUDED_URI_DOMAINS].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_EXCLUDED_URI_DOMAINS)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_EXCLUDED_URI_DOMAINS)
		if err != nil {
			return
		}
		x509temp.ExcludedURIDomains = arrs
	}

	x509temp.IPAddresses = []net.IP{}
	valarr, ok = mapv[KEYWORD_IP_ADDRESSES].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_IP_ADDRESSES)
		x509temp.IPAddresses, err = get_ip_value(valarr, KEYWORD_IP_ADDRESSES)
		if err != nil {
			return
		}
	}

	x509temp.IsCA = false
	valb, ok = mapv[KEYWORD_IS_CA].(bool)
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_IS_CA)
		x509temp.IsCA = valb
	}

	x509temp.IssuingCertificateURL = []string{}
	valarr, ok = mapv[KEYWORD_ISSUING_CERTIFICATE_URL].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_ISSUING_CERTIFICATE_URL)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_ISSUING_CERTIFICATE_URL)
		if err != nil {
			return
		}
		x509temp.IssuingCertificateURL = arrs
	}

	x509temp.KeyUsage = 0
	valarr, ok = mapv[KEYWORD_KEY_USAGE].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_KEY_USAGE)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_KEY_USAGE)
		if err != nil {
			return
		}
		x509temp.KeyUsage, err = get_keyusage_value(arrs, KEYWORD_KEY_USAGE)
		if err != nil {
			return
		}
	}

	x509temp.MaxPathLen = 0
	intval, ok = mapv[KEYWORD_MAX_PATH_LEN]
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_MAX_PATH_LEN)
		x509temp.MaxPathLen, err = get_int_value(intval, KEYWORD_MAX_PATH_LEN)
		if err != nil {
			return
		}
	}

	x509temp.MaxPathLenZero = false
	valb, ok = mapv[KEYWORD_MAX_PATH_LEN_ZERO].(bool)
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_MAX_PATH_LEN_ZERO)
		x509temp.MaxPathLenZero = valb
	}

	x509temp.NotAfter = time.Now().AddDate(20, 0, 0)
	vals, ok = mapv[KEYWORD_NOT_AFTER].(string)
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_NOT_AFTER)
		x509temp.NotAfter, err = get_time_value(vals, KEYWORD_NOT_AFTER)
		if err != nil {
			return
		}
	}

	x509temp.NotBefore = time.Now()
	vals, ok = mapv[KEYWORD_NOT_BEFORE].(string)
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_NOT_BEFORE)
		x509temp.NotBefore, err = get_time_value(vals, KEYWORD_NOT_BEFORE)
		if err != nil {
			return
		}
	}

	x509temp.OCSPServer = []string{}
	valarr, ok = mapv[KEYWORD_O_C_S_P_SERVER].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_O_C_S_P_SERVER)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_O_C_S_P_SERVER)
		if err != nil {
			return
		}
		x509temp.OCSPServer = arrs
	}

	x509temp.PermittedDNSDomains = []string{}
	valarr, ok = mapv[KEYWORD_PERMITTED_DNS_DOMAINS].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_PERMITTED_DNS_DOMAINS)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_PERMITTED_DNS_DOMAINS)
		if err != nil {
			return
		}
		x509temp.PermittedDNSDomains = arrs
	}

	x509temp.PermittedDNSDomainsCritical = false
	valb, ok = mapv[KEYWORD_PERMITTED_DNS_DOMAINS_CRITICAL].(bool)
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_PERMITTED_DNS_DOMAINS_CRITICAL)
		x509temp.PermittedDNSDomainsCritical = valb
	}

	x509temp.PermittedEmailAddresses = []string{}
	valarr, ok = mapv[KEYWORD_PERMITTED_EMAIL_ADDRESSES].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_PERMITTED_EMAIL_ADDRESSES)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_PERMITTED_EMAIL_ADDRESSES)
		if err != nil {
			return
		}
		x509temp.PermittedEmailAddresses = arrs
	}

	x509temp.PermittedIPRanges = []*net.IPNet{}
	valarr, ok = mapv[KEYWORD_PERMITTED_IP_RANGES].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_PERMITTED_IP_RANGES)
		x509temp.PermittedIPRanges, err = get_netip_value(valarr, KEYWORD_PERMITTED_IP_RANGES)
		if err != nil {
			return
		}
	}

	x509temp.PermittedURIDomains = []string{}
	valarr, ok = mapv[KEYWORD_PERMITTED_URI_DOMAINS].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_PERMITTED_URI_DOMAINS)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_PERMITTED_URI_DOMAINS)
		if err != nil {
			return
		}
		x509temp.PermittedURIDomains = arrs
	}

	x509temp.PolicyIdentifiers = []asn1.ObjectIdentifier{}
	valarr, ok = mapv[KEYWORD_POLICY_IDENTIFIERS].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_POLICY_IDENTIFIERS)
		x509temp.PolicyIdentifiers, err = get_objoids_value(valarr, KEYWORD_POLICY_IDENTIFIERS)
		if err != nil {
			return
		}
	}

	x509temp.Policies = []x509.OID{}
	valarr, ok = mapv[KEYWORD_POLICIES].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_POLICIES)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_POLICIES)
		if err != nil {
			return
		}
		x509temp.Policies, err = get_oids_value(arrs, KEYWORD_POLICIES)
		if err != nil {
			return
		}
	}

	x509temp.SerialNumber = big.NewInt(0)
	vals, ok = mapv[KEYWORD_SERIAL_NUMBER].(string)
	if ok {
		x509temp.SerialNumber, err = get_bigint_value(vals, KEYWORD_SERIAL_NUMBER)
		if err != nil {
			return
		}
	}

	x509temp.SignatureAlgorithm = x509.SHA256WithRSA
	vals, ok = mapv[KEYWORD_SIGNATURE_ALGORITHM].(string)
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_SIGNATURE_ALGORITHM)
		x509temp.SignatureAlgorithm, err = get_algorithm_value(vals, KEYWORD_SIGNATURE_ALGORITHM)
		if err != nil {
			return
		}
	}

	x509temp.Subject = pkix.Name{}
	valmap, ok = mapv[KEYWORD_SUBJECT].(map[string]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_SUBJECT)
		x509temp.Subject, err = get_pkixname_value(valmap, KEYWORD_SUBJECT)
		if err != nil {
			return
		}
	}

	x509temp.SubjectKeyId = []byte{}
	valarr, ok = mapv[KEYWORD_SUBJECT_KEY_ID].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_SUBJECT_KEY_ID)
		x509temp.SubjectKeyId, err = get_bytes_value(valarr, KEYWORD_SUBJECT_KEY_ID)
		if err != nil {
			return
		}
	}

	x509temp.URIs = []*url.URL{}
	valarr, ok = mapv[KEYWORD_URIS].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_URIS)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_URIS)
		if err != nil {
			return
		}
		x509temp.URIs, err = get_urls_value(arrs, KEYWORD_URIS)
		if err != nil {
			return
		}
	}

	x509temp.UnknownExtKeyUsage = []asn1.ObjectIdentifier{}
	valarr, ok = mapv[KEYWORD_UNKNOWN_EXT_KEY_USAGE].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_UNKNOWN_EXT_KEY_USAGE)
		x509temp.UnknownExtKeyUsage, err = get_objoids_value(valarr, KEYWORD_UNKNOWN_EXT_KEY_USAGE)
		if err != nil {
			return
		}
	}

	err = nil
	return

}
