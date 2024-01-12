package executil

import (
	"bytes"
	"dbgutil"
	"logutil"
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
	if timeout > 0 {
		stime = time.Now()
		etime = stime.Add(time.Duration(timeout) * time.Second)
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
