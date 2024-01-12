package main

import (
	"executil"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"logutil"
	"os"
)

func init() {
	Runtimeout_handler(nil, nil, nil)
}

func Runtimeout_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var timeout int
	var exitcode int
	var outs string
	var errs string
	err = nil
	if ns == nil {
		return
	}
	err = logutil.InitLog(ns)
	if err != nil {
		return
	}
	sarr = ns.GetArray("subnargs")
	timeout = ns.GetInt("timeout")
	outs, errs, exitcode, err = executil.RunCmdTimeout(sarr, timeout)
	if err != nil {
		return
	}

	logutil.Debug("%v exitcode %d", exitcode)
	logutil.Debug("%v outs\n%s", sarr, outs)
	logutil.Debug("%v errs\n%s", sarr, errs)
	return
}

func main() {
	var parser *extargsparse.ExtArgsParse
	var err error
	var commandline string = `
	{
		"input|i" : null,
		"output|o" : null,
		"timeout|t" : 0,
		"runtimeout<Runtimeout_handler>##args ...##" : {
			"$" : "+"
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

	err = logutil.PrepareLog(parser)
	if err != nil {
		fmt.Fprintf(os.Stderr, "PrepareLog error [%s]", err.Error())
		os.Exit(5)
	}

	_, err = parser.ParseCommandLine(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error [%s]", err.Error())
		os.Exit(4)
	}
	os.Exit(0)
}
