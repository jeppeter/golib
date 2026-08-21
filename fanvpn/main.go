package main

import (
	"dbgutil"
	"fmt"
	"fileop"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"logutil"
	"os"
)



func Decodejson_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var pemstr string
	var datastr string
	var parser *VPNParse = nil
	var cfg *VPNConfig = nil
	var idx int
	err = nil
	if ns == nil {
		return nil
	}
	err = logutil.InitLog(ns)
	if err != nil {
		logutil.Error("can not Initlog err[%s]", err.Error())
		return err
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = dbgutil.FormatError("need pemfile jsonfile ...")
		return
	}


	pemstr, err = fileop.ReadFile(sarr[0])
	if err != nil {
		return
	}

	parser, err = NewVPNParse(pemstr)
	if err != nil {
		return
	}

	for idx = 1; idx < len(sarr); idx += 1 {
		datastr, err = fileop.ReadFile(sarr[idx])
		if err != nil {
			return
		}
		cfg, err = parser.GetConfig(datastr)
		if err != nil {
			return
		}
		cfg = cfg
	}


	err = nil
	return
}


func init() {
	Decodejson_handler(nil, nil, nil)
}
func main() {
	var commandline string
	var err error
	var parser *extargsparse.ExtArgsParse
	var ns *extargsparse.NameSpaceEx

	commandline = `{
		"input|i" : null,
		"output|o" : null,
		"decjson<Decodejson_handler>##pemfile jsonfile  with decode json for config##" : {
			"$" : "+"
		}

	}`

	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not make parser err[%s]\n", err.Error())
		atexit.Exit(5)
	}

	err = logutil.PrepareLog(parser)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not set [%s]\n", err.Error())
		atexit.Exit(5)
	}
	err = parser.LoadCommandLineString(commandline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not parse %s\n", commandline)
		atexit.Exit(5)
	}


	ns, err = parser.ParseCommandLineEx(nil, nil, nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not use parse command line [%s]\n", err.Error())
		atexit.Exit(4)
	}
	if len(ns.GetString("subcommand")) == 0 {
		fmt.Fprintf(os.Stderr, "can not get subcommand\n")
		atexit.Exit(5)
	}
	//fmt.Fprintf(os.Stdout, "subcommand [%s] succ\n", ns.GetString("subcommand"))
	atexit.Exit(0)
	return
}
