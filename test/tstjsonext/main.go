package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"jsonext"
	"os"
)

func init() {
	Getarrayidx_handler(nil, nil, nil)
	Getjson_handler(nil, nil, nil)
}

func Getjson_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	err = nil
	if ns == nil {
		return
	}
	return
}

func Getarrayidx_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	err = nil
	if ns == nil {
		return
	}
	return
}

func main() {
	var parser *extargsparse.ExtArgsParse
	var err error
	var commandline string = `
	{
		"input|i" : null,
		"getjson<Getjson_handler>##path type : type can be int float array map string##" : {
			"$" : "+"
		},
		"getarrayidx<Getarrayidx_handler>##path type idx to get array index type can be int float array map string##" : {
			"$" : 2
		}
	}
	`
	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not new parser [%s]", err.Error())
		os.Exit(5)
	}
	err = parser.LoadCommandLineString(commandline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load [%s] error [%s]", commandline, err.Error())
		os.Exit(5)
	}

	_, err = parser.ParseCommandLine(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error [%s]", err.Error())
		os.Exit(4)
	}
	os.Exit(0)
}
