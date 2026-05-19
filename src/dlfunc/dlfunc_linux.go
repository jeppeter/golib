package dlfunc

import (
	"fmt"
	"runtime"
	"unsafe"
)

//#cgo LDFLAGS: -ldl
//#include <dlfcn.h>
//#include <stdlib.h>
// typedef unsigned long (*call_0_func_t)();
//
// unsigned long
// bridge_call_0(call_0_func_t f)
// {
//     return f();
// }
//
// typedef unsigned long (*call_1_func_t)(unsigned long);
//
// unsigned long
// bridge_call_1(call_1_func_t f ,unsigned long a0)
// {
//     return f(a0);
// }
//
// typedef unsigned long (*call_2_func_t)(unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_2(call_2_func_t f ,unsigned long a0 ,unsigned long a1)
// {
//     return f(a0,a1);
// }
//
// typedef unsigned long (*call_3_func_t)(unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_3(call_3_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2)
// {
//     return f(a0,a1,a2);
// }
//
// typedef unsigned long (*call_4_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_4(call_4_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3)
// {
//     return f(a0,a1,a2,a3);
// }
//
// typedef unsigned long (*call_5_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_5(call_5_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3 ,unsigned long a4)
// {
//     return f(a0,a1,a2,a3,a4);
// }
//
// typedef unsigned long (*call_6_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_6(call_6_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3 ,unsigned long a4 ,unsigned long a5)
// {
//     return f(a0,a1,a2,a3,a4,a5);
// }
//
// typedef unsigned long (*call_7_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_7(call_7_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3 ,unsigned long a4 ,unsigned long a5 ,unsigned long a6)
// {
//     return f(a0,a1,a2,a3,a4,a5,a6);
// }
//
// typedef unsigned long (*call_8_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_8(call_8_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3 ,unsigned long a4 ,unsigned long a5 ,unsigned long a6 ,unsigned long a7)
// {
//     return f(a0,a1,a2,a3,a4,a5,a6,a7);
// }
//
// typedef unsigned long (*call_9_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_9(call_9_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3 ,unsigned long a4 ,unsigned long a5 ,unsigned long a6 ,unsigned long a7 ,unsigned long a8)
// {
//     return f(a0,a1,a2,a3,a4,a5,a6,a7,a8);
// }
//
// typedef unsigned long (*call_10_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_10(call_10_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3 ,unsigned long a4 ,unsigned long a5 ,unsigned long a6 ,unsigned long a7 ,unsigned long a8 ,unsigned long a9)
// {
//     return f(a0,a1,a2,a3,a4,a5,a6,a7,a8,a9);
// }
//
// typedef unsigned long (*call_11_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_11(call_11_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3 ,unsigned long a4 ,unsigned long a5 ,unsigned long a6 ,unsigned long a7 ,unsigned long a8 ,unsigned long a9 ,unsigned long a10)
// {
//     return f(a0,a1,a2,a3,a4,a5,a6,a7,a8,a9,a10);
// }
//
// typedef unsigned long (*call_12_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_12(call_12_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3 ,unsigned long a4 ,unsigned long a5 ,unsigned long a6 ,unsigned long a7 ,unsigned long a8 ,unsigned long a9 ,unsigned long a10 ,unsigned long a11)
// {
//     return f(a0,a1,a2,a3,a4,a5,a6,a7,a8,a9,a10,a11);
// }
//
// typedef unsigned long (*call_13_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_13(call_13_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3 ,unsigned long a4 ,unsigned long a5 ,unsigned long a6 ,unsigned long a7 ,unsigned long a8 ,unsigned long a9 ,unsigned long a10 ,unsigned long a11 ,unsigned long a12)
// {
//     return f(a0,a1,a2,a3,a4,a5,a6,a7,a8,a9,a10,a11,a12);
// }
//
// typedef unsigned long (*call_14_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_14(call_14_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3 ,unsigned long a4 ,unsigned long a5 ,unsigned long a6 ,unsigned long a7 ,unsigned long a8 ,unsigned long a9 ,unsigned long a10 ,unsigned long a11 ,unsigned long a12 ,unsigned long a13)
// {
//     return f(a0,a1,a2,a3,a4,a5,a6,a7,a8,a9,a10,a11,a12,a13);
// }
//
// typedef unsigned long (*call_15_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_15(call_15_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3 ,unsigned long a4 ,unsigned long a5 ,unsigned long a6 ,unsigned long a7 ,unsigned long a8 ,unsigned long a9 ,unsigned long a10 ,unsigned long a11 ,unsigned long a12 ,unsigned long a13 ,unsigned long a14)
// {
//     return f(a0,a1,a2,a3,a4,a5,a6,a7,a8,a9,a10,a11,a12,a13,a14);
// }
//
// typedef unsigned long (*call_16_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_16(call_16_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3 ,unsigned long a4 ,unsigned long a5 ,unsigned long a6 ,unsigned long a7 ,unsigned long a8 ,unsigned long a9 ,unsigned long a10 ,unsigned long a11 ,unsigned long a12 ,unsigned long a13 ,unsigned long a14 ,unsigned long a15)
// {
//     return f(a0,a1,a2,a3,a4,a5,a6,a7,a8,a9,a10,a11,a12,a13,a14,a15);
// }
//
// typedef unsigned long (*call_17_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_17(call_17_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3 ,unsigned long a4 ,unsigned long a5 ,unsigned long a6 ,unsigned long a7 ,unsigned long a8 ,unsigned long a9 ,unsigned long a10 ,unsigned long a11 ,unsigned long a12 ,unsigned long a13 ,unsigned long a14 ,unsigned long a15 ,unsigned long a16)
// {
//     return f(a0,a1,a2,a3,a4,a5,a6,a7,a8,a9,a10,a11,a12,a13,a14,a15,a16);
// }
//
// typedef unsigned long (*call_18_func_t)(unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long ,unsigned long);
//
// unsigned long
// bridge_call_18(call_18_func_t f ,unsigned long a0 ,unsigned long a1 ,unsigned long a2 ,unsigned long a3 ,unsigned long a4 ,unsigned long a5 ,unsigned long a6 ,unsigned long a7 ,unsigned long a8 ,unsigned long a9 ,unsigned long a10 ,unsigned long a11 ,unsigned long a12 ,unsigned long a13 ,unsigned long a14 ,unsigned long a15 ,unsigned long a16 ,unsigned long a17)
// {
//     return f(a0,a1,a2,a3,a4,a5,a6,a7,a8,a9,a10,a11,a12,a13,a14,a15,a16,a17);
// }
//
import "C"

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
	libptr = unsafe.Pointer(C.dlopen(so_name, C.RTLD_LAZY))

	retptr = nil
	if libptr == unsafe.Pointer(uintptr(0)) {
		err = fmt.Errorf("can not load %s", name)
		return
	}

	retptr = &DllLib{}
	retptr.hdl = libptr
	retptr.Name = name
	runtime.SetFinalizer(retptr, (*DllLib).Close)
	err = nil
	return
}

func (ptr *DllLib) Close() {
	if ptr.hdl != unsafe.Pointer(uintptr(0)) {
		C.dlclose(ptr.hdl)
	}
	ptr.hdl = unsafe.Pointer(uintptr(0))
	return
}

func (ptr *DllLib) GetFunc(procname string) (retptr *DllFunc, err error) {
	var funcaddr unsafe.Pointer
	retptr = nil
	if ptr.hdl == unsafe.Pointer(uintptr(0)) {
		err = fmt.Errorf("%s freed", ptr.Name)
		return
	}
	func_name := C.CString(procname)
	defer C.free(unsafe.Pointer(func_name))

	funcaddr = C.dlsym(ptr.hdl, func_name)
	if funcaddr == unsafe.Pointer(uintptr(0)) {
		err = fmt.Errorf("can not find %s in %s", procname, ptr.Name)
		return
	}

	retptr = &DllFunc{}
	retptr.funcptr = funcaddr
	retptr.libptr = ptr
	retptr.procname = procname
	runtime.SetFinalizer(retptr, (*DllFunc).Close)
	err = nil
	return
}

func make_input_var(num int, a ...uintptr) (retvar []C.ulong) {
	retvar = []C.ulong{}
	for _, ai := range a {
		retvar = append(retvar, C.ulong(ai))
	}

	for len(retvar) < num {
		retvar = append(retvar, C.ulong(0))
	}
	return
}

func (fptr *DllFunc) Close() {
	fptr.libptr = nil
	fptr.funcptr = nil
	fptr.procname = ""
	return
}

func (fptr *DllFunc) Name() string {
	return fptr.procname
}

func (fptr *DllFunc) CallN(num int, a ...uintptr) (ret uintptr, err error) {

	var input_var []C.ulong
	var retval C.ulong

	input_var = make_input_var(num, a...)

	if num == 0 {
		func_ptr := C.call_0_func_t(fptr.funcptr)
		retval = C.bridge_call_0(func_ptr)
	} else if num == 1 {
		func_ptr := C.call_1_func_t(fptr.funcptr)
		retval = C.bridge_call_1(func_ptr, input_var[0])
	} else if num == 2 {
		func_ptr := C.call_2_func_t(fptr.funcptr)
		retval = C.bridge_call_2(func_ptr, input_var[0], input_var[1])
	} else if num == 3 {
		func_ptr := C.call_3_func_t(fptr.funcptr)
		retval = C.bridge_call_3(func_ptr, input_var[0], input_var[1], input_var[2])
	} else if num == 4 {
		func_ptr := C.call_4_func_t(fptr.funcptr)
		retval = C.bridge_call_4(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3])
	} else if num == 5 {
		func_ptr := C.call_5_func_t(fptr.funcptr)
		retval = C.bridge_call_5(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3], input_var[4])
	} else if num == 6 {
		func_ptr := C.call_6_func_t(fptr.funcptr)
		retval = C.bridge_call_6(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3], input_var[4], input_var[5])
	} else if num == 7 {
		func_ptr := C.call_7_func_t(fptr.funcptr)
		retval = C.bridge_call_7(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3], input_var[4], input_var[5], input_var[6])
	} else if num == 8 {
		func_ptr := C.call_8_func_t(fptr.funcptr)
		retval = C.bridge_call_8(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3], input_var[4], input_var[5], input_var[6], input_var[7])
	} else if num == 9 {
		func_ptr := C.call_9_func_t(fptr.funcptr)
		retval = C.bridge_call_9(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3], input_var[4], input_var[5], input_var[6], input_var[7], input_var[8])
	} else if num == 10 {
		func_ptr := C.call_10_func_t(fptr.funcptr)
		retval = C.bridge_call_10(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3], input_var[4], input_var[5], input_var[6], input_var[7], input_var[8], input_var[9])
	} else if num == 11 {
		func_ptr := C.call_11_func_t(fptr.funcptr)
		retval = C.bridge_call_11(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3], input_var[4], input_var[5], input_var[6], input_var[7], input_var[8], input_var[9], input_var[10])
	} else if num == 12 {
		func_ptr := C.call_12_func_t(fptr.funcptr)
		retval = C.bridge_call_12(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3], input_var[4], input_var[5], input_var[6], input_var[7], input_var[8], input_var[9], input_var[10], input_var[11])
	} else if num == 13 {
		func_ptr := C.call_13_func_t(fptr.funcptr)
		retval = C.bridge_call_13(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3], input_var[4], input_var[5], input_var[6], input_var[7], input_var[8], input_var[9], input_var[10], input_var[11], input_var[12])
	} else if num == 14 {
		func_ptr := C.call_14_func_t(fptr.funcptr)
		retval = C.bridge_call_14(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3], input_var[4], input_var[5], input_var[6], input_var[7], input_var[8], input_var[9], input_var[10], input_var[11], input_var[12], input_var[13])
	} else if num == 15 {
		func_ptr := C.call_15_func_t(fptr.funcptr)
		retval = C.bridge_call_15(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3], input_var[4], input_var[5], input_var[6], input_var[7], input_var[8], input_var[9], input_var[10], input_var[11], input_var[12], input_var[13], input_var[14])
	} else if num == 16 {
		func_ptr := C.call_16_func_t(fptr.funcptr)
		retval = C.bridge_call_16(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3], input_var[4], input_var[5], input_var[6], input_var[7], input_var[8], input_var[9], input_var[10], input_var[11], input_var[12], input_var[13], input_var[14], input_var[15])
	} else if num == 17 {
		func_ptr := C.call_17_func_t(fptr.funcptr)
		retval = C.bridge_call_17(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3], input_var[4], input_var[5], input_var[6], input_var[7], input_var[8], input_var[9], input_var[10], input_var[11], input_var[12], input_var[13], input_var[14], input_var[15], input_var[16])
	} else if num == 18 {
		func_ptr := C.call_18_func_t(fptr.funcptr)
		retval = C.bridge_call_18(func_ptr, input_var[0], input_var[1], input_var[2], input_var[3], input_var[4], input_var[5], input_var[6], input_var[7], input_var[8], input_var[9], input_var[10], input_var[11], input_var[12], input_var[13], input_var[14], input_var[15], input_var[16], input_var[17])
	} else {
		err = fmt.Errorf("%d extend 18", num)
		return
	}

	ret = uintptr(retval)
	err = nil
	return
}
