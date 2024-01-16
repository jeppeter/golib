package executil

import (
	"bytes"
	"dbgutil"
	"logutil"
	"os"
	"os/exec"
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
