package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"dbgutil"
	"encoding/asn1"
	"encoding/pem"
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"logutil"
	"reflect"
)

func X509create_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var rsakey *rsa.PrivateKey
	var tempx509 x509.Certificate
	var pubkey *rsa.PublicKey
	var capem *bytes.Buffer
	var cabytes []byte
	var cafile string
	var keyfile string
	var rsabytes []byte
	var pkany any
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
		err = dbgutil.FormatError("to get template json file for Certifacate")
		return
	}

	keyfile = ns.GetString("keyfile")
	if len(keyfile) == 0 {
		err = dbgutil.FormatError("need keyfile")
		return
	}

	rsabytes, err = read_pem_or_der(keyfile)
	if err != nil {
		return
	}

	pkany, err = x509.ParsePKCS8PrivateKey(rsabytes)
	if err != nil {
		err = dbgutil.FormatError("[%s] not valid rsa %s", keyfile, err.Error())
		return
	}
	switch pkany.(type) {
	case *rsa.PrivateKey:
		rsakey = pkany.(*rsa.PrivateKey)
	default:
		err = dbgutil.FormatError("key is not rsakey type [%s]", reflect.TypeOf(pkany))
		return
	}

	tempx509, err = get_certificate_file(sarr[0])
	if err != nil {
		return
	}

	pubkey = &(rsakey.PublicKey)
	//cabytes, err = x509.CreateCertificate(rand.Reader, &tempx509, &tempx509, pubkey, rsakey)
	cabytes, err = createCertificate(rand.Reader, &tempx509, &tempx509, pubkey, rsakey)
	if err != nil {
		err = dbgutil.FormatError("output certificate %s", err.Error())
		return
	}
	capem = new(bytes.Buffer)
	err = pem.Encode(capem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cabytes,
	})
	if err != nil {
		err = dbgutil.FormatError("encode CERTIFICATE error %s", err.Error())
		return
	}
	cafile = ns.GetString("certfile")
	_, err = fileop.WriteFileBytes(cafile, capem.Bytes())
	if err != nil {
		err = dbgutil.FormatError("can not write certfile [%s] %s", cafile, err.Error())
		return
	}

	err = nil
	return
}

func X509parse_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var tempx509 *x509.Certificate
	var x509bytes []byte
	var f string
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
		err = dbgutil.FormatError("to get template json file for Certifacate")
		return
	}

	for _, f = range sarr {
		x509bytes, err = read_pem_or_der(f)
		if err != nil {
			return
		}
		tempx509, err = parseCertificate(x509bytes)
		if err != nil {
			return
		}
		display_x509_cert(tempx509)
	}

	err = nil
	return
}

func X509vfy_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var tempx509 *x509.Certificate
	var x509bytes []byte
	var vfyopt VerifyOptionsF
	var f string
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
		err = dbgutil.FormatError("to verify ")
		return
	}

	f = ns.GetString("vfyopt")
	if len(f) == 0 {
		err = dbgutil.FormatError("need vfyopt set")
		return
	}
	vfyopt, err = get_x509_verify_options(f)
	if err != nil {
		return
	}

	for _, f = range sarr {
		x509bytes, err = read_pem_or_der(f)
		if err != nil {
			return
		}
		tempx509, err = x509.ParseCertificate(x509bytes)
		if err != nil {
			return
		}

		_, err = Verify_Certificate(tempx509, vfyopt)
		//_, err = tempx509.Verify(vfyopt)
		if err != nil {
			return
		}

		fmt.Printf("%s verified\n", f)
	}

	err = nil
	return
}

func X509reqvfy_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var x509bytes []byte
	var req *x509.CertificateRequest
	var f string
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
		err = dbgutil.FormatError("to verify ")
		return
	}

	for _, f = range sarr {
		x509bytes, err = read_pem_or_der(f)
		if err != nil {
			return
		}
		req, err = x509.ParseCertificateRequest(x509bytes)
		if err != nil {
			return
		}

		err = x509req_CheckSignature(req)
		if err != nil {
			err = dbgutil.FormatError("not check [%s] valid", f)
			return
		}

		fmt.Printf("%s verified\n", f)
	}

	err = nil
	return
}

func dump_x509_req(req *x509.CertificateRequest) {
	display_buffer(req.Raw, "Raw", 0)
	display_buffer(req.RawTBSCertificateRequest, "RawTBSCertificateRequest", 0)
	display_buffer(req.RawSubjectPublicKeyInfo, "RawSubjectPublicKeyInfo", 0)
	display_buffer(req.RawSubject, "RawSubject", 0)
	fmt.Println("Version %d", req.Version)
	display_buffer(req.Signature, "Signature", 0)
	fmt.Println("SignatureAlgorithm %d", req.SignatureAlgorithm)
	fmt.Println("PublicKeyAlgorithm %d", req.PublicKeyAlgorithm)
	debug_pkix_name(&req.Subject, "Subject", 0)
	display_attri_array(req.Attributes, "Attributes", 0)
	display_extensions(req.Extensions, "Extensions", 0)
	display_extensions(req.ExtraExtensions, "ExtraExtensions", 0)
	display_strings(req.DNSNames, "DNSNames", 0)
	display_strings(req.EmailAddresses, "EmailAddresses", 0)
	display_ips(req.IPAddresses, "IPAddresses", 0)
	display_urls(req.URIs, "URIs", 0)
}

func X509reqparse_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var req *x509.CertificateRequest
	var x509bytes []byte
	var f string
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
		err = dbgutil.FormatError("to verify ")
		return
	}

	for _, f = range sarr {
		x509bytes, err = read_pem_or_der(f)
		if err != nil {
			return
		}
		req, err = Data_ParseCertificateRequest(x509bytes)
		if err != nil {
			return
		}
		dump_x509_req(req)
	}

	err = nil
	return
}

func Bitstrdec_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var bitstr asn1.BitString
	var x509bytes []byte
	var f string
	var nbytes []byte
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
		err = dbgutil.FormatError("to verify ")
		return
	}

	for _, f = range sarr {
		x509bytes, err = read_pem_or_der(f)
		if err != nil {
			return
		}

		_, err = asn1.Unmarshal(x509bytes, &bitstr)
		if err != nil {
			err = dbgutil.FormatError("decode [%s] error %s", f, err.Error())
			return
		}
		logutil.DebugBuffer(bitstr.Bytes, "bytes BitLength %d", bitstr.BitLength)
		nbytes = bitstr.RightAlign()
		logutil.DebugBuffer(nbytes, "nbytes RightAlign")
	}

	err = nil
	return
}

func Reqcreate_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var outfile string
	var nbytes []byte
	var keyfile string
	var jsonfile string
	var privkey *rsa.PrivateKey
	var certreq *x509.CertificateRequest
	var capem *bytes.Buffer
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
		err = dbgutil.FormatError("need json file")
		return
	}

	keyfile = ns.GetString("keyfile")
	privkey, err = get_rsa_private(keyfile)
	if err != nil {
		return
	}

	jsonfile = sarr[0]
	certreq, err = parse_x509_req_json(jsonfile)
	if err != nil {
		return
	}

	nbytes, err = x509.CreateCertificateRequest(rand.Reader, certreq, privkey)
	if err != nil {
		err = dbgutil.FormatError("can not create CertificateRequest error %s", err.Error())
		return
	}

	capem = new(bytes.Buffer)
	block := &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: nbytes,
	}
	err = pem.Encode(capem, block)
	outfile = ns.GetString("output")
	_, err = fileop.WriteFileBytes(outfile, capem.Bytes())
	if err != nil {
		return
	}
	err = nil
	return
}

func init() {
	X509create_handler(nil, nil, nil)
	X509parse_handler(nil, nil, nil)
	X509vfy_handler(nil, nil, nil)
	X509reqvfy_handler(nil, nil, nil)
	X509reqparse_handler(nil, nil, nil)
	Bitstrdec_handler(nil, nil, nil)
	Reqcreate_handler(nil, nil, nil)
}

func load_x509_handler(parser *extargsparse.ExtArgsParse) (err error) {
	var commandline = `{
		"x509create<X509create_handler>##jsonfile to set x509 from template file by keyfile##":{
			"$":1
		},
		"x509parse<X509parse_handler>##pemfile ... to parse x509.Certificate##":{
			"$":"+"
		},
		"x509vfy<X509vfy_handler>##pemfile ... to parse in vfyopt file##":{
			"$":"+"
		},
		"x509reqvfy<X509reqvfy_handler>##pemfile ... to verify CertificateRequest##":{
			"$" : "+"
		},
		"x509reqparse<X509reqparse_handler>##pemfile ... to decode CertificateRequest##" : {
			"$" : "+"
		},
		"bitstrdec<Bitstrdec_handler>##binfile ... to decode asn1.BitString##" : {
			"$" : "+"
		},
		"reqcreate<Reqcreate_handler>##req.json to create CertificateRequest with ##" : {
			"$" : 1
		}
	}`
	err = parser.LoadCommandLineString(commandline)
	return
}
