package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"dbgutil"
	"encoding/asn1"
	"fileop"
	"jsonext"
	"logutil"
	"math/big"
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
)

func get_extra_names(valarr []interface{}) (retv []pkix.AttributeTypeAndValue, err error) {
	var curmap map[string]interface{}
	var idx, jdx int
	var ok bool
	var carr []interface{}
	var curoid asn1.ObjectIdentifier
	var curattr pkix.AttributeTypeAndValue
	var curi int
	var curf float64
	var iarr asn1.RawContent
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
			curi = 0
			curi, ok = carr[jdx].(int)
			if !ok {
				curf, ok = carr[jdx].(float64)
				if ok {
					curi = int(curf)
				}
			} else {
				curi = 0
			}
			logutil.Debug("curi %d", curi)
			curoid = append(curoid, curi)
		}

		curattr.Type = curoid
		carr, ok = curmap[KEYWORD_VALUE].([]interface{})
		if !ok {
			curattr.Value = nil
		} else {
			iarr = asn1.RawContent{}

			for jdx = 0; jdx < len(carr); jdx += 1 {
				curi = 0
				curi, ok = carr[jdx].(int)
				if !ok {
					curf, ok = carr[jdx].(float64)
					if ok {
						curi = int(curf)
					}
				} else {
					curi = 0
				}
				logutil.Debug("curi %d", curi)
				iarr = append(iarr, byte(curi))
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

func get_pkix_name_mapv(mapv map[string]interface{}) (name pkix.Name, err error) {

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
	KEYWORD_EXTKEYUSAGE                                    = "extkeyusage"
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

	KEYWORD_MD2_WITH_RSA = "md2withrsa"
	KEYWORD_MD5_WITH_RSA = "md5withrsa"
	KEYWORD_SHA1_WITH_RSA = "sha1withrsa"
	KEYWORD_SHA256_WITH_RSA = "sha256withrsa"
	SHA384WithRSA
	SHA512WithRSA
	DSAWithSHA1   // Unsupported.
	DSAWithSHA256 // Unsupported.
	ECDSAWithSHA1 // Only supported for signing, and verification of CRLs, CSRs, and OCSP responses.
	ECDSAWithSHA256
	ECDSAWithSHA384
	ECDSAWithSHA512
	SHA256WithRSAPSS
	SHA384WithRSAPSS
	SHA512WithRSAPSS
	PureEd25519


	KEYWORD_AUTHORITY_KEY_ID        = "authoritykeyid"
	KEYWORD_BASIC_CONSTRAINT_VALID  = "basiccontraintsvalid"
	KEYWORD_CRL_DISTRIBUTION_POINTS = "crldistributionpoints"
	KEYWORD_DNS_NAMES               = "dnsnames"
	KEYWORD_EMAIL_ADDRESSES         = "emailaddresses"
	KEYWORD_EXCLUDED_DNS_DOMAINS    = "excludeddnsdomains"
	KEYWORD_EXCLUDED_EMAILADDRESSES = "excludeemailaddresses"
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

func get_key_ext_usage(arrs []string) (retv []x509.ExtKeyUsage, err error) {
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
			err = dbgutil.FormatError("[%s] not supported", arrs[idx])
			return
		}
	}
	err = nil
	return
}

func get_byte_array(valarr []interface{}, note string) (retv []byte, err error) {
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
		err = dbgutil.FormatError("[%s] [%s] parse error %s", note, times, err.Error())
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
		curnet, err = net.ParseCIDR(s)
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
		curnet, err = net.ParseIP(s)
		if err != nil {
			err = dbgutil.FormatError("[%s].[%d] parse [%s] error %s", note, idx, s, err.Error())
			return
		}
		retv = append(retv, curnet)
	}
	err = nil
	return
}

func get_objoids_value(val []interface{}, note string) (retv []asn1.ObjectIdentifier, err error) {
	var idx int
	var curoid asn1.ObjectIdentifier
	var curinter []interface{}
	var nval []interface{}
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

func get_bigint_value(val string, note string) (retv big.Int, err error) {
	var base int = 10
	var inputs string = val
	if strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
		inputs = val[2:]
		base = 16
	} else if strings.HasPrefix(val, "x") || strings.HasPrefix(val, "X") {
		inputs = val[1:]
		base = 16
	}
	retv, ok = big.SetString(inputs, base)
	if !ok {
		err = dbgutil.FormatError("[%s] [%s] not valid big.Int ", note, val)
		return
	}
	err = nil
	return
}

func get_algorithm_value(val string, note string) (retv int, err error) {
	retv = x509.UnknownSignatureAlgorithm
	if 
}

func get_certificate_file(f string) (tempx509 x509.Certificate, err error) {
	var s string
	var ok bool
	var mapv map[string]interface{}
	var arrs []string
	var valarr []interface{}
	var valbool bool
	var intval interface{}
	tempx509 = &x509.Certificate{}
	s, err = fileop.ReadFile(f)
	if err != nil {
		return
	}
	mapv, err = jsonext.GetJsonMap(s)
	if err != nil {
		return
	}

	tempx509.ExtKeyUsage = []x509.ExtKeyUsage{}
	arrs, ok = mapv[KEYWORD_EXTKEYUSAGE].([]string)
	if ok {
		logutil.Debug("[%s]=%v", KEYWORD_EXTKEYUSAGE, arrs)
		tempx509.ExtKeyUsage, err = get_key_ext_usage(arrs)
		if err != nil {
			return
		}
	}

	tempx509.AuthorityKeyId = []byte{}

	valarr, ok = mapv[KEYWORD_AUTHORITY_KEY_ID].([]interface{})
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_AUTHORITY_KEY_ID)
		tempx509.AuthorityKeyId, err = get_byte_array(valarr, KEYWORD_AUTHORITY_KEY_ID)
		if err != nil {
			return
		}
	}

	tempx509.BasicConstraintsValid = false
	valbool, ok = mapv[KEYWORD_BASIC_CONSTRAINT_VALID].(bool)
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_BASIC_CONSTRAINT_VALID)
		tempx509.BasicConstraintsValid = valbool
	}

	tempx509.CRLDistributionPoints = []string{}
	arrs, ok = mapv[KEYWORD_CRL_DISTRIBUTION_POINTS].([]string)
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_CRL_DISTRIBUTION_POINTS)
		tempx509.CRLDistributionPoints = arrs
	}

	tempx509.DNSNames = []string{}
	arrs, ok = mapv[KEYWORD_DNS_NAMES].([]string)
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_DNS_NAMES)
		tempx509.DNSNames = arrs
	}

	tempx509.EmailAddresses = []string{}
	arrs, ok = mapv[KEYWORD_EMAIL_ADDRESSES].([]string)
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_EMAIL_ADDRESSES)
		tempx509.EmailAddresses = arrs
	}

	tempx509.ExcludedDNSDomains = []string{}
	arrs, ok = mapv[KEYWORD_EXCLUDED_DNS_DOMAINS].([]string)
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_EXCLUDED_DNS_DOMAINS)
		tempx509.ExcludedDNSDomains = arrs
	}

	tempx509.ExcludedEmailAddresses = []string{}
	arrs, ok = mapv[KEYWORD_EXCLUDED_EMAILADDRESSES].([]string)
	if ok {
		logutil.Debug("[%s] parse")
	}

	err = nil
	return

}
