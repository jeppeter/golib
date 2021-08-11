package main

import (
	"os"
	"sync/atomic"
	"syscall"
)

func NewProcWait(proc *os.Process) (retp *ProcWait, err error) {
	const da = syscall.STANDARD_RIGHTS_READ |
		syscall.PROCESS_QUERY_INFORMATION | syscall.SYNCHRONIZE
	var h syscall.Handle
	retp = &ProcWait{}
	retp.pid = proc.Pid
	retp.handle = uintptr(0)
	retp.isdone = 0
	retp.exitcode = 1
	h, err = syscall.OpenProcess(da, false, uint32(retp.pid))
	if err != nil {
		return
	}
	retp.handle = uintptr(h)
	err = nil
	return
}

func (p *ProcWait) IsExited() bool {
	if p.isdone > 0 {
		return true
	}
	return false
}

func (p *ProcWait) GetExitcode() int {
	return p.exitcode
}

func (p *ProcWait) WaitExitTimeout(mills int) bool {
	if p.isdone > 0 {
		return true
	}
	hv := atomic.LoadUintptr(&p.handle)
	s, e := syscall.WaitForSingleObject(syscall.Handle(hv), uint32(mills))
	switch s {
	case syscall.WAIT_OBJECT_0:
		break
	case syscall.WAIT_FAILED:
		return false
	default:
		return false
	}

	var ec uint32
	e = syscall.GetExitCodeProcess(syscall.Handle(p.handle), &ec)
	if e != nil {
		return false
	}
	p.isdone = 1
	p.exitcode = int(ec)

	return true
}

func (p *ProcWait) Close() {
	if p.handle != uintptr(0) {
		syscall.CloseHandle(syscall.Handle(p.handle))
		p.handle = uintptr(0)
	}
	return
}
