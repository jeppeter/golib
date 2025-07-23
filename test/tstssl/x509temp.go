package main

import (
	"crypto/x509"
	"dbgutil"
	"fileop"
	"jsonext"
	"logutil"
)

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

	KEYWORD_AUTHORITY_KEY_ID        = "authoritykeyid"
	KEYWORD_BASIC_CONSTRAINT_VALID  = "basiccontraintsvalid"
	KEYWORD_CRL_DISTRIBUTION_POINTS = "crldistributionpoints"
	KEYWORD_DNS_NAMES               = "dnsnames"
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
		}

		retv = append(retv, byte(vali))
	}

	err = nil
	return

}

func get_certificate_file(f string) (tempx509 *x509.Certificate, err error) {
	var s string
	var ok bool
	var mapv map[string]interface{}
	var arrs []string
	var valarr []interface{}
	var valbool bool
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

	arrs, ok = mapv[KEYWORD_CRL_DISTRIBUTION_POINTS].([]string)
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_CRL_DISTRIBUTION_POINTS)
		tempx509.CRLDistributionPoints = arrs
	}

	arrs, ok = mapv[KEYWORD_DNS_NAMES].([]string)
	if ok {
		logutil.Debug("[%s] parse", KEYWORD_DNS_NAMES)
		tempx509.DNSNames = arrs
	}

	err = nil
	return

}
