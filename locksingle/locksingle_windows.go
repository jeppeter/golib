package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32        = syscall.NewLazyDLL("kernel32.dll")
	procCreateMutex = kernel32.NewProc("CreateMutexW")
	procCloseHandle = kernel32.NewProc("CloseHandle")
)

type SingleLock struct {
	ptr uintptr
}

func lock_single(name string) (*SingleLock, error) {
	ret, _, err := procCreateMutex.Call(
		0,
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name))),
	)
	switch int(err.(syscall.Errno)) {
	case 0:
		sl := &SingleLock{}
		sl.ptr = ret
		return sl, nil
	default:
		return nil, err
	}
}

func unlock_single(sl *SingleLock) {
	if sl != nil {
		procCloseHandle.Call(sl.ptr)
	}
	return
}
