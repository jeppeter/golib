package executil

import (
	"bytes"
	"dbgutil"
	"fileop"
	"fmt"
	"logutil"
	"os"
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
	var cmd *exec.Cmd
	var waiting int = 1
	var wpid int
	for i, c := range cmds {
		logutil.Trace("[%d]=[%s]", i, c)
	}
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
	if timeout > 0 {
		stime = time.Now()
		etime = stime.Add(time.Duration(timeout) * time.Millisecond)
	}
	waiting = 1
	for waiting > 0 {
		ctime = time.Now()
		if timeout > 0 && ctime.After(etime) {
			err = dbgutil.FormatError("run %v timeout ", cmds)
			return
		}

		wpid, _ = syscall.Wait4(cmd.Process.Pid, nil, syscall.WNOHANG, nil)
		if wpid == cmd.Process.Pid {
			break
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
	}

	outstr = string(outb.Bytes())
	errstr = string(errb.Bytes())
	return
}

func Deamon() (err error) {
	var id uintptr
	id, _, _ = syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if err != nil {
		err = dbgutil.FormatError("[%d]fork error [%s]", os.Getpid(), err.Error())
		return
	}
	if id != 0 {
		os.Exit(0)
	}

	_, err = syscall.Setsid()
	if err != nil {
		err = dbgutil.FormatError("[%d]Setsid error [%s]", os.Getpid(), err.Error())
		return
	}
	err = nil
	return

}

func GetTicksFromBoot() (retv uint64, err error) {
	var s string
	s, err = fileop.ReadFile("/proc/uptime")
	if err != nil {
		return
	}
	var sarr []string
	var restr string = "\\s+"
	var reg *regexp.Regexp
	var totalf float64
	reg, err = regexp.Compile(restr)
	if err != nil {
		err = dbgutil.FormatError("compile [%s] error [%s]", restr, err.Error())
		return
	}
	sarr = reg.Split(s, -1)
	if len(sarr) <= 1 {
		err = dbgutil.FormatError("split [%s] not valid", s)
		return
	}

	totalf, err = strconv.ParseFloat(sarr[1], 64)
	if err != nil {
		err = dbgutil.FormatError("parse [%s] not valid [%s]", sarr[1], err.Error())
		return
	}

	retv = uint64(totalf * 1000)
	err = nil
	return

}

func GetChildProcs(pid int) (retp *ChildProcs, err error) {
	var outs string
	var exitcode int
	var cmds []string = []string{"/bin/ps", "-ef"}
	var sarr []string
	var l string
	var i int
	var pidx, ppidx int
	var carr []string
	var spreg *regexp.Regexp
	var sps string
	var vmap map[string][]int
	var curpid, curppid int
	var nv string
	var ok bool

	outs, _, exitcode, err = GetOutputCmd(cmds)
	if err != nil {
		return
	}

	if exitcode != 0 {
		err = dbgutil.FormatError("%v exitcode %d", cmds, exitcode)
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
		logutil.Debug("[%d]=[%s]", i, l)
		if i == 0 {
			pidx = 1
			ppidx = 2
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
						logutil.Debug("add [%s] = %v", nv, vmap[nv])
					}
				}
			}
		}
	}

	return find_childs(pid, vmap)
}
