package main

import (
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"logutil"
	"os"
	"path/filepath"
)

func init() {
	Tempfile_handler(nil, nil, nil)
	Tempdir_handler(nil, nil, nil)
	Walkdir_handler(nil, nil, nil)
}

func Tempfile_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var idx int
	var pattern string
	var realfile string
	var dirn string
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	dirn = ns.GetString("directory")

	sarr = ns.GetArray("subnargs")

	for idx = 0; idx < len(sarr); idx += 1 {
		pattern = sarr[idx]
		realfile, err = fileop.Mktempfile(dirn, pattern)
		if err != nil {
			return
		}
		fmt.Printf("[%s].[%s] => [%s]\n", dirn, pattern, realfile)
	}
	err = nil
	return

}

func Tempdir_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var idx int
	var pattern string
	var realfile string
	var dirn string
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	dirn = ns.GetString("directory")

	sarr = ns.GetArray("subnargs")

	for idx = 0; idx < len(sarr); idx += 1 {
		pattern = sarr[idx]
		realfile, err = fileop.Mktempdir(dirn, pattern)
		if err != nil {
			return
		}
		fmt.Printf("[%s].[%s] => [%s]\n", dirn, pattern, realfile)
	}
	err = nil
	return
}

func Walkdir_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var idx int
	var curdir string
	var totaln int = 0
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")

	for idx = 0; idx < len(sarr); idx += 1 {
		curdir = sarr[idx]
		err = filepath.Walk(curdir, func(curn string, info os.FileInfo, err2 error) error {
			if !info.IsDir() {
				totaln += 1
			}
			fmt.Printf("curdir [%s] curn %s\n", curdir, curn)
			return nil
		})
	}
	fmt.Printf("totaln %d\n", totaln)
	err = nil
	return
}

func main() {
	var commandline string
	var err error
	var parser *extargsparse.ExtArgsParse
	var ns *extargsparse.NameSpaceEx

	commandline = `{
		"input|i" : null,
		"output|o" : null,
		"directory|D" : null,
		"tempfile<Tempfile_handler>##pattern to create file##" : {
			"$" : "+"
		},
		"tempdir<Tempdir_handler>##pattern to create file##" : {
			"$" : "+"
		},
		"walkdir<Walkdir_handler>##dir... to scan dir##" : {
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
