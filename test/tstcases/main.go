package main

import (
	"bytes"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func go_chan(c chan string, e chan int, timeout int) {
	var s string
	var icnt int = 0
	for s = range c {
		fmt.Fprintf(os.Stdout, "%s\n", s)
		icnt++
		time.Sleep(time.Duration(timeout) * time.Millisecond)
	}
	e <- icnt
	return
}

func Chan_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var s string
	var i int
	var schan chan string
	var ichan chan int
	err = nil
	if ns == nil {
		return
	}

	schan = make(chan string, 100)
	ichan = make(chan int, 1)
	go go_chan(schan, ichan, ns.GetInt("timeout"))

	for _, s = range ns.GetArray("subnargs") {
		schan <- s
	}
	close(schan)
	schan = nil
	i = <-ichan

	fmt.Fprintf(os.Stdout, "end value [%d]\n", i)

	return
}

func get_code(sarr []string) (retbyte []byte, err error) {
	var s string
	var v int
	retbyte = []byte{}
	err = nil

	for _, s = range sarr {
		base = 10
		if strings.HasPrefix("0x") ||
			strings.HasPrefix("0X") {
			base = 16
			s = s[2:]
		} else if strings.HasPrefix("x") ||
			strings.HasPrefix("X") {
			base = 16
			s = s[1:]
		}
		v, err = strconv.ParseInt(s, base)
		if err != nil {
			return
		}
		retbyte = append(retbyte, byte(v))
	}

	err = nil
	return
}

func gbk_to_utf8(inbytes []byte) (outbytes []byte, err error) {
	var rd *transform.Reader
	rd = transform.NewReader(bytes.NewReader(inbytes), simplifiedchinese.GBK.NewDecoder())
	outbytes, err = ioutil.ReadAll(rd)
	if err != nil {
		return
	}
	return
}

func utf8_to_gbk(inbytes []byte) (outbytes []byte, err error) {
	var rd *transform.Reader
	rd = transform.NewReader(bytes.NewReader(inbytes), simplifiedchinese.GBK.NewEncoder())
	outbytes, err = ioutil.ReadAll(rd)
	if err != nil {
		return
	}
	return
}

func out_bytes(inbytes []byte, a ...interface{}) (outs string) {
	var lasti, i int
	var b byte
	outs = ""
	lasti = 0
	i = 0
	outs += fmt.Sprintf("bytes [%d:0x%x] ", len(inbytes), len(inbytes))
	outs += fmt.Sprintf(a...)
	for i, b = range inbytes {
		if (i % 16) == 0 {
			if i > 0 {
				outs += "    "
				for lasti != i {
					if inbytes[lasti] >= byte(' ') &&
						inbytes[lasti] <= byte('~') {
						outs += fmt.Sprintf("%c", inbytes[lasti])
					} else {
						outs += "."
					}
					lasti++
				}
			}
			outs += fmt.Sprintf("\n0x%08x:", i)
		}
		outs += fmt.Sprintf(" 0x%02x", b)
	}

	if lasti != i {
		for (i % 16) != 0 {
			outs += fmt.Sprintf("     ")
			i++
		}
		outs += fmt.Sprintf("    ")

		for lasti != len(inbytes) {
			if inbytes[lasti] >= byte(' ') &&
				inbytes[lasti] <= byte('~') {
				outs += fmt.Sprintf("%c", inbytes[lasti])
			} else {
				outs += "."
			}
			lasti++
		}

		outs += fmt.Sprintf("\n")
	}

	return
}

func Gbktoutf8_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var inbytes, outbytes []byte
	err = nil
	if ns == nil {
		return
	}
	inbytes, err = get_code(ns.GetArray("subnargs"))
	if err != nil {
		return
	}

	outbytes, err = gbk_to_utf8(inbytes)
	if err != nil {
		return
	}

	fmt.Printf(os.Stdout, "%s", out_bytes(inbytes, "input bytes"))
	fmt.Fprintf(os.Stdout, "%s", out_bytes(outbytes, "output bytes"))
	err = nil

	return
}

func Utf8togbk_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var inbytes, outbytes []byte
	err = nil
	if ns == nil {
		return
	}
	inbytes, err = get_code(ns.GetArray("subnargs"))
	if err != nil {
		return
	}

	outbytes, err = utf8_to_gbk(inbytes)
	if err != nil {
		return
	}

	fmt.Printf(os.Stdout, "%s", out_bytes(inbytes, "input bytes"))
	fmt.Fprintf(os.Stdout, "%s", out_bytes(outbytes, "output bytes"))
	err = nil

	return
}

func init() {
	Chan_handler(nil, nil, nil)
	Ansitogbk_handler(nil, nil, nil)
	Gbktoansi_handler(nil, nil, nil)
}

func main() {
	var commandline string
	var err error
	var parser *extargsparse.ExtArgsParse

	commandline = `{
		"timeout|t" : 500,
		"chan<Chan_handler>##outstr ... to set out string##" : {
			"$" : "+"
		},
		"utf8togbk<Utf8togbk_handler>## codes ... to get codes from utf-8 to ansi##"  : {
			"$" : "+"
		},
		"Gbktoutf8<Gbktoutf8_handler>## codes ... to get codes from ansi to utf-8##" : {
			"$" : "+"
		}
	}`

	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not make parser err[%s]\n", err.Error())
		os.Exit(5)
	}
	err = parser.LoadCommandLineString(commandline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not parse %s\n", commandline)
		os.Exit(5)
	}

	_, err = parser.ParseCommandLineEx(nil, nil, nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not use parse command line [%s]\n", err.Error())
		os.Exit(4)
	}
	return
}
