package main

import (
	"dbgutil"
	"crypto/x509"
	"crypto/rsa"
	"reflect"
	"crypto/sha256"
	"crypto/sha1"
	"crypto/md5"
	"crypto/sha512"
	"hash"
	"crypto"
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

func get_crypto_hash(origdata []byte, hashtype string) (hashed []byte, hashalgo crypto.Hash, err error) {
	var hasher hash.Hash
	if hashtype == "sha256" {
		hasher = sha256.New()
		hashalgo = crypto.SHA256
	} else if hashtype == "sha1" {
		hasher = sha1.New()
		hashalgo = crypto.SHA1
	} else if hashtype == "md5" {
		hasher = md5.New()
		hashalgo = crypto.MD5
	} else if hashtype == "sha224" {
		hasher = sha256.New224()
		hashalgo = crypto.SHA224
	} else if hashtype == "sha512" {
		hasher = sha512.New()
		hashalgo = crypto.SHA512
	} else if hashtype == "sha384" {
		hasher = sha512.New384()
		hashalgo = crypto.SHA384
	} else {
		err = dbgutil.FormatError("not support hashtype {}",hashtype)
		return
	}
	hasher.Write(origdata)
	hashed = hasher.Sum(nil)
	err = nil
	return
}