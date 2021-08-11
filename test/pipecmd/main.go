package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"golang.org/x/sys/windows"
	"io/ioutil"
	"os"
	"syscall"
	"time"
)

func init() {
	Pipecmd_handler(nil, nil, nil)
}

const (
	DEBUG_OUT_BYTES = 8192
)

var gl_exitmode int = 0

func pipecmd_go_func(ns *extargsparse.NameSpaceEx, exitch chan int, exitout chan int) {
	var input string
	var inputbytes []byte
	var outbytes, errbytes []byte
	var exitcode int = 1
	var err error
	var sarr []string
	var fin *os.File
	var outlen int

	defer func() {
		Error("exit %d", exitcode)
		exitout <- exitcode
	}()

	inputbytes = []byte{}
	input = ns.GetString("input")
	if len(input) > 0 {
		fin, err = os.Open(input)
		if err != nil {
			Error("open [%s] error[%s]", input, err.Error())
			return
		}
		defer fin.Close()
		inputbytes, err = ioutil.ReadAll(fin)
		if err != nil {
			err = fmt.Errorf("read [%s] [%s]", input, err.Error())
			Error("%s", err.Error())
			return
		}
	}
	Error("inputbytes [%d]", len(inputbytes))

	sarr = ns.GetArray("subnargs")

	outbytes, errbytes, exitcode, err = run_cmd_output(exitch, sarr, inputbytes)
	if err != nil {
		Error("%s", err.Error())
		return
	}

	Error("run cmd %v exitcode %d output [%d] errout [%d]", sarr, exitcode, len(outbytes), len(errbytes))
	outlen = len(outbytes)
	if outlen < DEBUG_OUT_BYTES {
		DebugBuffer(outbytes, "run cmd %v output", sarr)
	} else {

		DebugBuffer(outbytes[:DEBUG_OUT_BYTES], "run cmd %v output first %d", sarr, DEBUG_OUT_BYTES)
		DebugBuffer(outbytes[(outlen-DEBUG_OUT_BYTES):], "run cmd %v output last %d", sarr, DEBUG_OUT_BYTES)
	}

	outlen = len(errbytes)
	if outlen < DEBUG_OUT_BYTES {
		DebugBuffer(errbytes, "run cmd %v errout", sarr)
	} else {

		DebugBuffer(errbytes[:DEBUG_OUT_BYTES], "run cmd %v errout first %d", sarr, DEBUG_OUT_BYTES)
		DebugBuffer(errbytes[(outlen-DEBUG_OUT_BYTES):], "run cmd %v errout last %d", sarr, DEBUG_OUT_BYTES)
	}
	return
}

func pipecmd_go(ns *extargsparse.NameSpaceEx) (err error) {
	var exitch, exitout chan int

	exitch = make(chan int, 10)
	exitout = make(chan int, 10)
	go pipecmd_go_func(ns, exitch, exitout)

	for gl_exitmode == 0 {
		select {
		case <-exitout:
			err = nil
			Error("get exitout")
			return
		case <-time.After(time.Duration(300) * time.Millisecond):
			err = nil
		}
	}

	Error("gl_exitmode %d", gl_exitmode)
	exitch <- 1
	Error("wait exitout")
	<-exitout
	Error("wait exitout over")
	err = nil

	return

}

func Pipecmd_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	err = nil
	if ns == nil {
		return
	}

	err = InitLog(ns)
	if err != nil {
		return
	}

	kernel32 := windows.NewLazySystemDLL("kernel32.dll")

	setConsoleCtrlHandler := kernel32.NewProc("SetConsoleCtrlHandler")

	setConsoleCtrlHandler.Call(
		syscall.NewCallback(func(controlType uint) uint {
			gl_exitmode = 1
			return 1
		}),
		1)

	return pipecmd_go(ns)
}

func Pipecmd_Load(parser *extargsparse.ExtArgsParse) (err error) {
	var commandline string
	commandline = `{
		"input|i" : null,
		"pipecmd<Pipecmd_handler>## args ... to run command##" : {
			"$" : "+"
		}
	}`
	return parser.LoadCommandLineString(commandline)
}

func main() {
	var parser *extargsparse.ExtArgsParse
	var err error

	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		Error("%s", err.Error())
		atexit.Exit(4)
	}

	err = PrepareLog(parser)
	if err != nil {
		Error("%s", err.Error())
		atexit.Exit(4)
	}

	err = Pipecmd_Load(parser)
	if err != nil {
		Error("%s", err.Error())
		atexit.Exit(4)
	}

	_, err = parser.ParseCommandLine(nil, nil)
	if err != nil {
		Error("%s", err.Error())
		atexit.Exit(5)
	}

	Debug("[%d]all succ", os.Getpid())
	atexit.Exit(0)
	return
}
