package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"dbgutil"
	"encoding/pem"
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"logutil"
	"reflect"
	"strconv"
)

func rsa_sign_sha256(keydata []byte, indata []byte) (signdata []byte, err error) {
	var hashed []byte
	var pkany any
	var rsakey *rsa.PrivateKey

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

	hasher := sha256.New()
	hasher.Write(indata)

	hashed = hasher.Sum(nil)
	signdata, err = rsa.SignPKCS1v15(nil, rsakey, crypto.SHA256, hashed)
	if err != nil {
		return
	}
	return
}

func Rsasign_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var input string
	var sarr []string
	var keyfile string
	var signfile string
	var signdata []byte
	var indata []byte
	var keydata []byte
	err = nil
	if ns == nil {
		return nil
	}
	err = logutil.InitLog(ns)
	if err != nil {
		logutil.Error("can not Initlog err[%s]", err.Error())
		return err
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 3 {
		err = dbgutil.FormatError("need keyfile inputfile signfile")
		return
	}

	keyfile = sarr[0]
	input = sarr[1]
	signfile = sarr[2]

	keydata, err = read_pem_or_der(keyfile)
	if err != nil {
		return
	}

	if err != nil {
		return
	}

	indata, err = fileop.ReadFileBytes(input)
	if err != nil {
		return
	}

	signdata, err = rsa_sign_sha256(keydata, indata)
	if err != nil {
		return
	}

	_, err = fileop.WriteFileBytes(signfile, signdata)
	if err != nil {
		return
	}

	return
}

func rsa_pss_sign_sha256(keydata []byte, indata []byte, psslen int) (signdata []byte, err error) {
	var hashed []byte
	var pkany any
	var rsakey *rsa.PrivateKey
	var pssopt *rsa.PSSOptions

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

	hasher := sha256.New()
	hasher.Write(indata)

	hashed = hasher.Sum(nil)

	pssopt = &rsa.PSSOptions{}
	pssopt.SaltLength = psslen
	pssopt.Hash = crypto.SHA256

	signdata, err = rsa.SignPSS(rand.Reader, rsakey, crypto.SHA256, hashed, pssopt)
	if err != nil {
		return
	}
	return
}

func Rsapsssign_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var input string
	var sarr []string
	var keyfile string
	var signfile string
	var signdata []byte
	var indata []byte
	var keydata []byte
	var psslen int
	err = nil
	if ns == nil {
		return nil
	}
	err = logutil.InitLog(ns)
	if err != nil {
		logutil.Error("can not Initlog err[%s]", err.Error())
		return err
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 3 {
		err = dbgutil.FormatError("need keyfile inputfile signfile")
		return
	}

	psslen = ns.GetInt("psslength")

	keyfile = sarr[0]
	input = sarr[1]
	signfile = sarr[2]

	keydata, err = read_pem_or_der(keyfile)
	if err != nil {
		return
	}

	if err != nil {
		return
	}

	indata, err = fileop.ReadFileBytes(input)
	if err != nil {
		return
	}

	signdata, err = rsa_pss_sign_sha256(keydata, indata, psslen)
	if err != nil {
		return
	}

	_, err = fileop.WriteFileBytes(signfile, signdata)
	if err != nil {
		return
	}

	return
}

func rsa_verify_sha256(keydata []byte, indata []byte, signdata []byte) (err error) {
	var hashed []byte
	var pubkey *rsa.PublicKey
	var pkany any
	var rsakey *rsa.PrivateKey

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

	pubkey = &rsakey.PublicKey

	hasher := sha256.New()
	hasher.Write(indata)

	hashed = hasher.Sum(nil)

	err = rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, hashed, signdata)
	if err != nil {
		err = dbgutil.FormatError("verify error %s", err.Error())
		return
	}
	return
}

func Rsavfy_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var input string
	var indata, signdata []byte
	var sarr []string
	var keyfile string
	var signfile string
	var keydata []byte

	err = nil
	if ns == nil {
		return nil
	}
	err = logutil.InitLog(ns)
	if err != nil {
		logutil.Error("can not Initlog err[%s]", err.Error())
		return err
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 1 {
		err = dbgutil.FormatError("need TAG")
		return
	}

	keyfile = sarr[0]
	input = sarr[1]
	signfile = sarr[2]

	keydata, err = read_pem_or_der(keyfile)
	if err != nil {
		return
	}

	indata, err = fileop.ReadFileBytes(input)
	if err != nil {
		return
	}

	signdata, err = fileop.ReadFileBytes(signfile)
	if err != nil {
		return
	}

	err = rsa_verify_sha256(keydata, indata, signdata)
	if err != nil {
		return
	}

	fmt.Printf("verify [%s] [%s] [%s] succ\n", keyfile, input, signfile)
	return
}

func rsa_pss_verify_sha256(keydata []byte, indata []byte, signdata []byte, psslen int) (err error) {
	var hashed []byte
	var pubkey *rsa.PublicKey
	var pkany any
	var rsakey *rsa.PrivateKey
	var pssopt *rsa.PSSOptions

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

	pubkey = &rsakey.PublicKey

	hasher := sha256.New()
	hasher.Write(indata)

	hashed = hasher.Sum(nil)

	pssopt = &rsa.PSSOptions{}
	pssopt.SaltLength = psslen
	pssopt.Hash = crypto.SHA256

	err = rsa.VerifyPSS(pubkey, crypto.SHA256, hashed, signdata, pssopt)
	if err != nil {
		err = dbgutil.FormatError("verify error %s", err.Error())
		return
	}
	return
}

func Rsapssvfy_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var input string
	var indata, signdata []byte
	var sarr []string
	var keyfile string
	var signfile string
	var keydata []byte
	var psslen int

	err = nil
	if ns == nil {
		return nil
	}
	err = logutil.InitLog(ns)
	if err != nil {
		logutil.Error("can not Initlog err[%s]", err.Error())
		return err
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 1 {
		err = dbgutil.FormatError("need TAG")
		return
	}
	psslen = ns.GetInt("psslength")

	keyfile = sarr[0]
	input = sarr[1]
	signfile = sarr[2]

	keydata, err = read_pem_or_der(keyfile)
	if err != nil {
		return
	}

	indata, err = fileop.ReadFileBytes(input)
	if err != nil {
		return
	}

	signdata, err = fileop.ReadFileBytes(signfile)
	if err != nil {
		return
	}

	err = rsa_pss_verify_sha256(keydata, indata, signdata, psslen)
	if err != nil {
		return
	}

	fmt.Printf("verify [%s] [%s] [%s] psslen %d succ\n", keyfile, input, signfile, psslen)
	return
}

func Rsagen_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var rsabits int = 2048
	var rsakey *rsa.PrivateKey
	var rsabytes []byte
	var capem *bytes.Buffer
	var keyfile string
	err = nil
	if ns == nil {
		return nil
	}
	err = logutil.InitLog(ns)
	if err != nil {
		logutil.Error("can not Initlog err[%s]", err.Error())
		return err
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) > 0 {
		rsabits, err = strconv.Atoi(sarr[0])
		if err != nil {
			err = dbgutil.FormatError("[%s] not valid bits", sarr[0])
			return
		}
	}

	rsakey, err = rsa.GenerateKey(rand.Reader, rsabits)
	if err != nil {
		err = dbgutil.FormatError("generate %d error %s", rsabits, err.Error())
		return
	}

	rsabytes, err = x509.MarshalPKCS8PrivateKey(rsakey)
	if err != nil {
		err = dbgutil.FormatError("MarshalPKCS8PrivateKey error %s", err.Error())
		return
	}

	capem = new(bytes.Buffer)

	err = pem.Encode(capem, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: rsabytes,
	})
	if err != nil {
		err = dbgutil.FormatError("encode RSA PRIVATE KEY error %s", err.Error())
		return
	}
	keyfile = ns.GetString("keyfile")
	_, err = fileop.WriteFileBytes(keyfile, capem.Bytes())
	if err != nil {
		err = dbgutil.FormatError("can not write keyfile [%s] %s", keyfile, err.Error())
		return
	}
	err = nil
	return
}

func init() {
	Rsasign_handler(nil, nil, nil)
	Rsavfy_handler(nil, nil, nil)
	Rsapsssign_handler(nil, nil, nil)
	Rsapssvfy_handler(nil, nil, nil)
	Rsagen_handler(nil, nil, nil)
}

func load_rsa_command(parser *extargsparse.ExtArgsParse) (err error) {

	var commandline string = `{
		"psslength##pss length set 0 for auto -1 for equal hash##" : 0,
		"rsagen<Rsagen_handler>##bits to generate rsa key file##" : {
			"$" : 1
		},
		"rsasign<Rsasign_handler>##privkeyfile inputfile signfile to sign data##" : {
			"$": 3
		},
		"rsavfy<Rsavfy_handler>##privkeyfile inputfile signfile to verify data##" : {
			"$" : 3
		},
		"rsapsssign<Rsapsssign_handler>##privkeyfile inputdata signfile to sign data##" : {
			"$" : 3
		},
		"rsapssvfy<Rsapssvfy_handler>##privkeyfile inputfile signfile to verify data##" : {
			"$" : 3
		}

	}`
	err = parser.LoadCommandLineString(commandline)
	return
}
