package main

import (
	"dbgutil"
	"encoding/json"
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"logutil"
	"slices"
)

func init() {
	Add_handler(nil, nil, nil)
	Del_handler(nil, nil, nil)
}

func LoadParser(parser *extargsparse.ExtArgsParse) (err error) {
	var commandline_fmt string
	var commandline string
	commandline_fmt = `{
		"input|i" : null,
		"output|o" : null,
		"add<Add_handler>##to add string into ##" : {
			"$" : "+"
		},
		"del<Del_handler>##to delete string value##" : {
			"$" : "+"
		}
	}`

	commandline = fmt.Sprintf(commandline_fmt)
	err = parser.LoadCommandLineString(commandline)
	return
}

func Add_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var infile string
	var outfile string
	var strs []string = []string{}
	var inb []byte
	var sarr []string
	var i int
	var outb []byte
	err = nil
	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	infile = ns.GetString("input")
	outfile = ns.GetString("output")
	sarr = ns.GetArray("subnargs")

	inb, err = fileop.ReadFileBytes(infile)
	if err == nil && len(inb) > 0 {
		err = json.Unmarshal(inb, &strs)
		if err != nil {
			err = dbgutil.FormatError("can not parse [%s] error %s", string(inb), err.Error())
			return
		}
	}

	for i = 0; i < len(sarr); i += 1 {
		strs = slices.Insert(strs, len(strs), sarr[i])
	}

	outb, err = json.Marshal(&strs)
	if err != nil {
		err = dbgutil.FormatError("can not %v format json error %s", strs, err.Error())
		return
	}

	_, err = fileop.WriteFileBytes(outfile, outb)
	if err != nil {
		return
	}

	return
}

func Del_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var infile string
	var outfile string
	var strs []string = []string{}
	var inb []byte
	var sarr []string
	var outb []byte
	err = nil
	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	infile = ns.GetString("input")
	outfile = ns.GetString("output")
	sarr = ns.GetArray("subnargs")

	inb, err = fileop.ReadFileBytes(infile)
	if err == nil && len(inb) > 0 {
		err = json.Unmarshal(inb, &strs)
		if err != nil {
			err = dbgutil.FormatError("can not parse [%s] error %s", string(inb), err.Error())
			return
		}
	}

	strs = slices.DeleteFunc(strs, func(vi string) bool {
		var ci int
		for ci = 0; ci < len(sarr); ci += 1 {
			if vi == sarr[ci] {
				return true
			}
		}
		return false
	})

	outb, err = json.Marshal(&strs)
	if err != nil {
		err = dbgutil.FormatError("can not %v format json error %s", strs, err.Error())
		return
	}

	_, err = fileop.WriteFileBytes(outfile, outb)
	if err != nil {
		return
	}

	return
}

func main() {
	var parser *extargsparse.ExtArgsParse
	var err error
	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		logutil.Error("%s", err.Error())
		atexit.Exit(5)
	}

	err = LoadParser(parser)
	if err != nil {
		logutil.Error("%s", err.Error())
		atexit.Exit(5)
	}

	err = logutil.PrepareLog(parser)
	if err != nil {
		logutil.Error("%s", err.Error())
		atexit.Exit(5)
	}

	_, err = parser.ParseCommandLine(nil, nil)
	if err != nil {
		logutil.Error("%s", err.Error())
		atexit.Exit(4)
	}
	atexit.Exit(0)
}
