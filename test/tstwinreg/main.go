package main

import (
	"dbgutil"
	"executil"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
	"runtime"
	"strconv"
	"strings"
	"winpriv"
	"winreg"
)

func init() {
	Readregstring_handler(nil, nil, nil)
	Writeregstring_handler(nil, nil, nil)
	Createregkey_handler(nil, nil, nil)
	Deleteregkey_handler(nil, nil, nil)
	Deleteregvalue_handler(nil, nil, nil)
	Clearroutetable_handler(nil, nil, nil)
	Setpriv_handler(nil, nil, nil)
}

func LoadRegCmdFlags(parser *extargsparse.ExtArgsParse) (err error) {
	var commandline_fmt string
	var commandline string
	commandline_fmt = `{
		"regsubkey" : "SYSTEM\\CurrentControlSet\\Control\\Idvtools\\BTVMTOOL",
		"regpath" : "routetable",
		"ReadRegString<Readregstring_handler>## root path key to read registry root can be(HKLM|HKCU|HKCR|HKU|HKCC)##" : {
			"$" : "+"
		},
		"WriteRegString<Writeregstring_handler>## root path key value to write registry root can be(HKLM|HKCU|HKCR|HKU|HKCC)##" : {
			"$" : "+"
		},
		"CreateRegKey<Createregkey_handler>## root path key [accesstype] [existok] to create registry root can be(HKLM|HKCU|HKCR|HKU|HKCC) accesstype can be(ALL|EXECUTE|QUERY_VALUE|READ|SET_VALUE|WRITE|ENUMERATE_SUB_KEYS) default(ALL) existed ok default True##" : {
			"$" : "+"
		},
		"DeleteRegKey<Deleteregkey_handler>## root path to delete registry key root can be(HKLM|HKCU|HKCR|HKU|HKCC)##" : {
			"$" : "+"
		},
		"DeleteRegValue<Deleteregvalue_handler>## root path value to delete registry value root can be(HKLM|HKCU|HKCR|HKU|HKCC)##" : {
			"$" : "+"
		},
		"Clearroutetable<Clearroutetable_handler>## to clear root table default for regsubkey regpath ##"  : {
			"$" : 0
		},
		"setpriv<Setpriv_handler>##privname ... to set priv##" : {
			"$" : "+"
		}
	}`

	commandline = fmt.Sprintf(commandline_fmt)
	err = parser.LoadCommandLineString(commandline)
	return
}

/*ReadRegString*/
func Readregstring_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var args []string
	if ns == nil {
		err = nil
		return
	}
	args = ns.GetArray("subnargs")
	if len(args) < 3 {
		err = dbgutil.FormatError(" need root path key")
		return
	}
	root := args[0]
	path := args[1]
	key := args[2]
	value, valtype, err := winreg.ReadRegString(root, path, key)
	if err != nil {
		err = dbgutil.FormatError("read %s %s\\%s error(%s)", root, path, key, err.Error())
		return
	}
	fmt.Fprintf(os.Stdout, "[%s]%s\\%s=%s(%s)\n", root, path, key, value, valtype)
	return nil
}

/*WriteRegString*/
func Writeregstring_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var args []string
	if ns == nil {
		err = nil
		return
	}
	args = ns.GetArray("subnargs")
	if len(args) < 4 {
		err = dbgutil.FormatError("need root path key value")
		return
	}
	root := args[0]
	path := args[1]
	key := args[2]
	value := args[3]
	typestr := "SZ"
	if len(args) > 4 {
		typestr = args[4]
	}
	err = winreg.WriteRegString(root, path, key, value, typestr)
	if err != nil {
		err = dbgutil.FormatError("write %s %s\\%s %s(%s) error(%s)", root, path, key, value, typestr, err.Error())
		return
	}
	fmt.Fprintf(os.Stdout, "write [%s]%s\\%s %s(%s) succ\n", root, path, key, value, typestr)
	return nil
}

/*CreateRegKey*/
func Createregkey_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var args []string
	if ns == nil {
		err = nil
		return
	}
	args = ns.GetArray("subnargs")
	if len(args) < 2 {
		err = dbgutil.FormatError("root path [accesstype] [existok]")
		return
	}
	root := args[0]
	path := args[1]
	accesstype := "ALL"
	existok := true
	if len(args) > 2 {
		accesstype = args[2]
	}
	if len(args) > 3 {
		existok = false
	}
	err = winreg.CreateRegKey(root, path, accesstype, existok)
	if err != nil {
		err = dbgutil.FormatError("createkey[%s] %s  accesstype(%s) existok(%v) error(%s)", root, path, accesstype, existok, err.Error())
		return
	}
	fmt.Fprintf(os.Stdout, "createkey[%s] %s  accesstype(%s) existok(%v) succ\n", root, path, accesstype, existok)
	return nil
}

/*DeleteRegKey*/
func Deleteregkey_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var args []string
	if ns == nil {
		err = nil
		return
	}
	args = ns.GetArray("subnargs")
	if len(args) < 2 {
		err = dbgutil.FormatError("need root path ")
		return
	}
	root := args[0]
	path := args[1]
	err = winreg.DeleteRegKey(root, path)
	if err != nil {
		err = dbgutil.FormatError("deletekey[%s] %s  error(%s)", root, path, err.Error())
		return
	}
	fmt.Fprintf(os.Stdout, "deletekey[%s] %s succ\n", root, path)
	return nil
}

/*DeleteRegValue*/
func Deleteregvalue_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var args []string
	if ns == nil {
		err = nil
		return
	}
	args = ns.GetArray("subnargs")
	if len(args) < 3 {
		err = dbgutil.FormatError("need root path value")
		return
	}
	root := args[0]
	path := args[1]
	val := args[2]
	err = winreg.DeleteRegValue(root, path, val)
	if err != nil {
		err = dbgutil.FormatError("deletevalue[%s] %s\\%s  error(%s)", root, path, val, err.Error())
		return
	}
	fmt.Fprintf(os.Stdout, "deletekey[%s] %s\\%s succ\n", root, path, val)
	return nil
}

/*Clearroutetable*/
func Clearroutetable_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var root string = "HKLM"
	var path string
	var key string
	var sarr []string
	var carr []string
	var s string
	if ns == nil {
		err = nil
		return
	}

	path = ns.GetString("regsubkey")
	key = ns.GetString("regpath")

	/*now first to2 get value*/
	value, valtype, err := winreg.ReadRegString(root, path, key)
	if err != nil {
		/*nothing to handle the value*/
		err = nil
		return
	}

	if valtype != "SZ" && valtype != "EXPAND_SZ" {
		err = dbgutil.FormatError("[%s].[%s] value type [%s] not valid", path, key, valtype)
		return
	}

	sarr = strings.Split(value, ";")
	if len(sarr) > 0 {
		for _, s = range sarr {
			carr = strings.Split(s, ",")
			if len(carr) >= 3 {
				executil.GetOutputCmd([]string{"route.exe", "delete", carr[0]})
			}
		}
	}

	err = winreg.DeleteRegValue(root, path, key)
	if err != nil {
		return
	}

	return nil
}

func Setpriv_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var i int
	var priv string
	err = nil
	if ns == nil {
		err = nil
		return
	}

	sarr = ns.GetArray("subnargs")
	for i = 0; i < len(sarr); i++ {
		priv = sarr[i]
		err = winpriv.SetPrivLedge(priv, true)
		if err != nil {
			return
		}
		err = winpriv.SetPrivLedge(priv, false)
		if err != nil {
			return
		}
		fmt.Printf("set/unset [%s] succ\n", priv)
	}

	return nil
}

var (
	global_verbose_mode int
)

func get_verbose_mode() int {
	return global_verbose_mode
}

func Set_verbose_callback(ns *extargsparse.NameSpaceEx, validx int, keycls *extargsparse.ExtKeyParse, params []string) (step int, err error) {
	if ns == nil {
		return 0, nil
	}
	if (validx + 1) > len(params) {
		return 0, dbgutil.FormatError("[%d+1] > len(%d) %v", validx, len(params), params)
	}
	global_verbose_mode, err = strconv.Atoi(params[validx])
	if err != nil {
		err = dbgutil.FormatError("parse [%s] not valid number", params[validx])
		return 0, err
	}
	return 1, nil
}

func VerboseLoadFlag(parser *extargsparse.ExtArgsParse) (err error) {
	var commandline string
	commandline = `
	{
		"verbose|V!optparse=Set_verbose_callback!##verbose mode set##" : 0
	}
	`
	err = parser.LoadCommandLineString(commandline)
	return
}

func Error(format string, a ...interface{}) int {
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	s += "\n"
	fmt.Fprint(os.Stderr, s)
	return len(s)
}

func main() {
	var parser *extargsparse.ExtArgsParse
	var err error
	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		Error("%s", err.Error())
		os.Exit(5)
	}

	err = VerboseLoadFlag(parser)
	if err != nil {
		Error("%s", err.Error())
		os.Exit(5)
	}
	err = LoadRegCmdFlags(parser)
	if err != nil {
		Error("%s", err.Error())
		os.Exit(5)
	}

	_, err = parser.ParseCommandLine(nil, nil)
	if err != nil {
		Error("%s", err.Error())
		os.Exit(4)
	}
	os.Exit(0)
}
