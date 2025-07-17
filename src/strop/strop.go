package strop

import (
	"dbgutil"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func SplitLines(ins string) (retsarr []string) {
	var sarr []string
	var s string
	sarr = strings.Split(ins, "\n")
	retsarr = []string{}
	for _, s = range sarr {
		retsarr = append(retsarr, strings.TrimRight(s, "\r"))
	}
	return
}

func QuoteString(s string) (rets string) {
	var b []byte
	var retb []byte
	var i int
	b = []byte(s)
	retb = append(retb, byte('"'))
	for i = 0; i < len(b); i++ {
		if b[i] == byte('"') || b[i] == byte('\\') {
			retb = append(retb, byte('\\'))
		}
		retb = append(retb, b[i])
	}

	retb = append(retb, byte('"'))

	rets = string(retb)
	return
}

func EncodeBase64(inb []byte) (rets string) {
	rets = base64.StdEncoding.EncodeToString(inb)
	return
}

func DecodeBase64(ins string) (retb []byte, err error) {
	retb, err = base64.StdEncoding.DecodeString(ins)
	return
}

func Base64SplitLines(ins string, linelen int) (outs string) {
	/*all are ascii so make bytes*/
	var inb []byte
	var ridx int
	var cursize int
	var curs string
	inb = []byte(ins)
	outs = ""
	ridx = 0
	for ridx < len(inb) {
		if ridx > 0 {
			outs += "\n"
		}
		cursize = linelen
		if (cursize + ridx) > len(inb) {
			cursize = len(inb) - ridx
		}
		curs = string(inb[ridx:(ridx + cursize)])
		outs += curs
		ridx += cursize
	}
	return
}

func Base64CompactLine(ins string) (outs string) {
	var sarr []string
	var idx int
	var curs string
	outs = ""
	sarr = strings.Split(ins, "\n")
	for idx = 0; idx < len(sarr); idx += 1 {
		curs = strings.TrimRight(sarr[idx], "\r")
		outs += curs
	}
	return
}

func StringToUnicode(ins string) (outb []byte) {
	var inb []byte
	var idx int
	var r rune
	var retn int
	var i int
	var curval int
	inb = []byte(ins)
	outb = []byte{}

	idx = 0
	for idx < len(inb) {
		r, retn = utf8.DecodeRune(inb[idx:])
		for i = 0; i < 2; i++ {
			curval = int(r)
			curval = (curval >> (i * 8)) & 0xff
			outb = append(outb, byte(curval))
		}
		idx += retn
	}
	return
}

func UnicodeToString(inbytes []byte) (outs string, err error) {
	var s string
	var idx, j, retn int
	var r rune
	var rs []rune
	var ps string
	var buf []byte
	var curval int
	var outbytes []byte
	outbytes = []byte{}
	err = nil

	ps = "\""
	for idx = 0; idx < (len(inbytes) - 1); idx += 2 {
		curval = 0
		curval += int(inbytes[idx])
		curval += (int(inbytes[idx+1]) << 8)
		ps += fmt.Sprintf("\\u%04x", curval)
	}
	ps += "\""
	s, err = strconv.Unquote(ps)
	if err != nil {
		err = dbgutil.FormatError("[%s] error [%s]", ps, err.Error())
		return
	}

	buf = make([]byte, 10)

	idx = 0
	rs = []rune(s)
	for idx = 0; idx < len(rs); idx++ {
		r = rs[idx]
		retn = utf8.EncodeRune(buf, r)
		for j = 0; j < retn; j++ {
			outbytes = append(outbytes, buf[j])
		}
	}
	outs = string(outbytes)
	err = nil
	return
}

func Parseu64(val string) (vali uint64, err error) {
	var ss string = val
	var base int = 10
	if strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
		ss = val[2:]
		base = 16
	} else if strings.HasPrefix(val, "x") || strings.HasPrefix(val, "X") {
		ss = val[1:]
		base = 16
	}
	vali, err = strconv.ParseUint(ss, base, 64)
	if err != nil {
		err = dbgutil.FormatError("parse [%s] error [%s]", val, err.Error())
		return
	}
	return
}
