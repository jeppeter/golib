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
	"reflect"
	"regexp"
	"strings"
	"time"
	"unsafe"
)

type HuigouInfo struct {
	SecurityCode       string  `json:"SECURITY_CODE"`
	DeriveSecurity     string  `json:"DERIVE_SECURITY_CODE"`
	SecurityName       string  `json:"SECURITY_NAME"`
	ChangeDate         string  `json:"CHANGE_DATE"`
	PersonName         string  `json:"PERSON_NAME"`
	ChangeShares       int64   `json:"CHANGE_SHARES"`
	AveragePrice       float64 `json:"AVERAGE_PRICE"`
	ChangeAmount       float64 `json:"CHANGE_AMOUNT"`
	ChangeReason       string  `json:"CHANGE_REASON"`
	ChangeRatio        float64 `json:"CHANGE_RATIO"`
	ChangeAfterHoldNum int64   `json:"CHANGE_AFTER_HOLDNUM"`
	HoldType           string  `json:"HOLD_TYPE"`
	DsePersonName      string  `json:"DSE_PERSON_NAME"`
	PositionName       string  `json:"POSITION_NAME"`
	PersonDseRelation  string  `json:"PERSON_DSE_RELATION"`
	OrgCode            string  `json:"ORG_CODE"`
	GGEid              int64   `json:"GGEID"`
	BeginHoldNum       int64   `json:"BEGIN_HOLD_NUM"`
	EndHoldNum         int64   `json:"END_HOLD_NUM"`
}

func (hg *HuigouInfo) get_epoch(s string) (rval int64, err error) {
	var ntime time.Time
	rval = 0

	ntime, err = time.Parse(time.DateTime, s)
	logutil.Debug("ntime %v s %s", ntime, s)
	rval = ntime.Unix()

	err = nil
	return
}

func (hg *HuigouInfo) FormatCreateTable(table string) (outs string, err error) {
	var kname string
	var pt reflect.Type
	var rf reflect.Value
	var i int
	var types string
	var reg *regexp.Regexp
	var bmatch bool

	reg, err = regexp.Compile(".*date.*")

	err = nil
	outs = ""
	outs += fmt.Sprintf("create table %s (", table)

	rf = reflect.ValueOf(hg).Elem()
	pt = rf.Type()
	for i = 0; i < rf.NumField(); i += 1 {
		kname = strings.ToLower(pt.Field(i).Name)
		types = pt.Field(i).Type.String()

		if i > 0 {
			outs += ","
		}

		if types == "int64" {
			outs += fmt.Sprintf("%s int64", kname)
		} else if types == "string" {
			/**/
			bmatch = reg.MatchString(kname)
			if bmatch {
				outs += fmt.Sprintf("%s int64", kname)
			} else {
				outs += fmt.Sprintf("%s text", kname)
			}

		} else if types == "float64" {
			outs += fmt.Sprintf("%s double", kname)
		}

	}
	outs += ")"

	err = nil
	return
}

func (hg *HuigouInfo) FormatInsert() (keys string, vals string, err error) {
	var kname string
	var pt reflect.Type
	var rf, vrf reflect.Value
	var i int
	var types string
	var reg *regexp.Regexp
	var bmatch bool
	var iv64 int64
	var fv64 float64
	var s string
	var ok bool

	reg, err = regexp.Compile(".*date.*")

	err = nil
	keys = "("
	vals = "("

	rf = reflect.ValueOf(hg).Elem()
	pt = rf.Type()
	for i = 0; i < rf.NumField(); i += 1 {
		kname = strings.ToLower(pt.Field(i).Name)
		types = pt.Field(i).Type.String()

		if i > 0 {
			keys += ","
			vals += ","
		}

		vrf = rf.Field(i)
		vrf = reflect.NewAt(vrf.Type(), unsafe.Pointer(vrf.UnsafeAddr())).Elem()

		if types == "int64" {
			iv64, ok = vrf.Interface().(int64)
			if !ok {
				err = dbgutil.FormatError("%s not int64", kname)
				return
			}
			keys += kname
			vals += fmt.Sprintf("%d", iv64)
		} else if types == "string" {
			/**/
			s, ok = vrf.Interface().(string)
			if !ok {
				err = dbgutil.FormatError("%s not string", kname)
				return
			}
			bmatch = reg.MatchString(kname)
			if bmatch {
				iv64, err = hg.get_epoch(s)
				if err != nil {
					return
				}
				keys += kname
				vals += fmt.Sprintf("%d", iv64)
			} else {
				keys += kname
				vals += fmt.Sprintf(`"%s"`, s)
			}

		} else if types == "float64" {
			fv64, ok = vrf.Interface().(float64)
			if !ok {
				err = dbgutil.FormatError("%s not float64", kname)
				return
			}
			keys += kname
			vals += fmt.Sprintf("%f", fv64)
		}

	}

	keys += ")"
	vals += ")"
	return
}

func init() {
	Testjsonload_handler(nil, nil, nil)
	Repack_handler(nil, nil, nil)
	Refvals_handler(nil, nil, nil)
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
		},
		"refvals<Refvals_handler>##infile to reflect values##" : {
			"$" : 1
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

func Refvals_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var fname string
	var fdata []byte
	var sarr []string
	var keys, vals string
	var creates string
	var retp *HuigouInfo
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

	fname = sarr[0]
	fdata, err = fileop.ReadFileBytes(fname)
	if err != nil {
		return
	}

	retp = &HuigouInfo{}

	err = json.Unmarshal(fdata, retp)
	if err != nil {
		return
	}

	keys, vals, err = retp.FormatInsert()
	if err != nil {
		return
	}

	creates, err = retp.FormatCreateTable("huigou")
	if err != nil {
		return
	}

	fmt.Printf("keys \n%s\nvals\n%s\n", keys, vals)
	fmt.Printf("create sql\n%s\n", creates)

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
