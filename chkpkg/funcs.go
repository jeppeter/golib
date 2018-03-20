package main


import (
	"fmt"
	"syscall"
	"flag"
	"strconv"
	"unsafe"
	"runtime"
	"os"
)


var __syscall_outputdebugstring *syscall.Proc

type verbosemode struct {
	vmode int
}

var (
	global_verbose_mode verbosemode
)

func (v *verbosemode) String() string {
	return fmt.Sprintf("%d", v.vmode)
}
func (v *verbosemode) Set(val string) error {
	ival, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	v.vmode = ival
	return nil
}


func init() {
	flag.Var(&global_verbose_mode, "verbose", "specify verbose mode 3 Debug 2 Info default Error")
	flag.Var(&global_verbose_mode, "V", "specify verbose mode 3 Debug 2 Info default Error")
	__syscall_outputdebugstring = nil
	d, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return
	}
	__syscall_outputdebugstring, err = d.FindProc("OutputDebugStringW")
	if err != nil {
		__syscall_outputdebugstring = nil
		return
	}
	return
}

func InnerDebugOutput(s string) error {
	if __syscall_outputdebugstring != nil {
		p := syscall.StringToUTF16Ptr(s)
		__syscall_outputdebugstring.Call(uintptr(unsafe.Pointer(p)))
	}
	return nil
}



func Info(format string, a ...interface{}) int {
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	s += "\n"
	if global_verbose_mode.vmode > 1 {
		fmt.Fprint(os.Stdout, s)
		InnerDebugOutput(s)
	}
	return len(s)
}

func Debug(format string, a ...interface{}) int {
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	s += "\n"
	if global_verbose_mode.vmode > 2 {
		fmt.Fprint(os.Stdout, s)
		InnerDebugOutput(s)
	}
	return len(s)
}

func Error(format string, a ...interface{}) int {
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	s += "\n"
	fmt.Fprint(os.Stderr, s)
	InnerDebugOutput(s)
	return len(s)
}
