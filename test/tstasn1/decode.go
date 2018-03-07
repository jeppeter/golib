package main

import (
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
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
			s += fmt.Sprintf("Integer[%d]", i)
		} else {
			fmt.Fprintf(os.Stderr, "err[%s]\n", err.Error())
		}
	} else if len(self.Value.Bytes) <= 8 {
		var i int64
		_, err = asn1.Unmarshal(self.Value.FullBytes, &i)
		if err == nil {
			s += fmt.Sprintf("Integer[%d]", i)
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
			s += fmt.Sprintf("Integer[%s]", bint.B.String())
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

func (self *Asn1Seq) formatValue() string {
	var s string = ""
	switch self.Value.Tag {
	case asn1.TagInteger:
		s += self.formatIntValue()
	case asn1.TagOID:
		s += self.formatOIDValue()
	}
	return s
}

func (self *Asn1Seq) formatClassType() string {
	var s string = ""
	s += fmt.Sprintf("Class:0x%x;Tag:0x%x", self.Value.Class, self.Value.Tag)
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
