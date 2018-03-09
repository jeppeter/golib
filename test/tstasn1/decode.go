package main

import (
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	//"golang.org/x/crypto/pkcs12"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"time"
)

func DecodePem(infile string) (ps []*pem.Block, err error) {
	var data []byte
	var p *pem.Block

	ps = make([]*pem.Block, 0)
	p = nil
	data, err = ioutil.ReadFile(infile)
	if err != nil {
		return
	}

	for len(data) > 0 {
		p, data = pem.Decode(data)
		ps = append(ps, p)
	}
	err = nil
	return
}

func DecodeDer(infile string) (data []byte, err error) {
	data, err = ioutil.ReadFile(infile)
	return
}

type Asn1Seq struct {
	Value   *asn1.RawValue
	Child   []*Asn1Seq
	Verbose int
}

func NewAsn1Seq(v *asn1.RawValue) *Asn1Seq {
	var self *Asn1Seq
	self = &Asn1Seq{}
	self.Value = v
	self.Child = make([]*Asn1Seq, 0)
	self.Verbose = 0
	return self
}

func formatTabs(tabs int, fmtstr string, a ...interface{}) string {
	var s string
	var i int
	for i = 0; i < tabs; i++ {
		s += "    "
	}
	s += fmt.Sprintf(fmtstr, a...)
	s += "\n"
	return s
}

func FormatBytes(data []byte, verbose int) string {
	var s string = ""
	var b byte
	var i int
	s += "["
	if len(data) < 16 || verbose >= 3 {
		for i, b = range data {
			if i > 0 {
				s += ","
			}
			s += fmt.Sprintf("0x%02x", b)
		}
	} else {
		i = 0
		for i < 4 {
			if i > 0 {
				s += ","
			}
			s += fmt.Sprintf("0x%02x", data[i])
			i++
		}
		s += "..."
		i = len(data) - 4
		for i < len(data) {
			if i > (len(data) - 4) {
				s += ","
			}
			s += fmt.Sprintf("0x%02x", data[i])
			i++
		}
	}
	s += "]"
	return s
}

type bigint struct {
	B *big.Int
}

func (self *Asn1Seq) formatIntValue() string {
	var s string = ""
	var err error
	if len(self.Value.Bytes) <= 4 {
		var i int
		_, err = asn1.Unmarshal(self.Value.FullBytes, &i)
		if err == nil {
			s += fmt.Sprintf("Integer[0x%x]", i)
		} else {
			fmt.Fprintf(os.Stderr, "err[%s]\n", err.Error())
		}
	} else if len(self.Value.Bytes) <= 8 {
		var i int64
		_, err = asn1.Unmarshal(self.Value.FullBytes, &i)
		if err == nil {
			s += fmt.Sprintf("Integer[0x%x]", i)
		} else {
			fmt.Fprintf(os.Stderr, "err[%s]\n", err.Error())
		}
	} else {
		var bint *bigint
		var newb []byte
		var needlen int
		var blen int
		var i int
		var j int
		var curs string

		blen = len(self.Value.FullBytes)
		needlen = blen
		if blen <= 0x7f {
			/*for tag and length*/
			needlen += 2
		} else if blen <= 0xffff {
			needlen += 4
		} else if blen <= 0xffffff {
			needlen += 5
		} else {
			needlen += 5
		}
		newb = make([]byte, needlen)
		newb[0] = byte(0x30)
		if blen <= 0x7f {
			newb[1] = byte(blen)
			i = 2
		} else if blen <= 0xffff {
			newb[1] = byte(0x82)
			newb[2] = byte((blen >> 8) & 0xff)
			newb[3] = byte(blen & 0xff)
			i = 4
		} else if blen <= 0xffffff {
			newb[1] = byte(0x83)
			newb[2] = byte((blen >> 16) & 0xff)
			newb[3] = byte((blen >> 8) & 0xff)
			newb[4] = byte(blen & 0xff)
			i = 5
		}
		for j = 0; i < needlen; j++ {
			newb[i] = self.Value.FullBytes[j]
			i++
		}

		bint = new(bigint)

		_, err = asn1.Unmarshal(newb, bint)
		if err == nil {
			if self.Verbose >= 3 {
				s += fmt.Sprintf("Integer[0x%s]", bint.B.Text(16))
			} else {
				curs = bint.B.Text(16)
				if len(curs) > 16 {
					s += fmt.Sprintf("Integer[0x%s...%s]", curs[:8], curs[(len(curs)-8):])
				} else {
					s += fmt.Sprintf("Integer[0x%s]", curs)
				}
			}

		} else {
			fmt.Fprintf(os.Stderr, "err[%s]\n", err.Error())
		}
	}

	return s
}

func (self *Asn1Seq) formatOIDValue() string {
	var s string = ""
	var err error
	var id asn1.ObjectIdentifier
	_, err = asn1.Unmarshal(self.Value.FullBytes, &id)
	if err == nil {
		s += fmt.Sprintf("%s", id.String())
	}

	return s
}

func (self *Asn1Seq) formatPrintableString() string {
	var s string = ""
	s += fmt.Sprintf("%s", string(self.Value.Bytes))
	return s
}

func (self *Asn1Seq) formatUTF8String() string {
	var s string = ""
	s += fmt.Sprintf("%s", string(self.Value.Bytes))
	return s
}

func (self *Asn1Seq) formatUTCTime() string {
	var s string = ""
	var t time.Time
	var err error
	var ds string = string(self.Value.Bytes)
	var year, mon, day int
	var min, hour int
	var month time.Month
	year, _ = strconv.Atoi(ds[:2])
	year += 2000
	mon, _ = strconv.Atoi(ds[2:4])
	switch mon {
	case 1:
		month = time.January
	case 2:
		month = time.February
	case 3:
		month = time.March
	case 4:
		month = time.April
	case 5:
		month = time.May
	case 6:
		month = time.June
	case 7:
		month = time.July
	case 8:
		month = time.August
	case 9:
		month = time.September
	case 10:
		month = time.October
	case 11:
		month = time.November
	case 12:
		month = time.December
	default:
		month = time.January
	}
	day, _ = strconv.Atoi(ds[4:6])
	hour, _ = strconv.Atoi(ds[6:8])
	min, _ = strconv.Atoi(ds[8:10])
	t = time.Date(year, month, day, hour, min, 0, 0, time.UTC)
	if err != nil {
		s += ds
		fmt.Fprintf(os.Stderr, "can not parse [%s] err[%s]\n", ds, err.Error())
	} else {
		s += t.Format(time.RFC1123Z)
	}
	return s
}

func (self *Asn1Seq) formatOctetString() string {
	var s string = ""
	s += FormatBytes(self.Value.Bytes, self.Verbose)
	return s
}

func (self *Asn1Seq) formatValue() string {
	var s string = ""
	switch self.Value.Tag {
	case asn1.TagInteger:
		s += self.formatIntValue()
	case asn1.TagOID:
		s += self.formatOIDValue()
	case asn1.TagPrintableString:
		s += self.formatPrintableString()
	case asn1.TagUTF8String:
		s += self.formatUTF8String()
	case asn1.TagUTCTime:
		s += self.formatUTCTime()
	case asn1.TagOctetString:
		s += self.formatOctetString()
	}
	return s
}

func (self *Asn1Seq) formatClassType() string {
	var s string = ""
	var clsstr string = fmt.Sprintf("0x%x", self.Value.Class)
	var tagstr string = fmt.Sprintf("0x%x", self.Value.Tag)
	switch self.Value.Class {
	case asn1.ClassUniversal:
		clsstr = "Universal"
	case asn1.ClassApplication:
		clsstr = "Application"
	case asn1.ClassPrivate:
		clsstr = "Private"
	case asn1.ClassContextSpecific:
		clsstr = "ContextSpecific"
	}

	switch self.Value.Tag {
	case asn1.TagBoolean:
		tagstr = "Boolean"
	case asn1.TagInteger:
		tagstr = "Integer"
	case asn1.TagBitString:
		tagstr = "BitString"
	case asn1.TagOctetString:
		tagstr = "OctetString"
	case asn1.TagNull:
		tagstr = "Null"
	case asn1.TagOID:
		tagstr = "OID"
	case asn1.TagEnum:
		tagstr = "Enum"
	case asn1.TagUTF8String:
		tagstr = "UTF8String"
	case asn1.TagSequence:
		tagstr = "Sequence"
	case asn1.TagSet:
		tagstr = "Set"
	case asn1.TagPrintableString:
		tagstr = "PrintableString"
	case asn1.TagT61String:
		tagstr = "T61String"
	case asn1.TagIA5String:
		tagstr = "IA5String"
	case asn1.TagUTCTime:
		tagstr = "UTCTime"
	case asn1.TagGeneralizedTime:
		tagstr = "GeneralizedTime"
	case asn1.TagGeneralString:
		tagstr = "GeneralString"
	}

	s += fmt.Sprintf("Class:%s;Tag:%s", clsstr, tagstr)
	return s
}

func (self *Asn1Seq) Format(tabs int) string {
	var s string
	var cur *Asn1Seq
	s = ""
	s += formatTabs(tabs, "{%s;IsCompound:%v;Length(%d:0x%x)}", self.formatClassType(), self.Value.IsCompound, len(self.Value.Bytes), len(self.Value.Bytes))
	s += formatTabs(tabs, "{Bytes:%s}", FormatBytes(self.Value.Bytes, self.Verbose))
	s += formatTabs(tabs, "{FullBytes:%s}", FormatBytes(self.Value.FullBytes, self.Verbose))
	s += formatTabs(tabs, "{Value:%s}", self.formatValue())
	if self.Value.IsCompound || len(self.Child) > 0 {
		for _, cur = range self.Child {
			s += cur.Format(tabs + 1)
		}
	}
	return s
}

func DecodeAsn(data []byte, verbose int) (seq []*Asn1Seq, err error) {
	var rdata []byte
	var rv *asn1.RawValue
	var i int
	var cv *Asn1Seq

	rdata = data
	i = 0
	seq = make([]*Asn1Seq, 0)
	for i = 0; len(rdata) > 0; i++ {
		rv = &asn1.RawValue{}
		rdata, err = asn1.Unmarshal(rdata, rv)
		if err != nil {
			return
		}
		cv = NewAsn1Seq(rv)
		cv.Verbose = verbose
		switch cv.Value.Tag {
		case asn1.TagSequence:
			cv.Child, err = DecodeAsn(cv.Value.Bytes, verbose)
			if err != nil {
				return
			}
		case 0:
			cv.Child, err = DecodeAsn(cv.Value.Bytes, verbose)
			if err != nil {
				return
			}
		case asn1.TagSet:
			cv.Child, err = DecodeAsn(cv.Value.Bytes, verbose)
			if err != nil {
				return
			}
		}
		seq = append(seq, cv)
	}
	err = nil
	return
}

func Pem(infile string, verbose int) error {
	var ps []*pem.Block
	var err error
	var i, j int
	var ap []*Asn1Seq
	var p *Asn1Seq
	ps, err = DecodePem(infile)
	if err != nil {
		return err
	}
	for i = 0; i < len(ps); i++ {
		fmt.Fprintf(os.Stdout, "[%s] decode [%d]\n", infile, i)
		ap, err = DecodeAsn(ps[i].Bytes, verbose)
		if err != nil {
			fmt.Fprintf(os.Stderr, "err[%s]\n", err.Error())
			return err
		}
		for j, p = range ap {
			fmt.Fprintf(os.Stdout, "[%d]\n", j)
			fmt.Fprintf(os.Stdout, "%s", p.Format(1))
		}
	}
	return nil
}

func Der(infile string, verbose int) error {
	var data []byte
	var err error
	var j int
	var ap []*Asn1Seq
	var p *Asn1Seq

	data, err = DecodeDer(infile)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "[%s] decode\n", infile)
	ap, err = DecodeAsn(data, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err[%s]\n", err.Error())
		return err
	}
	for j, p = range ap {
		fmt.Fprintf(os.Stdout, "[%d]\n", j)
		fmt.Fprintf(os.Stdout, "%s", p.Format(1))
	}
	return nil
}

type PemArgs struct {
	Verbose  int
	Password string
	Rsapriv  struct {
		Subnargs []string
	}
	Pem struct {
		Subnargs []string
	}
	Der struct {
		Subnargs []string
	}
	Pkcs12der struct {
		Subnargs []string
	}
	Args []string
}

func decode_rsa_priv(f string, password string, verbose int) error {
	var p *pem.Block
	var data []byte
	var err error
	var rdata []byte
	var ap []*Asn1Seq
	var cur *Asn1Seq
	var i int
	data, err = ioutil.ReadFile(f)
	if err != nil {
		Error("can not read [%s] [%s]", f, err.Error())
		return err
	}
	for len(data) > 0 {
		p, data = pem.Decode(data)
		if err != nil {
			Error("can not decode [%s] [%s]", f, err.Error())
			return err
		}
		rdata, err = x509.DecryptPEMBlock(p, []byte(password))
		if err != nil {
			Error("[%s]decrypt error [%s]\n", f, err.Error())
			return err
		}
		ap, err = DecodeAsn(rdata, verbose)
		if err != nil {
			Error("can not decode rawdata %v err[%s]", rdata, err.Error())
			return err
		}
		for i, cur = range ap {
			fmt.Fprintf(os.Stdout, "[%d]\n", i)
			fmt.Fprintf(os.Stdout, "%s", cur.Format(1))
		}
	}
	return nil
}

func Rsa_priv_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) error {
	var args *PemArgs
	var err error
	var f string
	if ns == nil || ostruct == nil {
		return nil
	}
	args = ostruct.(*PemArgs)
	err = InitLog(ns)
	if err != nil {
		return err
	}

	for _, f = range args.Rsapriv.Subnargs {
		err = decode_rsa_priv(f, args.Password, args.Verbose)
		if err != nil {
			return err
		}
	}
	os.Exit(0)
	return nil
}

func Pem_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) error {
	var args *PemArgs
	var f string
	var err error
	if ns == nil || ostruct == nil {
		return nil
	}
	args = ostruct.(*PemArgs)

	err = InitLog(ns)
	if err != nil {
		return err
	}

	for _, f = range args.Pem.Subnargs {
		err = Pem(f, args.Verbose)
		if err != nil {
			return err
		}
	}
	os.Exit(0)
	return nil
}

func Der_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) error {
	var args *PemArgs
	var f string
	var err error
	if ns == nil || ostruct == nil {
		return nil
	}
	args = ostruct.(*PemArgs)

	err = InitLog(ns)
	if err != nil {
		return err
	}

	for _, f = range args.Der.Subnargs {
		err = Der(f, args.Verbose)
		if err != nil {
			return err
		}
	}
	os.Exit(0)
	return nil
}

func decode_pkcs12_der(fname string, password string, verbose int) error {
	var data []byte
	var err error
	data, err = ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	data = data

	return nil
}

func Pkcs12_der_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) error {
	var args *PemArgs
	var f string
	var err error
	if ns == nil {
		return nil
	}
	args = ostruct.(*PemArgs)
	err = InitLog(ns)
	if err != nil {
		return err
	}

	for _, f = range args.Pkcs12der.Subnargs {
		err = decode_pkcs12_der(f, args.Password, args.Verbose)
		if err != nil {
			return err
		}
	}
	os.Exit(0)
	return nil
}

func init() {
	Rsa_priv_handler(nil, nil, nil)
	Pem_handler(nil, nil, nil)
	Der_handler(nil, nil, nil)
	Pkcs12_der_handler(nil, nil, nil)
}

func main() {
	var commandline = `{
			"password|p" : null,
			"rsapriv<Rsa_priv_handler>" : {
				"$" : "+"
			},
			"pem<Pem_handler>" : {
				"$" : "+"
			},
			"der<Der_handler>" : {
				"$" : "+"
			},
			"pkcs12der<Pkcs12_der_handler>" : {
				"$" : "+"
			}
		}`
	var parser *extargsparse.ExtArgsParse
	var p *PemArgs = &PemArgs{}
	var err error

	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "new args err [%s]\n", err.Error())
		os.Exit(5)
		return
	}
	err = parser.LoadCommandLineString(commandline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load command[%s] err [%s]\n", commandline, err.Error())
		os.Exit(5)
		return
	}

	err = PrepareLog(parser)
	if err != nil {
		fmt.Fprintf(os.Stderr, "prepare log err[%s]\n", err.Error())
		os.Exit(5)
		return
	}

	_, err = parser.ParseCommandLineEx(nil, parser, p, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse command err[%s]\n", err.Error())
		os.Exit(5)
		return
	}
	return

}
