package executil

import (
	"bytes"
	"dbgutil"
	"fmt"
	"logutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func RunCmdTimeout(cmds []string, timeout int) (outstr string, errstr string, exitcode int, err error) {
	var stime time.Time
	var etime time.Time
	var ctime time.Time
	var outb bytes.Buffer
	var errb bytes.Buffer
	var hdl syscall.Handle = syscall.InvalidHandle
	var evt uint32
	var cmd *exec.Cmd
	var waiting int = 1
	const da = syscall.STANDARD_RIGHTS_READ | syscall.PROCESS_QUERY_INFORMATION | syscall.SYNCHRONIZE
	logutil.Trace("cmds %v", cmds)
	cmd = &exec.Cmd{}
	cmd.Path = cmds[0]
	cmd.Args = cmds
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err = cmd.Start()
	if err != nil {
		err = dbgutil.FormatError("can not run %v error[%s]", cmds, err.Error())
		return
	}

	defer func() {
		if hdl != syscall.InvalidHandle {
			syscall.CloseHandle(hdl)
			hdl = syscall.InvalidHandle
		}
		cmd.Process.Kill()
		cmd = nil
	}()

	hdl, err = syscall.OpenProcess(da, false, uint32(cmd.Process.Pid))
	if err != nil {
		err = dbgutil.FormatError("can not open %v pid [%d] error[%s]", cmds, cmd.Process.Pid, err.Error())
		return
	}
	stime = time.Now()
	etime = stime.Add(time.Duration(timeout) * time.Millisecond)
	waiting = 1
	for waiting > 0 {
		ctime = time.Now()
		if ctime.After(etime) {
			err = dbgutil.FormatError("run %v timeout ", cmds)
			return
		}
		evt, err = syscall.WaitForSingleObject(hdl, uint32(100))
		if err != nil {
			err = dbgutil.FormatError("wait %v error [%s]", cmds, err.Error())
			return
		}
		switch evt {
		case syscall.WAIT_OBJECT_0:
			waiting = 0
			break
		case syscall.WAIT_TIMEOUT:
			break
		default:
			err = dbgutil.FormatError("wait %v evt %d", cmds, evt)
			return
		}
	}

	outstr = string(outb.Bytes())
	errstr = string(errb.Bytes())
	return
}

func Deamon() (err error) {
	err = dbgutil.FormatError("not supported Daemon")
	return
}

func GetTicksFromBoot() (retv uint64, err error) {
	var h syscall.Handle
	var proc uintptr
	var proc64 uintptr
	var retval uintptr
	retv = 0
	h, err = syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		err = dbgutil.FormatError("load kernel32.dll error %s", err.Error())
		return
	}
	defer func() {
		syscall.FreeLibrary(h)
	}()
	proc64, err = syscall.GetProcAddress(h, "GetTickCount64")
	if err != nil {
		proc, err = syscall.GetProcAddress(h, "GetTickCount")
		if err != nil {
			err = dbgutil.FormatError("GetTickCount search error %s", err.Error())
			return
		}
		retval, _, _ = syscall.Syscall(proc, 0, 0, 0, 0)
	} else {
		retval, _, _ = syscall.Syscall(proc64, 0, 0, 0, 0)
	}
	retv = uint64(retval)
	err = nil
	return
}

func GetChildProcs(pid int) (retp *ChildProcs, err error) {
	var outs string
	var exitcode int
	var cmds []string = []string{"wmic.exe", "process", "get", "ProcessId,ParentProcessId"}
	var sarr []string
	var carr []string
	var pidx int = -1
	var ppidx int = -1
	var l string
	var i int
	var spreg *regexp.Regexp
	var sps string
	var vmap map[string][]int
	var curpid, curppid int
	var nv string
	var ok bool
	var headone bool = false
	err = nil
	outs, _, exitcode, err = GetOutputCmd(cmds)
	if err != nil {
		return
	}
	if exitcode != 0 {
		err = dbgutil.FormatError("run %v exitcode %d", cmds, exitcode)
		return
	}
	sps = fmt.Sprintf("\\s+")
	spreg, err = regexp.Compile(sps)
	if err != nil {
		err = dbgutil.FormatError("compile [%s] error[%s]", sps, err.Error())
		return
	}

	vmap = make(map[string][]int)

	sarr = strings.Split(outs, "\n")
	for i = 0; i < len(sarr); i++ {
		l = sarr[i]
		l = strings.TrimRight(l, "\r")
		if !headone {
			carr = spreg.Split(l, -1)
			if len(carr) <= 1 {
				continue
			}
			headone = true
			//logutil.Debug("carr[0] [%s]", carr[0])
			if strings.ToLower(carr[0]) == "parentprocessid" {
				ppidx = 0
				pidx = 1
			} else if strings.ToLower(carr[0]) == "processid" {
				pidx = 0
				ppidx = 1
			} else {
				err = dbgutil.FormatError("[%s] not header", l)
				return
			}
		} else {
			carr = spreg.Split(l, -1)
			if len(carr) >= 2 {
				curppid, err = strconv.Atoi(carr[ppidx])
				if err == nil {
					curpid, err = strconv.Atoi(carr[pidx])
					if err == nil {
						nv = fmt.Sprintf("%d", curppid)
						_, ok = vmap[nv]
						if ok {
							vmap[nv] = append(vmap[nv], curpid)
						} else {
							vmap[nv] = []int{curpid}
						}
						//logutil.Debug("add [%s] = %v", nv, vmap[nv])
					}
				}
			}
		}
	}
	return find_childs(pid, vmap)
}
