package main

import (
	"dbgutil"
	"encoding/json"
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"logutil"
)

func init() {
	Testjsonload_handler(nil, nil, nil)
}

type AliSmsConfig struct {
	Cname        string   `json:"cname"`
	Note         string   `json:"note"`
	Hostname     string   `json:"hostname"`
	Count        int      `json:"count"`
	Phonenumber  []string `json:"phonenumber"`
	Templatecode string   `json:"templatecode"`
}

func LoadParser(parser *extargsparse.ExtArgsParse) (err error) {
	var commandline_fmt string
	var commandline string
	commandline_fmt = `{
		"testjsonload<Testjsonload_handler>##filenames ... to test json load##" : {
			"$" : "+"
		}
	}`

	commandline = fmt.Sprintf(commandline_fmt)
	err = parser.LoadCommandLineString(commandline)
	return
}

func Testjsonload_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var fname string
	var fdata []byte
	var outdata []byte
	var retp *AliSmsConfig = nil
	var sarr []string
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 1 {
		err = dbgutil.FormatError("need fname")
		return
	}

	for _, fname = range sarr {
		fdata, err = fileop.ReadFileBytes(fname)
		if err != nil {
			return
		}
		retp = &AliSmsConfig{}
		logutil.DebugBuffer(fdata, "read %s", fname)
		err = json.Unmarshal(fdata, retp)
		if err != nil {
			err = dbgutil.FormatError("unmarshal %s error %s", fname, err.Error())
			return
		}

		outdata, err = json.Marshal(retp)
		if err != nil {
			err = dbgutil.FormatError("marshal %s error %s", fname, err.Error())
			return
		}

		fmt.Printf("cname [%s] hostname [%s] phonenumber %v\n", retp.Cname, retp.Hostname, retp.Phonenumber)
		logutil.DebugBuffer(outdata, "output")
		fmt.Printf("%s output\n%s\n", fname, string(outdata))
	}

	err = nil
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
