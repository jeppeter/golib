package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"dbgutil"
	"encoding/asn1"
	"encoding/pem"
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"jsonext"
	"logutil"
	"math/big"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Genkeycert_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var rsakey *rsa.PrivateKey
	var rsabits int = 2048
	var tempx509 *x509.Certificate
	var pubkey *rsa.PublicKey
	var capem *bytes.Buffer
	var cabytes []byte
	var cafile string
	var keyfile string
	var rsabytes []byte
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
	rsakey = rsakey

	tempx509 = &x509.Certificate{
		ExtKeyUsage:                 []x509.ExtKeyUsage{},
		AuthorityKeyId:              []byte{},
		BasicConstraintsValid:       false,
		CRLDistributionPoints:       []string{},
		DNSNames:                    []string{},
		EmailAddresses:              []string{},
		ExcludedDNSDomains:          []string{},
		ExcludedEmailAddresses:      []string{},
		ExcludedIPRanges:            []*net.IPNet{},
		IPAddresses:                 []net.IP{},
		IsCA:                        false,
		IssuingCertificateURL:       []string{},
		KeyUsage:                    0,
		MaxPathLen:                  0,
		MaxPathLenZero:              false,
		NotAfter:                    time.Now().AddDate(20, 0, 0),
		NotBefore:                   time.Now(),
		OCSPServer:                  []string{},
		PermittedDNSDomains:         []string{},
		PermittedDNSDomainsCritical: false,
		PermittedEmailAddresses:     []string{},
		PermittedIPRanges:           []*net.IPNet{},
		PermittedURIDomains:         []string{},
		PolicyIdentifiers:           []asn1.ObjectIdentifier{},
		Policies:                    []x509.OID{},
		SerialNumber:                big.NewInt(0),
		SignatureAlgorithm:          x509.SHA256WithRSA,
		Subject: pkix.Name{
			CommonName: "NVIDIA GameStream Client",
			Names: []pkix.AttributeTypeAndValue{
				pkix.AttributeTypeAndValue{
					Type:  []int{2, 5, 4, 3},
					Value: "NVIDIA GameStream Client",
				},
			},
		},
		SubjectKeyId:       []byte{},
		URIs:               []*url.URL{},
		UnknownExtKeyUsage: []asn1.ObjectIdentifier{},
	}

	pubkey = &(rsakey.PublicKey)
	cabytes, err = x509.CreateCertificate(rand.Reader, tempx509, tempx509, pubkey, rsakey)
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

	capem = new(bytes.Buffer)
	rsabytes, err = x509.MarshalPKCS8PrivateKey(rsakey)
	if err != nil {
		err = dbgutil.FormatError("MarshalPKCS8PrivateKey error %s", err.Error())
		return
	}

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

func Pemtoder_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var input string
	var output string
	var inbytes []byte
	var block *pem.Block
	err = nil
	if ns == nil {
		return nil
	}
	err = logutil.InitLog(ns)
	if err != nil {
		logutil.Error("can not Initlog err[%s]", err.Error())
		return err
	}

	input = ns.GetString("input")
	output = ns.GetString("output")

	inbytes, err = fileop.ReadFileBytes(input)
	if err != nil {
		return
	}
	block, _ = pem.Decode(inbytes)
	if block == nil {
		err = dbgutil.FormatError("can not decode [%s]", input)
		return
	}

	_, err = fileop.WriteFileBytes(output, block.Bytes)
	if err != nil {
		return
	}
	return
}

func Dertopem_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var input string
	var output string
	var inbytes []byte
	var sarr []string
	var capem *bytes.Buffer
	var pemtag string
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

	pemtag = sarr[0]
	input = ns.GetString("input")
	output = ns.GetString("output")

	inbytes, err = fileop.ReadFileBytes(input)
	if err != nil {
		return
	}

	capem = new(bytes.Buffer)
	err = pem.Encode(capem, &pem.Block{
		Type:  strings.ToUpper(pemtag),
		Bytes: inbytes,
	})

	_, err = fileop.WriteFileBytes(output, capem.Bytes())
	if err != nil {
		return
	}
	return
}

func get_pkix_name(f string) (name pkix.Name, err error) {
	var mapv map[string]interface{}
	var s string
	s, err = fileop.ReadFile(f)
	if err != nil {
		return
	}

	logutil.Debug("[%s]\n%s", f, s)
	mapv, err = jsonext.GetJsonMap(s)
	if err != nil {
		return
	}
	name, err = get_pkixname_value(mapv, "")
	return
}

func Pkixname_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var curname pkix.Name
	var sarr []string
	var f string
	var outb []byte

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
	for _, f = range sarr {
		curname, err = get_pkix_name(f)
		if err != nil {
			return
		}

		logutil.Debug("[%s] Country %v", f, curname.Country)

		outb, err = asn1.Marshal(curname.ToRDNSequence())
		if err != nil {
			return
		}

		logutil.DebugBuffer(outb, "[%s] bytes", f)
	}

	err = nil
	return
}

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

func init() {
	Genkeycert_handler(nil, nil, nil)
	Pemtoder_handler(nil, nil, nil)
	Dertopem_handler(nil, nil, nil)
	Pkixname_handler(nil, nil, nil)
	X509create_handler(nil, nil, nil)
	X509parse_handler(nil, nil, nil)
	X509vfy_handler(nil, nil, nil)
}
func main() {
	var commandline string
	var err error
	var parser *extargsparse.ExtArgsParse
	var ns *extargsparse.NameSpaceEx

	commandline = `{
		"input|i" : null,
		"output|o" : null,
		"certfile" : null,
		"keyfile" : null,
		"vfyopt" : null,
		"digesttype" : "sha256",
		"genkeycert<Genkeycert_handler>##[rsabits] to generate keyfile and certfile default rsabits 2048##" : {
			"$" : "?"
		},
		"pemtoder<Pemtoder_handler>##from input to output with pem to der##" : {
			"$" : 0
		},
		"dertopem<Dertopem_handler>##TAG from input to output to writeout##" : {
			"$" : 1
		},
		"pkixname<Pkixname_handler>##[inputfile] ... to format pkix.Name asn1.Marshal inputfile is json file##" : {
			"$" : 1
		},
		"x509create<X509create_handler>##jsonfile to set x509 from template file by keyfile##" : {
			"$" : 1
		},
		"x509parse<X509parse_handler>##pemfile ... to parse x509.Certificate##" : {
			"$" : "+"
		},
		"x509vfy<X509vfy_handler>##pemfile ... to parse in vfyopt file##" : {
			"$" : "+"
		}

	}`

	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not make parser err[%s]\n", err.Error())
		atexit.Exit(5)
	}

	err = logutil.PrepareLog(parser)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not set [%s]\n", err.Error())
		atexit.Exit(5)
	}
	err = parser.LoadCommandLineString(commandline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not parse %s\n", commandline)
		atexit.Exit(5)
	}

	err = load_rsa_command(parser)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		atexit.Exit(5)
	}

	ns, err = parser.ParseCommandLineEx(nil, nil, nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not use parse command line [%s]\n", err.Error())
		atexit.Exit(4)
	}
	if len(ns.GetString("subcommand")) == 0 {
		fmt.Fprintf(os.Stderr, "can not get subcommand\n")
		atexit.Exit(5)
	}
	//fmt.Fprintf(os.Stdout, "subcommand [%s] succ\n", ns.GetString("subcommand"))
	atexit.Exit(0)
	return
}
