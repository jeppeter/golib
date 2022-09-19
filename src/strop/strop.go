package strop

import (
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
