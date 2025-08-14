package main

import (
	"encoding/pem"
	"fileop"
	"logutil"
)

func read_pem_or_der(infile string) (outb []byte, err error) {
	var inb []byte
	var block *pem.Block
	var leftbytes []byte
	inb, err = fileop.ReadFileBytes(infile)
	if err != nil {
		return
	}
	block, leftbytes = pem.Decode(inb)
	if len(leftbytes) > 0 {
		logutil.DebugBuffer(leftbytes, "leftbytes in [%s]", infile)
		outb = inb
		err = nil
		return
	}
	outb = block.Bytes
	return
}
