package main

import (
	"aesext"
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
	var key []byte
	var iv []byte
	var data []byte
	var encbytes []byte
	var ofile string
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
	if len(sarr) < 3 {
		err = dbgutil.FormatError("need aesfile ivfile datafile")
		return
	}

	ofile = ns.GetString("output")

	key, err = fileop.ReadFileBytes(sarr[0])
	if err != nil {
		return
	}

	iv, err = fileop.ReadFileBytes(sarr[1])
	if err != nil {
		return
	}


	data, err = fileop.ReadFileBytes(sarr[2])
	if err != nil {
		return
	}

	encbytes, err = aesext.AesEncCbc(data,key,iv)
	if err != nil {
		return
	}

	_, err = fileop.WriteFileBytes(ofile,encbytes)
	if err != nil {
		return
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
		"decjson<Decodejson_handler>##jsonfile ... with decode json for config##" : {
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
