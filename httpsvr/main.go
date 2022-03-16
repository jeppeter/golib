package main

import (
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
)

func main() {
	var commandline = `
	{
		"server<Server_handler>## port listen on port##"  : {
			"$" : 1
		}
	}
	`
	var parser *extargsparse.ExtArgsParse
	var err error

	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		Error("can not make parser err[%s]", err.Error())
		atexit.Exit(5)
	}
	err = parser.LoadCommandLineString(commandline)
	if err != nil {
		Error("can not parse %s", commandline)
		atexit.Exit(5)
	}

	err = PrepareLog(parser)
	if err != nil {
		Error("can not prepare log")
		atexit.Exit(5)
	}
	_, err = parser.ParseCommandLineEx(nil, nil, nil, nil)
	if err != nil {
		Error("can not use parse command line [%s]", err.Error())
		atexit.Exit(4)
	}
	atexit.Exit(0)

}
