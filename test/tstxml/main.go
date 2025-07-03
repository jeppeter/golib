package main

import (
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"logutil"
	"os"
	"strings"
	"xmlext"
)

func Xmlparse_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var f string
	var ins string
	var xext *xmlext.XmlExt
	var attrs map[string]string
	err = nil
	if ns == nil {
		return
	}
	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	for _, f = range sarr {
		ins, err = fileop.ReadFile(f)
		if err != nil {
			return
		}
		xext, err = xmlext.ParseXmlExt(ins)
		if err != nil {
			return
		}
		attrs, err = xext.GetAttrs("")
		if err != nil {
			return
		}
		fmt.Printf("format\n%s\nindent\n%s\n", xext.String(), xext.Ident(0))
		for k, v := range attrs {
			fmt.Printf("[%s]=[%s]\n", k, v)
		}
	}

	return
}

func Getset_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var f string
	var ins string
	var xext *xmlext.XmlExt
	var input string
	var output string
	var kvarr []string
	var pattr []string
	var val string
	var outs string
	err = nil
	if ns == nil {
		return
	}
	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	input = ns.GetString("input")
	output = ns.GetString("output")
	ins, err = fileop.ReadFile(input)
	if err != nil {
		return
	}
	xext, err = xmlext.ParseXmlExt(ins)
	if err != nil {
		return
	}

	for _, f = range sarr {
		kvarr = strings.SplitN(f, "=", 2)
		if len(kvarr) == 1 {
			pattr = strings.SplitN(kvarr[0], ".", 2)
			if len(pattr) == 2 {
				val, err = xext.GetAttrValue(pattr[0], pattr[1])
				if err != nil {
					logutil.Error("can not get [%s].[%s] error %s", pattr[0], pattr[1], err.Error())
				} else {
					fmt.Printf("[%s].[%s] value [%s]\n", pattr[0], pattr[1], val)
				}
			} else {
				val, err = xext.GetValue(pattr[0])
				if err != nil {
					logutil.Error("can not get [%s] error %s", pattr[0], err.Error())
				} else {
					fmt.Printf("[%s] value [%s]\n", pattr[0], val)
				}
			}
		} else {
			pattr = strings.SplitN(kvarr[0], ".", 2)
			if len(pattr) == 2 {
				val, err = xext.SetAttr(pattr[0], pattr[1], kvarr[1])
				if err != nil {
					logutil.Error("can not set [%s].[%s] = [%s] error %s", pattr[0], pattr[1], kvarr[1], err.Error())
				} else {
					fmt.Printf("[%s].[%s] = [%s] retval [%s]\n", pattr[0], pattr[1], kvarr[1], val)
				}
			} else {
				val, err = xext.SetValue(pattr[0], kvarr[1])
				if err != nil {
					logutil.Error("can not set [%s] = [%s] error %s", pattr[0], kvarr[1], err.Error())
				} else {
					fmt.Printf("[%s] = [%s] retval [%s]\n", pattr[0], kvarr[1], val)
				}
			}
		}
	}
	outs = xext.String()

	if len(output) != 0 {
		_, err = fileop.WriteFile(output, outs)
		if err != nil {
			return
		}
	} else {
		fmt.Printf("%s", outs)
	}

	err = nil
	return
}

func init() {
	Xmlparse_handler(nil, nil, nil)
	Getset_handler(nil, nil, nil)
}
func main() {
	var commandline string
	var err error
	var parser *extargsparse.ExtArgsParse
	var ns *extargsparse.NameSpaceEx

	commandline = `{
		"input|i" : null,
		"output|o" : null,
		"xmlparse<Xmlparse_handler>##files ... to parse xml##" : {
			"$" : "+"
		},
		"getset<Getset_handler>##key[=value] to set attr like root.status_code=200 to read attr root.status_code to set value root=cc get valueroot##" : {
			"$" : "*"
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
