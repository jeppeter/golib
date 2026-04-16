package main

import (
	"dbgutil"
	"encoding/json"
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"jsonext"
	"logutil"
)

func init() {
	Testjsonload_handler(nil, nil, nil)
	Repack_handler(nil, nil, nil)
}

type AliSmsConfig struct {
	Cname        string   `json:"cname"`
	Note         string   `json:"note"`
	Hostname     string   `json:"hostname"`
	Count        int      `json:"count"`
	Phonenumber  []string `json:"phonenumber"`
	Templatecode string   `json:"templatecode"`
}

type JsonFrom struct {
	Fromid  string `json:"fromid"`
	Fromval string `json:"fromval"`
}

func LoadParser(parser *extargsparse.ExtArgsParse) (err error) {
	var commandline_fmt string
	var commandline string
	commandline_fmt = `{
		"testjsonload<Testjsonload_handler>##filenames ... to test json load##" : {
			"$" : "+"
		},
		"repack<Repack_handler>##fname fromid fromval to repack##" : {
			"$" : 3
		}
	}`

	commandline = fmt.Sprintf(commandline_fmt)
	err = parser.LoadCommandLineString(commandline)
	return
}

func repack_json(ins string, f *JsonFrom) (outs string, err error) {
	var vmap map[string]interface{}
	outs = ""
	vmap, err = jsonext.GetJsonMap(ins)
	if err != nil {
		return
	}

	if f != nil {
		if len(f.Fromid) > 0 {
			vmap, err = jsonext.SetJsonValue("fromid", "string", f.Fromid, vmap)
			if err != nil {
				return
			}
		}

		if len(f.Fromval) > 0 {
			vmap, err = jsonext.SetJsonValue("fromval", "string", f.Fromval, vmap)
			if err != nil {
				return
			}
		}
	}

	outs, err = jsonext.FormatJsonValue(0, "", vmap)
	return
}

func Testjsonload_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var fname string
	var fdata []byte
	var outdata []byte
	var retp *AliSmsConfig = nil
	var fromptr *JsonFrom = nil
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
		fromptr = &JsonFrom{}
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

		err = json.Unmarshal(fdata, fromptr)
		if err != nil {
			err = dbgutil.FormatError("unmarshal %s JsonFrom error %s", fname, err.Error())
			return
		}

		fmt.Printf("cname [%s] hostname [%s] phonenumber %v\n", retp.Cname, retp.Hostname, retp.Phonenumber)
		fmt.Printf("Fromid [%s] Fromval [%s]\n", fromptr.Fromid, fromptr.Fromval)
		logutil.DebugBuffer(outdata, "output")
		fmt.Printf("%s output\n%s\n", fname, string(outdata))
	}

	err = nil
	return
}

func Repack_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var fname string
	var fdata []byte
	var outdata []byte
	var retp *AliSmsConfig = nil
	var fromptr *JsonFrom = nil
	var outs string
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
	if len(sarr) < 3 {
		err = dbgutil.FormatError("need fname fromid fromval")
		return
	}

	fname = sarr[0]
	fdata, err = fileop.ReadFileBytes(fname)
	if err != nil {
		return
	}

	retp = &AliSmsConfig{}

	err = json.Unmarshal(fdata, retp)
	if err != nil {
		return
	}

	fromptr = &JsonFrom{}

	fromptr.Fromid = sarr[1]
	fromptr.Fromval = sarr[2]
	outdata, err = json.Marshal(retp)
	if err != nil {
		err = dbgutil.FormatError("marshal error %s", err.Error())
		return
	}
	outs, err = repack_json(string(outdata), fromptr)
	if err != nil {
		return
	}

	fmt.Printf("fname %s\n%s\nrepack\n%s\n", fname, string(fdata), outs)
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
