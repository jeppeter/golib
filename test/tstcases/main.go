package main

import (
	"bytes"
	"dbgutil"
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"logutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func go_chan(c chan string, e chan int, timeout int) {
	var s string
	var icnt int = 0
	for s = range c {
		fmt.Fprintf(os.Stdout, "%s\n", s)
		icnt++
		time.Sleep(time.Duration(timeout) * time.Millisecond)
	}
	e <- icnt
	return
}

func Chan_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var s string
	var i int
	var schan chan string
	var ichan chan int
	err = nil
	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	schan = make(chan string, 100)
	ichan = make(chan int, 1)
	go go_chan(schan, ichan, ns.GetInt("timeout"))

	for _, s = range ns.GetArray("subnargs") {
		schan <- s
	}
	close(schan)
	schan = nil
	i = <-ichan

	fmt.Fprintf(os.Stdout, "end value [%d]\n", i)

	return
}

func get_code(sarr []string) (retbyte []byte, err error) {
	var s string
	var v int64
	var base int
	retbyte = []byte{}
	err = nil

	for _, s = range sarr {
		base = 10
		if strings.HasPrefix(s, "0x") ||
			strings.HasPrefix(s, "0X") {
			base = 16
			s = s[2:]
		} else if strings.HasPrefix(s, "x") ||
			strings.HasPrefix(s, "X") {
			base = 16
			s = s[1:]
		}
		v, err = strconv.ParseInt(s, base, 64)
		if err != nil {
			return
		}
		retbyte = append(retbyte, byte(v))
	}

	err = nil
	return
}

func gbk_to_utf8(inbytes []byte) (outbytes []byte, err error) {
	var rd *transform.Reader
	rd = transform.NewReader(bytes.NewReader(inbytes), simplifiedchinese.GBK.NewDecoder())
	outbytes, err = ioutil.ReadAll(rd)
	if err != nil {
		return
	}
	return
}

func utf8_to_gbk(inbytes []byte) (outbytes []byte, err error) {
	var rd *transform.Reader
	rd = transform.NewReader(bytes.NewReader(inbytes), simplifiedchinese.GBK.NewEncoder())
	outbytes, err = ioutil.ReadAll(rd)
	if err != nil {
		return
	}
	return
}

func utf8_to_uni(inbytes []byte) (outbytes []byte, err error) {
	var idx, retn, i, curval int
	var r rune
	outbytes = []byte{}
	err = nil

	idx = 0
	for idx < len(inbytes) {
		r, retn = utf8.DecodeRune(inbytes[idx:])
		for i = 0; i < 2; i++ {
			curval = int(r)
			curval = (curval >> (i * 8)) & 0xff
			outbytes = append(outbytes, byte(curval))
		}
		idx += retn
	}
	err = nil
	return
}

func uni_to_utf8(inbytes []byte) (outbytes []byte, err error) {
	var s string
	var idx, j, retn int
	var r rune
	var rs []rune
	var ps string
	var buf []byte
	var curval int
	outbytes = []byte{}
	err = nil

	ps = "\""
	for idx = 0; idx < (len(inbytes) - 1); idx += 2 {
		curval = 0
		curval += int(inbytes[idx])
		curval += (int(inbytes[idx+1]) << 8)
		ps += fmt.Sprintf("\\u%04x", curval)
	}
	ps += "\""
	s, err = strconv.Unquote(ps)
	if err != nil {
		err = dbgutil.FormatError("[%s] error [%s]", ps, err.Error())
		return
	}

	buf = make([]byte, 10)

	idx = 0
	rs = []rune(s)
	for idx = 0; idx < len(rs); idx++ {
		r = rs[idx]
		retn = utf8.EncodeRune(buf, r)
		for j = 0; j < retn; j++ {
			outbytes = append(outbytes, buf[j])
		}
	}

	err = nil
	return
}

func out_bytes(inbytes []byte, fmtstr string, a ...interface{}) (outs string) {
	var lasti, i int
	var b byte
	outs = ""
	lasti = 0
	i = 0
	outs += fmt.Sprintf("bytes [%d:0x%x] ", len(inbytes), len(inbytes))
	outs += fmt.Sprintf(fmtstr, a...)
	for i, b = range inbytes {
		if (i % 16) == 0 {
			if i > 0 {
				outs += "    "
				for lasti != i {
					if inbytes[lasti] >= byte(' ') &&
						inbytes[lasti] <= byte('~') {
						outs += fmt.Sprintf("%c", inbytes[lasti])
					} else {
						outs += "."
					}
					lasti++
				}
			}
			outs += fmt.Sprintf("\n0x%08x:", i)
		}
		outs += fmt.Sprintf(" 0x%02x", b)
	}

	if lasti != i {
		for (i % 16) != 0 {
			outs += fmt.Sprintf("     ")
			i++
		}
		outs += fmt.Sprintf("    ")

		for lasti != len(inbytes) {
			if inbytes[lasti] >= byte(' ') &&
				inbytes[lasti] <= byte('~') {
				outs += fmt.Sprintf("%c", inbytes[lasti])
			} else {
				outs += "."
			}
			lasti++
		}

		outs += fmt.Sprintf("\n")
	}

	return
}

func Gbktoutf8_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var inbytes, outbytes []byte
	err = nil
	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	inbytes, err = get_code(ns.GetArray("subnargs"))
	if err != nil {
		return
	}

	logutil.Debug("inbyte %v", inbytes)

	outbytes, err = gbk_to_utf8(inbytes)
	if err != nil {
		return
	}

	fmt.Fprintf(os.Stdout, "%s", out_bytes(inbytes, "input bytes"))
	fmt.Fprintf(os.Stdout, "%s", out_bytes(outbytes, "output bytes"))
	err = nil

	return
}

func Utf8togbk_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var inbytes, outbytes []byte
	err = nil
	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	inbytes, err = get_code(ns.GetArray("subnargs"))
	if err != nil {
		return
	}

	outbytes, err = utf8_to_gbk(inbytes)
	if err != nil {
		return
	}

	fmt.Fprintf(os.Stdout, "%s", out_bytes(inbytes, "input bytes"))
	fmt.Fprintf(os.Stdout, "%s", out_bytes(outbytes, "output bytes"))
	err = nil

	return
}

func Utf8touni_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var inbytes, outbytes []byte
	err = nil
	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	inbytes, err = get_code(ns.GetArray("subnargs"))
	if err != nil {
		return
	}

	outbytes, err = utf8_to_uni(inbytes)
	if err != nil {
		return
	}

	fmt.Fprintf(os.Stdout, "%s", out_bytes(inbytes, "input bytes"))
	fmt.Fprintf(os.Stdout, "%s", out_bytes(outbytes, "output bytes"))
	err = nil

	return
}

func Unitoutf8_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var inbytes, outbytes []byte
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	inbytes, err = get_code(ns.GetArray("subnargs"))
	if err != nil {
		return
	}

	outbytes, err = uni_to_utf8(inbytes)
	if err != nil {
		return
	}

	fmt.Fprintf(os.Stdout, "%s", out_bytes(inbytes, "input bytes"))
	fmt.Fprintf(os.Stdout, "%s", out_bytes(outbytes, "output bytes"))
	err = nil

	return
}

func Readfilebyte_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var outbytes []byte
	var sarr []string
	var i int
	var s string
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) == 0 {
		outbytes, err = fileop.ReadFileBytes("")
		if err != nil {
			return
		}
		logutil.DebugBuffer(outbytes, "read stdin")
	} else {
		for i, s = range sarr {
			outbytes, err = fileop.ReadFileBytes(s)
			if err != nil {
				return
			}
			logutil.DebugBuffer(outbytes, "read [%d][%s]", i, s)
		}
	}
	err = nil
	return
}

func Writefilebyte_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var output string = ""
	var s string
	var outs string
	var nret int
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	output = ns.GetString("output")
	outs = ""
	for _, s = range sarr {
		outs += fmt.Sprintf("%s\n", s)
	}

	nret, err = fileop.WriteFileBytes(output, []byte(outs))
	if err != nil {
		return
	}
	fmt.Fprintf(os.Stderr, "write [%s] nret [%d]\n", output, nret)
	return
}

func Readfile_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var outs string
	var sarr []string
	var i int
	var s string
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) == 0 {
		outs, err = fileop.ReadFile("")
		if err != nil {
			return
		}
		logutil.Debug("read file [stdin]\n%s", outs)
	} else {
		for i, s = range sarr {
			outs, err = fileop.ReadFile(s)
			if err != nil {
				return
			}
			logutil.Debug("read [%d][%s]\n%s", i, s, outs)
		}
	}
	err = nil
	return
}

func Writefile_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var output string = ""
	var s string
	var outs string
	var nret int
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	output = ns.GetString("output")
	outs = ""
	for _, s = range sarr {
		outs += fmt.Sprintf("%s\n", s)
	}

	nret, err = fileop.WriteFile(output, outs)
	if err != nil {
		return
	}
	fmt.Fprintf(os.Stderr, "write [%s] nret [%d]\n", output, nret)
	return
}

func Deletefile_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var s string
	var i int
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	for i, s = range sarr {
		err = fileop.DeleteFile(s)
		if err != nil {
			err = fmt.Errorf("delete [%d].[%s] error[%s]", i, s, err.Error())
			return
		}
		fmt.Printf("delete [%s] succ\n", s)
	}

	return
}

func Parseu64_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var i int
	var retv uint64
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
		retv, err = strconv.ParseUint(sarr[i], 10, 64)
		if err != nil {
			err = dbgutil.FormatError("parse [%s] error[%s]", sarr[i], err.Error())
			return
		}
		fmt.Printf("[%s]=[%d]\n", sarr[i], retv)
	}
	err = nil
	return
}

func Mkdirsafe_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var i int
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
		err = fileop.MkdirSafe(sarr[i], -1)
		if err != nil {
			return
		}
		fmt.Printf("create [%s] succ\n", sarr[i])
	}
	err = nil
	return
}

type MemInfo struct {
	Startaddr uintptr
	Endaddr   uintptr
}

func get_func_addr(name string, startaddr uintptr, endaddr uintptr) (findptr *runtime.Func, err error) {
	//var names []string
	//var searchaddr []uintptr
	var curaddr uintptr
	var stepaddr uintptr = (1 << 10)
	var curfunc *runtime.Func = nil
	findptr = nil
	err = fmt.Errorf("not foud %s", name)
	//names = strings.Split(name, ".")
	for findptr == nil && stepaddr >= 32 {
		curaddr = startaddr
		for {
			curfunc = runtime.FuncForPC(curaddr)
			if curfunc != nil {
				fmt.Printf("0x%x addr %s\n", curaddr, curfunc.Name())
				if curfunc.Name() == name {
					findptr = curfunc
					fmt.Printf("get %s addr 0x%x\n", name, curfunc.Entry())
					err = nil
					break
				}
			}
			curaddr += stepaddr
			if curaddr > endaddr {
				break
			}
		}
		stepaddr = stepaddr >> 1
	}
	return

}

func Querymem_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var meminfo []MemInfo
	var i int
	var info MemInfo
	err = nil

	if ns == nil {
		return
	}
	meminfo, err = get_current_process_exec_info(os.Getpid())
	if err != nil {
		return
	}
	for i, info = range meminfo {
		fmt.Printf("[%d].start 0x%x end 0x%x\n", i, info.Startaddr, info.Endaddr)
	}
	_, err = get_func_addr("main.Querymem_handler", meminfo[0].Startaddr, meminfo[0].Endaddr)

	return
}

func init() {
	Chan_handler(nil, nil, nil)
	Utf8togbk_handler(nil, nil, nil)
	Gbktoutf8_handler(nil, nil, nil)
	Utf8touni_handler(nil, nil, nil)
	Unitoutf8_handler(nil, nil, nil)
	Readfilebyte_handler(nil, nil, nil)
	Writefilebyte_handler(nil, nil, nil)
	Readfile_handler(nil, nil, nil)
	Writefile_handler(nil, nil, nil)
	Deletefile_handler(nil, nil, nil)
	Parseu64_handler(nil, nil, nil)
	Mkdirsafe_handler(nil, nil, nil)
	Goversioncheck_handler(nil, nil, nil)
	Querymem_handler(nil, nil, nil)
}

func Goversioncheck_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	err = nil

	if ns == nil {
		return
	}

	fmt.Printf("version %s\n", runtime.Version())
	return
}

func main() {
	var commandline string
	var err error
	var parser *extargsparse.ExtArgsParse
	var ns *extargsparse.NameSpaceEx

	commandline = `{
		"timeout|t" : 500,
		"input|i" : null,
		"output|o" : null,
		"chan<Chan_handler>##outstr ... to set out string##" : {
			"$" : "+"
		},
		"utf8togbk<Utf8togbk_handler>## codes ... to get codes from utf-8 to ansi##"  : {
			"$" : "+"
		},
		"gbktoutf8<Gbktoutf8_handler>## codes ... to get codes from ansi to utf-8##" : {
			"$" : "+"
		},
		"utf8touni<Utf8touni_handler>## codes ... to get codes from utf-8 to unicode##" : {
			"$" : "+"
		},
		"unitoutf8<Unitoutf8_handler>## codes ... to get codes from utf-8 to unicode##" : {
			"$" : "+"
		},
		"readfilebyte<Readfilebyte_handler>## [fname] ... to read file default stdin ##" : {
			"$" : "*"
		},
		"writefilebyte<Writefilebyte_handler>## strs ... to write file output default stdout##" : {
			"$" : "+"
		},
		"readfile<Readfile_handler>## [fname] ... to read file default stdin ##" : {
			"$" : "*"
		},
		"writefile<Writefile_handler>## strs ... to write file output default stdout##" : {
			"$" : "+"
		},
		"deletefile<Deletefile_handler>## fname ... to delete file##" : {
			"$" : "+"
		},
		"parseu64<Parseu64_handler>##val ... to parse u64##" : {
			"$" : "+"
		},
		"mkdirsafe<Mkdirsafe_handler>##dir ... to make dir safe##" : {
			"$" : "+"
		},
		"goversioncheck<Goversioncheck_handler>##to check go compiler version##" : {
			"$" : 0
		},
		"querymem<Querymem_handler>##to list current process memory##" : {
			"$" : 0
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
	fmt.Fprintf(os.Stdout, "subcommand [%s] succ\n", ns.GetString("subcommand"))
	atexit.Exit(0)
	return
}
