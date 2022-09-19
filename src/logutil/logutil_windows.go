package logutil

import (
	"syscall"
	"unsafe"
)

var __syscall_outputdebugstring *syscall.Proc = nil

func init() {
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

type winbackGround struct {
}

func nativeGround() *winbackGround {
	p := &winbackGround{}
	return p
}

func (pb *winbackGround) LogDebugOutputBackGround(s string) error {
	if __syscall_outputdebugstring != nil {
		p := syscall.StringToUTF16Ptr(s)
		__syscall_outputdebugstring.Call(uintptr(unsafe.Pointer(p)))
	}
	return nil
}

func (pb *winbackGround) CloseDebugOutputBackGround() error {
	if __syscall_outputdebugstring != nil {
		__syscall_outputdebugstring = nil
	}
	return nil
}
