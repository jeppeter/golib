package dlfunc

import (
	"unsafe"
)

//#cgo LDFLAGS: -ldl
//#include <dlfcn.h>
//#include <stdlib.h>
import C


type DllLib struct {
	hdl  unsafe.Pointer
	Name string
}

type DllFunc struct {
	libptr   *DllLib
	funcptr  unsafe.Pointer
	procname string
}

func LoadDll(name string) (retptr *DllLib, err error) {

	var libptr unsafe.Pointer
	so_name := C.CString(name)
	defer C.free(unsafe.Pointer(so_name))
	libptr = unsafe.Pointer(C.dlopen(so_name,C.RTLD_LAZY))

	retptr = nil
	if libptr == unsafe.Pointer(0) {
		err = fmt.Errorf("can not load %s", name)
		return
	}

	retptr = &DllLib{}
	retptr.hdl = libptr
	retptr.Name = name
	err = nil
	return
}

func (ptr *DllLib) Close() {
	if ptr.hdl != unsafe.Pointer(0) {
		C.dlclose(ptr.hdl)
	}
	ptr.hdl = unsafe.Pointer(0)
	return
}

func (ptr *DllLib) GetFunc(procname string) (retptr *DllFunc, err error) {
	var funcaddr uintptr
	retptr = nil
	if ptr.hdl == unsafe.Pointer(0) {
		err = fmt.Errorf("%s freed", ptr.Name)
		return
	}
	func_name := C.CString(procname)
	defer C.free(unsafe.Pointer(func_name))

	funcaddr = C.dlsym(ptr.hdl, func_name)
	if funcaddr == uintptr(0) {
		err = fmt.Errorf("can not find %s in %s",procname, ptr.Name)
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

	var inputvar []uintptr

	inputvar = make_input_var(num,a...)

    if num == 0 {
        func_ptr :=  (*func() uintptr) fptr.funcptr
        retval = func_ptr()
    } else if num == 1 {
        func_ptr :=  (*func(*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0])
    } else if num == 2 {
        func_ptr :=  (*func(*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1])
    } else if num == 3 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2])
    } else if num == 4 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3])
    } else if num == 5 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3] ,input_var[4])
    } else if num == 6 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3] ,input_var[4] ,input_var[5])
    } else if num == 7 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3] ,input_var[4] ,input_var[5] ,input_var[6])
    } else if num == 8 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3] ,input_var[4] ,input_var[5] ,input_var[6] ,input_var[7])
    } else if num == 9 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3] ,input_var[4] ,input_var[5] ,input_var[6] ,input_var[7] ,input_var[8])
    } else if num == 10 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3] ,input_var[4] ,input_var[5] ,input_var[6] ,input_var[7] ,input_var[8] ,input_var[9])
    } else if num == 11 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3] ,input_var[4] ,input_var[5] ,input_var[6] ,input_var[7] ,input_var[8] ,input_var[9] ,input_var[10])
    } else if num == 12 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3] ,input_var[4] ,input_var[5] ,input_var[6] ,input_var[7] ,input_var[8] ,input_var[9] ,input_var[10] ,input_var[11])
    } else if num == 13 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3] ,input_var[4] ,input_var[5] ,input_var[6] ,input_var[7] ,input_var[8] ,input_var[9] ,input_var[10] ,input_var[11] ,input_var[12])
    } else if num == 14 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3] ,input_var[4] ,input_var[5] ,input_var[6] ,input_var[7] ,input_var[8] ,input_var[9] ,input_var[10] ,input_var[11] ,input_var[12] ,input_var[13])
    } else if num == 15 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3] ,input_var[4] ,input_var[5] ,input_var[6] ,input_var[7] ,input_var[8] ,input_var[9] ,input_var[10] ,input_var[11] ,input_var[12] ,input_var[13] ,input_var[14])
    } else if num == 16 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3] ,input_var[4] ,input_var[5] ,input_var[6] ,input_var[7] ,input_var[8] ,input_var[9] ,input_var[10] ,input_var[11] ,input_var[12] ,input_var[13] ,input_var[14] ,input_var[15])
    } else if num == 17 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3] ,input_var[4] ,input_var[5] ,input_var[6] ,input_var[7] ,input_var[8] ,input_var[9] ,input_var[10] ,input_var[11] ,input_var[12] ,input_var[13] ,input_var[14] ,input_var[15] ,input_var[16])
    } else if num == 18 {
        func_ptr :=  (*func(*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char ,*C.char) uintptr) fptr.funcptr
        retval = func_ptr(input_var[0] ,input_var[1] ,input_var[2] ,input_var[3] ,input_var[4] ,input_var[5] ,input_var[6] ,input_var[7] ,input_var[8] ,input_var[9] ,input_var[10] ,input_var[11] ,input_var[12] ,input_var[13] ,input_var[14] ,input_var[15] ,input_var[16] ,input_var[17])
    } else {
        err = fmt.Errorf("%%d >= %d",num)
        return
    }

    err = nil
	return
}

