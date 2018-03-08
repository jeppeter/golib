package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"os"
)

func debug_out_cmd(parser *extargsparse.ExtArgsParse, cmdname string, mode int) error {
	var subcmds []string
	var opts []*extargsparse.ExtKeyParse
	var curopt *extargsparse.ExtKeyParse
	var i int
	var curcmd, c string
	var err error
	opts, err = parser.GetCmdOpts(cmdname)
	if err != nil {
		return err
	}
	for i, curopt = range opts {
		if mode <= 0 {
			Error("%s.[%d]=%s", cmdname, i, curopt.Format())
		} else if mode == 1 {
			Warn("%s.[%d]=%s", cmdname, i, curopt.Format())
		} else if mode == 2 {
			Info("%s.[%d]=%s", cmdname, i, curopt.Format())
		} else if mode == 3 {
			Debug("%s.[%d]=%s", cmdname, i, curopt.Format())
		} else if mode >= 4 {
			Trace("%s.[%d]=%s", cmdname, i, curopt.Format())
		}
	}

	subcmds, err = parser.GetSubCommands(cmdname)
	if err != nil {
		return err
	}

	if len(subcmds) > 0 {
		for _, curcmd = range subcmds {
			c = cmdname
			if len(c) > 0 {
				c += "."
			}
			c += curcmd

			err = debug_out_cmd(parser, c, mode)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Cmd_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) error {
	var parser *extargsparse.ExtArgsParse
	var err error
	if ns == nil || ctx == nil {
		return nil
	}

	parser = ctx.(*extargsparse.ExtArgsParse)
	err = InitLog(ns)
	if err != nil {
		return err
	}

	err = debug_out_cmd(parser, "", ns.GetInt("verbose"))
	if err != nil {
		atexit.Exit(4)
	}
	atexit.Exit(0)
	return nil
}

func init() {
	Cmd_handler(nil, nil, nil)
}

func main() {
	var commandline = `{
		"cmd<cmd_handler>##cmd debug out##" : {
			"$" : "*"
		}
	}`

	var parser *extargsparse.ExtArgsParse
	var err error

	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	err = parser.LoadCommandLineString(commandline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load command err [%s]\n", err.Error())
		atexit.Exit(5)
		return
	}

	err = PrepareLog(parser)
	if err != nil {
		fmt.Fprintf(os.Stderr, "prepare log err[%s]\n", err.Error())
		atexit.Exit(5)
		return
	}

	_, err = parser.ParseCommandLine(nil, parser)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse command err[%s]\n", err.Error())
		atexit.Exit(5)
		return
	}
	atexit.Exit(0)
	return
}
