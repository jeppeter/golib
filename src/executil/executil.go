package executil

import (
	"bytes"
	"dbgutil"
	"logutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetExeDir() (dirname string, err error) {
	var paths []string
	var envpath string
	var curpath string
	var wholepath string
	var finfo os.FileInfo
	err = nil
	dirname = ""
	dirname, err = filepath.Abs(os.Args[0])
	if err == nil {
		finfo, err = os.Stat(dirname)
		if err == nil && !finfo.IsDir() && finfo.Size() > 0 {
			dirname = filepath.Dir(dirname)
			return
		}
	}

	/*it may be called by the path*/
	envpath = os.Getenv("PATH")
	if len(envpath) == 0 {
		err = dbgutil.FormatError("can not get PATH")
		return
	}

	if runtime.GOOS == "windows" {
		paths = strings.Split(envpath, ";")
	} else {
		paths = strings.Split(envpath, ":")
	}
	for _, curpath = range paths {
		wholepath = path.Join(curpath, os.Args[0])
		finfo, err = os.Stat(wholepath)
		if err != nil {
			continue
		}
		if !finfo.IsDir() && finfo.Size() > 0 {
			dirname = filepath.Dir(wholepath)
			err = nil
			return
		}
	}
	dirname = ""
	err = dbgutil.FormatError("can not get exe from path")
	return
}

func GetOutputCmdBytes(cmds []string) (outbs []byte, errbs []byte, exitcode int, err error) {
	var cmd *exec.Cmd
	var outb bytes.Buffer
	var errb bytes.Buffer
	var exiterr *exec.ExitError
	var ok bool
	logutil.Trace("cmds %v", cmds)
	cmd = &exec.Cmd{}
	cmd.Path = cmds[0]
	cmd.Args = cmds
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	outbs = []byte{}
	errbs = []byte{}
	exitcode = 0

	err = cmd.Run()
	if err != nil {
		exiterr, ok = err.(*exec.ExitError)
		if !ok {
			err = dbgutil.FormatError("run [%v] error[%s]", cmds, err.Error())
			return
		}
		exitcode = exiterr.ExitCode()
		err = nil
	}

	logutil.TraceBuffer(outb.Bytes(), "run %v outb", cmds)
	logutil.TraceBuffer(errb.Bytes(), "run %v errb", cmds)

	errbs = errb.Bytes()
	outbs = outb.Bytes()
	return
}

func GetOutputCmd(cmds []string) (outstr string, errstr string, exitcode int, err error) {
	var outbs, errbs []byte
	errstr = ""
	outstr = ""
	outbs, errbs, exitcode, err = GetOutputCmdBytes(cmds)
	if err == nil {
		errstr = string(errbs)
		outstr = string(outbs)
	}
	return
}

func StartAndDetach(cmds []string) (pid int, err error) {
	var cmd *exec.Cmd
	var exiterr *exec.ExitError
	var ok bool
	var exitcode int
	logutil.Trace("cmds [%v]", cmds)
	cmd = &exec.Cmd{}
	cmd.Path = cmds[0]
	cmd.Args = cmds

	err = cmd.Start()
	if err != nil {
		exiterr, ok = err.(*exec.ExitError)
		if !ok {
			err = dbgutil.FormatError("run [%v] error[%s]", cmds, err.Error())
			return
		}
		exitcode = exiterr.ExitCode()
		err = dbgutil.FormatError("run %v exit [%d]", cmds, exitcode)
		return
	}

	pid = cmd.Process.Pid
	err = nil
	cmd.Process.Release()
	return
}
