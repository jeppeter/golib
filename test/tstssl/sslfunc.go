package main

import (
	"dbgutil"
	"crypto/x509"
	"crypto/rsa"
	"reflect"
)

func get_rsa_private(keyfile string) (rsakey *rsa.PrivateKey,err error) {
	var pkany any
	var keydata []byte

	keydata , err = read_pem_or_der(keyfile)
	if err != nil {
		return
	}

	pkany, err = x509.ParsePKCS8PrivateKey(keydata)
	if err != nil {
		err = dbgutil.FormatError("keydata not valid rsa %s", err.Error())
		return
	}
	switch pkany.(type) {
	case *rsa.PrivateKey:
		rsakey = pkany.(*rsa.PrivateKey)
	default:
		err = dbgutil.FormatError("key is not rsakey type [%s]", reflect.TypeOf(pkany))
		return
	}

	return

}