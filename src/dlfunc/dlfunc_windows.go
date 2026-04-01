package dlfunc

import (
	"fmt"
	"syscall"
	"unsafe"
)

type DllLib struct {
	hdl  syscall.Handle
	Name string
}

type DllFunc struct {
	libptr   *DllLib
	funcptr  unsafe.Pointer
	procname string
}

func LoadDll(name string) (retptr *DllLib, err error) {
	var hdl syscall.Handle
	retptr = nil
	hdl, err = syscall.LoadLibrary(name)
	if err != nil {
		return
	}

	retptr = &DllLib{}
	retptr.hdl = hdl
	retptr.Name = name

	err = nil
	return
}

func (ptr *DllLib) Close() {
	if ptr.hdl != syscall.Handle(0) {
		syscall.FreeLibrary(ptr.hdl)
	}
	ptr.hdl = syscall.Handle(0)
	return
}

func (ptr *DllLib) GetFunc(procname string) (retptr *DllFunc, err error) {
	var funcaddr uintptr
	retptr = nil
	if ptr.hdl == syscall.Handle(0) {
		err = fmt.Errorf("%s freed", ptr.Name)
		return
	}
	funcaddr, err = syscall.GetProcAddress(ptr.hdl, procname)
	if err != nil {
		return
	}
	retptr = &DllFunc{}
	retptr.funcptr = unsafe.Pointer(funcaddr)
	retptr.libptr = ptr
	retptr.procname = procname
	err = nil
	return
}

func make_input_var(num int, a ...uintptr) (retvar []uintptr) {
	retvar = []uintptr{}
	for _, ai := range a {
		retvar = append(retvar, uintptr(ai))
	}

	for len(retvar) < num {
		retvar = append(retvar, uintptr(0))
	}
	return
}

func (fptr *DllFunc) Name() string {
	return fptr.procname
}

func (fptr *DllFunc) CallN(num int, a ...uintptr) (retval uintptr, err error) {

	var errv syscall.Errno
	var inputvar []uintptr

	if num <= 3 {
		inputvar = make_input_var(3, a...)
		retval, _, errv = syscall.Syscall(uintptr(fptr.funcptr), uintptr(num), inputvar[0], inputvar[1], inputvar[2])
	} else if num <= 6 {
		inputvar = make_input_var(6, a...)
		retval, _, errv = syscall.Syscall6(uintptr(fptr.funcptr), uintptr(num), inputvar[0], inputvar[1], inputvar[2], inputvar[3], inputvar[4], inputvar[5])
	} else if num <= 9 {
		inputvar = make_input_var(9, a...)
		retval, _, errv = syscall.Syscall9(uintptr(fptr.funcptr), uintptr(num), inputvar[0], inputvar[1], inputvar[2], inputvar[3], inputvar[4], inputvar[5], inputvar[6], inputvar[7], inputvar[8])
	} else if num <= 12 {
		inputvar = make_input_var(12, a...)
		retval, _, errv = syscall.Syscall12(uintptr(fptr.funcptr), uintptr(num), inputvar[0], inputvar[1], inputvar[2], inputvar[3], inputvar[4], inputvar[5], inputvar[6], inputvar[7], inputvar[8], inputvar[9], inputvar[10], inputvar[11])
	} else if num <= 15 {
		inputvar = make_input_var(15, a...)
		retval, _, errv = syscall.Syscall15(uintptr(fptr.funcptr), uintptr(num), inputvar[0], inputvar[1], inputvar[2], inputvar[3], inputvar[4], inputvar[5], inputvar[6], inputvar[7], inputvar[8], inputvar[9], inputvar[10], inputvar[11], inputvar[12], inputvar[13], inputvar[14])
	} else if num <= 18 {
		inputvar = make_input_var(18, a...)
		retval, _, errv = syscall.Syscall18(uintptr(fptr.funcptr), uintptr(num), inputvar[0], inputvar[1], inputvar[2], inputvar[3], inputvar[4], inputvar[5], inputvar[6], inputvar[7], inputvar[8], inputvar[9], inputvar[10], inputvar[11], inputvar[12], inputvar[13], inputvar[14], inputvar[15], inputvar[16], inputvar[17])
	} else {
		err = fmt.Errorf("num %d > 18", num)
		return
	}

	if errv != 0 {
		err = fmt.Errorf("%v", errv)
		return
	}
	err = nil
	return
}
