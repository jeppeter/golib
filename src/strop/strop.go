package strop

import (
	"encoding/base64"
	"strings"
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
