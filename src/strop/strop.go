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
