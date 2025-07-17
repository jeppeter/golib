package main

import (
	"dbgutil"
	"executil"
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"logutil"
	"npipepack"
	"os"
	"strings"
	"strop"
	"winpriv"
	"winreg"
)

func init() {
	Readregstring_handler(nil, nil, nil)
	Writeregstring_handler(nil, nil, nil)
	Readregbytes_handler(nil, nil, nil)
	Writeregbytes_handler(nil, nil, nil)
	Createregkey_handler(nil, nil, nil)
	Deleteregkey_handler(nil, nil, nil)
	Deleteregvalue_handler(nil, nil, nil)
	Clearroutetable_handler(nil, nil, nil)
	Setpriv_handler(nil, nil, nil)
	Loadhive_handler(nil, nil, nil)
	Unloadhive_handler(nil, nil, nil)
	Savehive_handler(nil, nil, nil)
	Npsvr_handler(nil, nil, nil)
	Npcli_handler(nil, nil, nil)
	Enumkeys_handler(nil, nil, nil)
	Enumvals_handler(nil, nil, nil)
}

func LoadRegCmdFlags(parser *extargsparse.ExtArgsParse) (err error) {
	var commandline_fmt string
	var commandline string
	commandline_fmt = `{
		"regsubkey" : "SYSTEM\\CurrentControlSet\\Control\\Idvtools\\BTVMTOOL",
		"regpath" : "routetable",
		"regkey" : "HKLM",
		"ReadRegString<Readregstring_handler>## root path key to read registry root can be(HKLM|HKCU|HKCR|HKU|HKCC)##" : {
			"$" : "+"
		},
		"WriteRegString<Writeregstring_handler>## root path key value to write registry root can be(HKLM|HKCU|HKCR|HKU|HKCC)##" : {
			"$" : "+"
		},
		"ReadRegBytes<Readregbytes_handler>## root path key to read registry root can be(HKLM|HKCU|HKCR|HKU|HKCC)##" : {
			"$" : "+"
		},
		"WriteRegBytes<Writeregbytes_handler>## root path key typestr val... to write registry root can be(HKLM|HKCU|HKCR|HKU|HKCC)##" : {
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
		},
		"loadhive<Loadhive_handler>##file subkey to load hive ##" : {
			"$" : 2
		},
		"unloadhive<Unloadhive_handler>##subkey to unload hive##" : {
			"$" : 1
		},
		"savehive<Savehive_handler>##file subkey to save hive##" : {
			"$" : 2
		},
		"npsvr<Npsvr_handler>##pipename to listen on pipe##" : {
			"$" : 1
		},
		"npcli<Npcli_handler>##pipename jsonfile ... to write json and wait##" : {
			"$" : "+"
		},
		"enumkeys<Enumkeys_handler>##root path to enumerate keys name##" : {
			"$" : 2
		},
		"enumvals<Enumvals_handler>##root path to enumerate values##" : {
			"$" : 2
		}

	}`

	commandline = fmt.Sprintf(commandline_fmt)
	err = parser.LoadCommandLineString(commandline)
	return
}

/*ReadRegString*/
func Readregstring_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var args []string
	var valb []byte
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
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
	valb = []byte(value)
	logutil.DebugBuffer(valb, "value byte")
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

/*ReadRegBytes*/
func Readregbytes_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var args []string
	var valb []byte
	var i, lasti int
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
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
	valb, valtype, err := winreg.ReadRegBytes(root, path, key)
	if err != nil {
		err = dbgutil.FormatError("read %s %s\\%s error(%s)", root, path, key, err.Error())
		return
	}
	fmt.Fprintf(os.Stdout, "read %s %s\\%s type [%s]", root, path, key, valtype)
	lasti = 0
	for i, _ = range valb {
		if (i % 16) == 0 {
			if i > 0 {
				fmt.Fprintf(os.Stdout, "   ")
				for lasti < i {
					if valb[lasti] >= byte(' ') && valb[lasti] <= byte('~') {
						fmt.Fprintf(os.Stdout, "%c", valb[lasti])
					} else {
						fmt.Fprintf(os.Stdout, ".")
					}
					lasti += 1
				}
			}
			fmt.Fprintf(os.Stdout, "\n0x%08x", i)
		}
		fmt.Fprintf(os.Stdout, " 0x%02x", valb[i])
	}

	if lasti != i {
		for (i % 16) != 0 {
			fmt.Fprintf(os.Stdout, "     ")
			i += 1
		}

		fmt.Fprintf(os.Stdout, "    ")
		for lasti < len(valb) {
			if valb[lasti] >= byte(' ') && valb[lasti] <= byte('~') {
				fmt.Fprintf(os.Stdout, "%c", valb[lasti])
			} else {
				fmt.Fprintf(os.Stdout, ".")
			}
			lasti += 1
		}
	}
	fmt.Fprintf(os.Stdout, "\n")
	return nil
}

/*WriteRegBytes*/
func Writeregbytes_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var args []string
	var curval uint64
	var vals string
	var valb []byte
	var idx int
	if ns == nil {
		err = nil
		return
	}
	args = ns.GetArray("subnargs")
	if len(args) < 4 {
		err = dbgutil.FormatError("need root path key typestr valb")
		return
	}
	root := args[0]
	path := args[1]
	key := args[2]
	typestr := args[3]
	valb = []byte{}

	for idx = 4; idx < len(args); idx += 1 {
		curval, err = strop.Parseu64(args[idx])
		if err != nil {
			return
		}
		valb = append(valb, byte(curval))
	}

	vals = ""
	for idx = 0; idx < len(valb); idx += 1 {
		if idx > 0 {
			vals += ","
		}
		vals += fmt.Sprintf("0x%02x", valb[idx])
	}

	err = winreg.WriteRegBytes(root, path, key, valb, typestr)
	if err != nil {
		err = dbgutil.FormatError("write %s %s\\%s %s(%s) error(%s)", root, path, key, vals, typestr, err.Error())
		return
	}
	fmt.Fprintf(os.Stdout, "write [%s]%s\\%s %s(%s) succ\n", root, path, key, vals, typestr)
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

func Loadhive_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var subkey string
	var file string
	var root string
	err = nil
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = dbgutil.FormatError("need file subkey")
		return
	}

	root = ns.GetString("regkey")
	subkey = sarr[1]
	file = sarr[0]

	err = winreg.LoadHive(file, root, subkey)
	if err != nil {
		return
	}
	fmt.Printf("load [%s] => [%s].[%s] succ\n", file, root, subkey)

	return nil
}

func Unloadhive_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var subkey string
	var root string
	err = nil
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
		err = dbgutil.FormatError("need subkey")
		return
	}

	root = ns.GetString("regkey")
	subkey = sarr[0]

	err = winreg.UnLoadHive(root, subkey)
	if err != nil {
		return
	}
	fmt.Printf("unload [%s].[%s] succ\n", root, subkey)

	return nil
}

func Savehive_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var subkey string
	var file string
	var root string
	err = nil
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = dbgutil.FormatError("need file subkey")
		return
	}

	root = ns.GetString("regkey")
	subkey = sarr[1]
	file = sarr[0]

	err = winreg.SaveHive(file, root, subkey)
	if err != nil {
		return
	}
	fmt.Printf("saves [%s].[%s] => [%s]  succ\n", root, subkey, file)

	return nil
}

func Npsvr_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var npsvr *npipepack.NpipeSock = nil
	var npacc *npipepack.NpipeSock = nil
	var ndata *npipepack.NpipeData = nil
	var pipename string
	err = nil
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	pipename = sarr[0]

try_bind_again:
	if npacc != nil {
		logutil.Debug("close acc [%s]", pipename)
		npacc.Close()
		npacc = nil
	}

	if npsvr != nil {
		logutil.Debug("close svr [%s]", pipename)
		npsvr.Close()
		npsvr = nil
	}

	npsvr, err = npipepack.BindPipe(pipename, 500)
	if err != nil {
		logutil.Error("can not bind [%s] error [%s]", pipename, err.Error())
		goto try_bind_again
	}

	logutil.Debug("listen on [%s]", pipename)
try_accept:
	npacc, err = npsvr.AcceptTimeout()
	if err != nil {
		logutil.Error("can not accept [%s] error [%s]", pipename, err.Error())
		goto try_bind_again
	}

	if npacc == nil {
		goto try_accept
	}

	logutil.Debug("accept [%s]", pipename)

	for {
		ndata, err = npacc.ReadpacketTimeout()
		if err != nil {
			logutil.Error("read [%s] error [%s]", pipename, err.Error())
			goto try_bind_again
		}

		if ndata == nil {
			continue
		}

		logutil.Debug("read [%s]\n%s", pipename, ndata.GetStr())
		err = npacc.WritePacket(ndata)
		if err != nil {
			logutil.Error("write [%s] error [%s]\n%s", pipename, err.Error(), ndata.GetStr())
			goto try_bind_again
		}
	}

	return
}

func Npcli_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var npcli *npipepack.NpipeSock = nil
	var ndata *npipepack.NpipeData
	var pipename string
	var f string
	var idx int
	var s string
	err = nil
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	pipename = sarr[0]

	npcli, err = npipepack.ConnPipe(pipename, 500, 500)
	if err != nil {
		return
	}
	logutil.Debug("connnect [%s]", pipename)

	defer npcli.Close()

	for idx = 1; idx < len(sarr); idx += 1 {
		f = sarr[idx]
		ndata = npipepack.NewNpipeData()
		s, err = fileop.ReadFile(f)
		if err != nil {
			return
		}
		ndata.SetStr(s)
		logutil.Debug("send [%s]\n%s", pipename, s)
		err = npcli.WritePacket(ndata)
		if err != nil {
			logutil.Error("[%s] write [%s]\n%s", pipename, err.Error(), s)
			return
		}

		for {
			ndata, err = npcli.ReadpacketTimeout()
			if err != nil {
				return
			}
			if ndata == nil {
				continue
			}

			logutil.Debug("read [%s]\n%s", pipename, ndata.GetStr())
			break
		}
	}

	err = nil

	return
}

func Enumkeys_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var keys []string
	var idx int
	var root string
	var path string
	err = nil
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}
	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = dbgutil.FormatError("need root path")
		return
	}
	root = sarr[0]
	path = sarr[1]
	keys, err = winreg.EnumerateRegKeys(root, path)
	if err != nil {
		return
	}

	for idx = 0; idx < len(keys); idx += 1 {
		fmt.Printf("[%s].[%s].[%d] = [%s]\n", root, path, idx, keys[idx])
	}
	err = nil

	return
}

func Enumvals_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var keys []string
	var idx int
	var root string
	var path string
	err = nil
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}
	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = dbgutil.FormatError("need root path")
		return
	}
	root = sarr[0]
	path = sarr[1]
	keys, err = winreg.EnumerateRegValueKeys(root, path)
	if err != nil {
		return
	}

	for idx = 0; idx < len(keys); idx += 1 {
		fmt.Printf("[%s].[%s].[%d]value = [%s]\n", root, path, idx, keys[idx])
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

	err = LoadRegCmdFlags(parser)
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
