package main

import (
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"jsonext"
	"logutil"
	"os"
	"strconv"
)

func init() {
	Getarrayidx_handler(nil, nil, nil)
	Getjson_handler(nil, nil, nil)
	Insertarray_handler(nil, nil, nil)
}

func Getjson_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var path string
	var jsonfile string
	var types string
	var f64 float64
	var ival int
	var vstr string
	var bval bool
	var aval []interface{}
	var mval map[string]interface{}
	var vmap map[string]interface{}
	err = nil
	if ns == nil {
		return
	}
	sarr = ns.GetArray("subnargs")
	if len(sarr) != 3 {
		err = fmt.Errorf("need jsonfile path type")
		return
	}
	jsonfile = sarr[0]
	path = sarr[1]
	types = sarr[2]

	vmap, err = jsonext.GetJson(jsonfile)
	if err != nil {
		return
	}
	switch types {
	case "string":
		vstr, err = jsonext.GetJsonValueString(path, vmap)
	case "float":
		f64, err = jsonext.GetJsonValueFloat(path, vmap)
		if err == nil {
			vstr = fmt.Sprintf("%f", f64)
		}
	case "int":
		ival, err = jsonext.GetJsonValueInt(path, vmap)
		if err == nil {
			vstr = fmt.Sprintf("%d", ival)
		}
	case "array":
		aval, err = jsonext.GetJsonValueArray(path, vmap)
		if err == nil {
			vstr, err = jsonext.FormJsonStruct(aval)
		}

	case "null":
		err = jsonext.GetJsonValueNull(path, vmap)
		if err == nil {
			vstr = "null"
		}
	case "bool":
		bval, err = jsonext.GetJsonValueBool(path, vmap)
		if err == nil {
			if bval {
				vstr = "true"
			} else {
				vstr = "false"
			}
		}
	case "map":
		mval, err = jsonext.GetJsonValueMap(path, vmap)
		if err == nil {
			vstr, err = jsonext.FormJsonStruct(mval)
		}
	default:
		err = fmt.Errorf("[%s] not supported type", types)
		return
	}
	if err != nil {
		return
	}
	fmt.Fprintf(os.Stdout, "[%s] [%s] => [%s]\n", jsonfile, path, vstr)
	err = nil
	return
}

func Getarrayidx_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var path string
	var jsonfile string
	var types string
	var f64 float64
	var ival int
	var vstr string
	var bval bool
	var aval []interface{}
	var mval map[string]interface{}
	var vmap map[string]interface{}
	var varr []interface{}
	var idx int
	err = nil
	if ns == nil {
		return
	}
	sarr = ns.GetArray("subnargs")
	if len(sarr) != 4 {
		err = fmt.Errorf("need jsonfile path type index")
		return
	}
	jsonfile = sarr[0]
	path = sarr[1]
	types = sarr[2]
	idx, err = strconv.Atoi(sarr[3])
	if err != nil {
		return
	}

	vmap, err = jsonext.GetJson(jsonfile)
	if err != nil {
		return
	}

	varr, err = jsonext.GetJsonValueArray(path, vmap)
	if err != nil {
		return
	}

	switch types {
	case "string":
		vstr, err = jsonext.GetJsonArrayItemString(varr, idx)
	case "float":
		f64, err = jsonext.GetJsonArrayItemFloat(varr, idx)
		if err == nil {
			vstr = fmt.Sprintf("%f", f64)
		}
	case "int":
		ival, err = jsonext.GetJsonArrayItemInt(varr, idx)
		if err == nil {
			vstr = fmt.Sprintf("%d", ival)
		}
	case "array":
		aval, err = jsonext.GetJsonArrayItemArray(varr, idx)
		if err == nil {
			vstr, err = jsonext.FormJsonStruct(aval)
		}

	case "null":
		err = jsonext.GetJsonArrayItemNull(varr, idx)
		if err == nil {
			vstr = "null"
		}
	case "bool":
		bval, err = jsonext.GetJsonArrayItemBool(varr, idx)
		if err == nil {
			if bval {
				vstr = "true"
			} else {
				vstr = "false"
			}
		}
	case "map":
		mval, err = jsonext.GetJsonArrayItemMap(varr, idx)
		if err == nil {
			vstr, err = jsonext.FormJsonStruct(mval)
		}
	default:
		err = fmt.Errorf("[%s] not supported type", types)
		return
	}
	if err != nil {
		return
	}
	fmt.Fprintf(os.Stdout, "[%s] [%s][%d] => [%s]\n", jsonfile, path, idx, vstr)
	err = nil
	return
}

func Insertarray_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var path string
	var jsonfile string
	var idx int
	var vmap map[string]interface{}
	var ins string
	var outs string
	var valstr string
	var typestr string
	var output string
	err = nil
	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 3 {
		err = fmt.Errorf("need path idx valstr ")
		return
	}

	jsonfile = ns.GetString("input")
	output = ns.GetString("output")
	path = sarr[0]
	idx, err = strconv.Atoi(sarr[1])
	if err != nil {
		return
	}
	typestr = sarr[2]
	valstr = sarr[3]

	ins, err = fileop.ReadFile(jsonfile)
	if err != nil {
		return
	}

	vmap, err = jsonext.SafeParseMessage(ins)
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	vmap, err = jsonext.InsertJsonArrayItem(path, idx, typestr, valstr, vmap)
	if err != nil {
		return
	}
	outs, err = jsonext.FormatJsonValue(0, "", vmap)
	if err != nil {
		return
	}

	_, err = fileop.WriteFile(output, outs)
	return
}

func main() {
	var parser *extargsparse.ExtArgsParse
	var err error
	var commandline string = `
	{
		"input|i" : null,
		"output|o" : null,
		"getjson<Getjson_handler>##jsonfile path type : type can be int float array map string null bool##" : {
			"$" : 3
		},
		"getarrayidx<Getarrayidx_handler>##jsonfile path type idx to get array index type can be int float array map string null bool##" : {
			"$" : 4
		},
		"insertarray<Insertarray_handler>##path idx typestr valuestr to append json idx -1 means append##" : {
			"$" : 4
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
