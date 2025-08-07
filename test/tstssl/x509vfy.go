package main

import (
	"crypto/x509"
	"dbgutil"
	"fileop"
	"jsonext"
	"logutil"
)

func load_certs_value(cert *x509.CertPool, arrs []string, note string) (err error) {
	var i int
	var s string
	var pemb []byte
	var ok bool
	for i, s = range arrs {
		pemb, err = fileop.ReadFileBytes(s)
		if err != nil {
			return
		}
		ok = cert.AppendCertsFromPEM(pemb)
		if !ok {
			err = dbgutil.FormatError("[%s].[%d] [%s] load cert error", note, i, s)
			return
		}
	}
	err = nil
	return
}

func get_x509_verify_options(cfgfile string) (retv x509.VerifyOptions, err error) {
	var s string
	var vmap map[string]interface{}
	var vals string
	var valarr []interface{}
	var arrs []string
	var ok bool
	var intval interface{}
	s, err = fileop.ReadFile(cfgfile)
	if err != nil {
		return
	}
	vmap, err = jsonext.GetJsonMap(s)
	if err != nil {
		return
	}

	retv = x509.VerifyOptions{}
	retv.Intermediates = x509.NewCertPool()
	retv.Roots = x509.NewCertPool()

	vals, ok = vmap[KEYWORD_DNSNAME].(string)
	if ok {
		logutil.Debug("[%s] parsed", KEYWORD_DNSNAME)
		retv.DNSName = vals
	}

	valarr, ok = vmap[KEYWORD_ROOTS].([]interface{})
	if ok {
		logutil.Debug("[%s] parsed", KEYWORD_ROOTS)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_ROOTS)
		if err != nil {
			return
		}
		err = load_certs_value(retv.Roots, arrs, KEYWORD_ROOTS)
		if err != nil {
			return
		}
	}

	valarr, ok = vmap[KEYWORD_INTERMEDIATES].([]interface{})
	if ok {
		logutil.Debug("[%s] parsed", KEYWORD_INTERMEDIATES)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_INTERMEDIATES)
		if err != nil {
			return
		}
		err = load_certs_value(retv.Intermediates, arrs, KEYWORD_INTERMEDIATES)
		if err != nil {
			return
		}
	}

	s, ok = vmap[KEYWORD_CURRENTTIME].(string)
	if ok {
		logutil.Debug("[%s] parsed", KEYWORD_CURRENTTIME)
		retv.CurrentTime, err = get_time_value(s, KEYWORD_CURRENTTIME)
		if err != nil {
			return
		}
	}

	valarr, ok = vmap[KEYWORD_EXT_KEY_USAGE].([]interface{})
	if ok {
		logutil.Debug("[%s] parsed", KEYWORD_EXT_KEY_USAGE)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_EXT_KEY_USAGE)
		if err != nil {
			return
		}
		retv.KeyUsages, err = get_key_ext_usage(arrs, KEYWORD_EXT_KEY_USAGE)
		if err != nil {
			return
		}
	}

	intval, ok = vmap[KEYWORD_MAX_CONSTRAINTS_COMPARISIONS].(interface{})
	if ok {
		logutil.Debug("[%s] parsed", KEYWORD_MAX_CONSTRAINTS_COMPARISIONS)
		retv.MaxConstraintComparisions, err = get_int_value(intval, KEYWORD_MAX_CONSTRAINTS_COMPARISIONS)
		if err != nil {
			return
		}
	}

	err = nil
	return
}
