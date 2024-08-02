package main

import (
	"executil"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"logutil"
	"os"
	"strconv"
)

func init() {
	Runtimeout_handler(nil, nil, nil)
	Getexec_handler(nil, nil, nil)
	Getbootticks_handler(nil, nil, nil)
	Childpid_handler(nil, nil, nil)
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

	logutil.Debug("%v exitcode %d", sarr, exitcode)
	logutil.Debug("%v outs\n%s", sarr, outs)
	logutil.Debug("%v errs\n%s", sarr, errs)
	return
}

func Getexec_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var exename string
	var dirname string
	err = nil
	if ns == nil {
		return
	}
	err = logutil.InitLog(ns)
	if err != nil {
		return
	}
	exename, err = executil.GetExeFile()
	if err != nil {
		return
	}
	dirname, err = executil.GetExeDir()
	if err != nil {
		return
	}
	fmt.Printf("exec %s dir %s\n", exename, dirname)
	return
}

func Getbootticks_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var ticks uint64
	err = nil
	if ns == nil {
		return
	}
	err = logutil.InitLog(ns)
	if err != nil {
		return
	}
	ticks, err = executil.GetTicksFromBoot()
	if err != nil {
		return
	}
	fmt.Printf("ticks %d\n", ticks)
	return
}

func Childpid_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var rp *executil.ChildProcs = nil
	var i int
	var pid int
	err = nil
	if ns == nil {
		return
	}
	err = logutil.InitLog(ns)
	if err != nil {
		return
	}
	sarr = ns.GetArray("subnargs")
	for i = 0; i < len(sarr); i++ {
		pid, err = strconv.Atoi(sarr[i])
		if err != nil {
			return
		}
		rp, err = executil.GetChildProcs(pid)
		if err != nil {
			return
		}
		fmt.Printf("ChildPids\n%s", rp.String())
	}
	err = nil
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
		},
		"getexec<Getexec_handler>##to get execname and file##" : {
			"$" : 0
		},
		"getbootticks<Getbootticks_handler>##to get ticks##" : {
			"$" :0
		},
		"childpid<Childpid_handler>##pid to list child##" : {
			"$": "+"
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
