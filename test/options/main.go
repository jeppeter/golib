package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	var c string
	var data []byte
	var err error
	var p *ExtArgsOptions
	for _, c = range os.Args[1:] {
		data, err = ioutil.ReadFile(c)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can not read [%s] err[%s]\n", c, err.Error())
			return
		}
		p, err = NewExtArgsOptions(string(data))
		if err != nil {
			fmt.Fprintf(os.Stderr, "can not parse data [%s] err[%s]\n", string(data), err.Error())
			return
		}

		fmt.Fprintf(os.Stdout, "%s\n", p.Format())
	}
}
