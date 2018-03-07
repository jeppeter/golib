package main

import (
	"encoding/asn1"
	"fmt"
	"math/big"
)

type ecdsa struct {
	R, S *big.Int
}

func main() {
	r, _ := new(big.Int).SetString("316eb3cad8b66fcf1494a6e6f9542c3555addbf337f04b62bf4758483fdc881d", 16)
	s, _ := new(big.Int).SetString("bf46d26cef45d998a2cb5d2d0b8342d70973fa7c3c37ae72234696524b2bc812", 16)
	sequence := ecdsa{r, s}
	encoding, _ := asn1.Marshal(sequence)
	fmt.Println(encoding)
	dec := new(ecdsa)
	asn1.Unmarshal(encoding, dec)
	fmt.Println(dec)
}
