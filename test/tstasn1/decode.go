package main

import (
	"encoding/asn1"
	"encoding/pem"
	"fmt"
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

type Asn1Seq struct {
	Value *asn1.RawValue
	Child []*Asn1Seq
}

func NewAsn1Seq(v *asn1.RawValue) *Asn1Seq {
	var self *Asn1Seq
	self = &Asn1Seq{}
	self.Value = v
	self.Child = make([]*Asn1Seq, 0)
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

func FormatBytes(data []byte) string {
	var s string = ""
	var b byte
	var i int
	s += "["
	for i, b = range data {
		if i > 0 {
			s += ","
		}
		s += fmt.Sprintf("0x%02x", b)
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
			s += fmt.Sprintf("Integer[%d:0x%x]", i, i)
		} else {
			fmt.Fprintf(os.Stderr, "err[%s]\n", err.Error())
		}
	} else if len(self.Value.Bytes) <= 8 {
		var i int64
		_, err = asn1.Unmarshal(self.Value.FullBytes, &i)
		if err == nil {
			s += fmt.Sprintf("Integer[%d:0x%x]", i, i)
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
			s += fmt.Sprintf("Integer[%s:0x%s]", bint.B.Text(10), bint.B.Text(16))
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
	s += formatTabs(tabs, "{Bytes:%s}", FormatBytes(self.Value.Bytes))
	s += formatTabs(tabs, "{FullBytes:%s}", FormatBytes(self.Value.FullBytes))
	s += formatTabs(tabs, "{Value:%s}", self.formatValue())
	if self.Value.IsCompound || len(self.Child) > 0 {
		for _, cur = range self.Child {
			s += cur.Format(tabs + 1)
		}
	}
	return s
}

func DecodeAsn(data []byte) (seq []*Asn1Seq, err error) {
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
		switch cv.Value.Tag {
		case asn1.TagSequence:
			cv.Child, err = DecodeAsn(cv.Value.Bytes)
			if err != nil {
				return
			}
		case 0:
			cv.Child, err = DecodeAsn(cv.Value.Bytes)
			if err != nil {
				return
			}
		case asn1.TagSet:
			cv.Child, err = DecodeAsn(cv.Value.Bytes)
			if err != nil {
				return
			}
		}
		seq = append(seq, cv)
	}
	err = nil
	return
}

func Pem(infile string) error {
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
		ap, err = DecodeAsn(ps[i].Bytes)
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

func main() {
	for _, c := range os.Args[1:] {
		Pem(c)
	}
}
