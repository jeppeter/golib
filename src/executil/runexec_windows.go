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
