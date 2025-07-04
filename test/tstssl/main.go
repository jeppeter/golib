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
	"logutil"
	"math/big"
	"net"
	"net/url"
	"os"
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

func init() {
	Genkeycert_handler(nil, nil, nil)
	Pemtoder_handler(nil, nil, nil)
	Dertopem_handler(nil, nil, nil)
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
		"genkeycert<Genkeycert_handler>##[rsabits] to generate keyfile and certfile default rsabits 2048##" : {
			"$" : "?"
		},
		"pemtoder<Pemtoder_handler>##from input to output with pem to der##" : {
			"$" : 0
		},
		"dertopem<Dertopem_handler>##TAG from input to output to writeout##" : {
			"$" : 1
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
