package main

import (
	"encoding/asn1"
	"fmt"
	"math/big"
	"os"
)

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

func makeBigInt(s string) error {
	var err error
	var i *bigint
	var j *bigint
	var data []byte
	i = new(bigint)
	i.B = big.NewInt(0)
	i.B.SetString(s, 16)
	data, err = asn1.Marshal(*i)
	if err == nil {
		fmt.Fprintf(os.Stdout, "format [%s] = %s\n", s, FormatBytes(data))
	} else {
		fmt.Fprintf(os.Stderr, "can not marshal [%s] err[%s]\n", s, err.Error())
		return err
	}

	j = new(bigint)
	j.B = big.NewInt(0)
	_, err = asn1.Unmarshal(data, j)
	if err == nil {
		fmt.Fprintf(os.Stdout, "decode %s = [%+v]\n", FormatBytes(data), j)
	} else {
		fmt.Fprintf(os.Stderr, "error [%s]\n", err.Error())
	}

	return nil
}

func main() {
	for _, c := range os.Args[1:] {
		makeBigInt(c)
	}
}
