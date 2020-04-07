package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
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

func init() {
	Chan_handler(nil, nil, nil)
}

func main() {
	var commandline string
	var err error
	var parser *extargsparse.ExtArgsParse

	commandline = `{
		"timeout|t" : 500,
		"chan<Chan_handler>##outstr ... to set out string##" : {
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
