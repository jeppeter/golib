package main

/* #include <fcntl.h>
 #include <sys/stat.h>
 #include <semaphore.h>
 #include <string.h>
 #include <stdlib.h>
 #include <errno.h>
 #include <stdio.h>

int _errno() {
    return errno;
}

void _set_errno(int e) {
	errno = e;
	return ;
}

void* _sem_open(char* name, int flags) {
	return sem_open(name,flags,0644,0);
}

void _sem_close(void* ptr){
	if (ptr != NULL) {
		sem_t* sem = (sem_t*) ptr;
		sem_close(sem);
	}
}

void _sem_unlink(char* name) {
	if (name != NULL) {
		sem_unlink(name);
	}
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type SingleLock struct {
	name string
	ptr  uintptr
}

func lock_single(name string) (*SingleLock, error) {
	var ptr uintptr
	cName := C.CString(name)
	flags := C.O_CREAT | C.O_EXCL
	C._set_errno(0)
	ptr = uintptr(C._sem_open(cName, C.int(flags)))
	if C._errno() != 0 {
		return nil, fmt.Errorf("can not make [%s] sema [%d]", name, C._errno())
	}
	sl := &SingleLock{}
	sl.name = name
	sl.ptr = ptr
	return sl, nil
}

func unlock_single(sl *SingleLock) {
	if sl != nil {
		C._sem_close(unsafe.Pointer(sl.ptr))
		cName := C.CString(sl.name)
		C._sem_unlink(cName)
	}

}
