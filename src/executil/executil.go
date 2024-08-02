package executil

import (
	"bytes"
	"dbgutil"
	"fmt"
	"logutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetExeFile() (exename string, err error) {
	var paths []string
	var envpath string
	var curpath string
	var wholepath string
	var finfo os.FileInfo
	err = nil
	exename = ""
	exename, err = filepath.Abs(os.Args[0])
	if err == nil {
		finfo, err = os.Stat(exename)
		if err == nil && !finfo.IsDir() && finfo.Size() > 0 {
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
			exename = wholepath
			err = nil
			return
		}
	}
	exename = ""
	err = dbgutil.FormatError("can not get exe from path")
	return
}

func GetExeDir() (dirname string, err error) {
	var exename string
	exename, err = GetExeFile()
	if err != nil {
		return
	}
	dirname = filepath.Dir(exename)
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
			err = dbgutil.FormatError("run %v error[%s]", cmds, err.Error())
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

type ChildProcs struct {
	Pid      int
	Children []*ChildProcs
}

func NewChildProcs(pid int) *ChildProcs {
	var retp *ChildProcs
	retp = &ChildProcs{}
	retp.Pid = pid
	retp.Children = []*ChildProcs{}
	return retp
}

func (p *ChildProcs) tab_string(tab int, s string) string {
	var rets string
	var i int
	for i = 0; i < tab; i++ {
		rets += fmt.Sprintf("    ")
	}
	rets += s
	rets += "\n"
	return rets
}

func (p *ChildProcs) format(tab int) string {
	var rets string
	var i int
	rets += p.tab_string(tab, fmt.Sprintf("%d", p.Pid))
	for i = 0; i < len(p.Children); i++ {
		rets += p.Children[i].format(tab + 1)
	}
	return rets
}

func (p *ChildProcs) String() string {
	return p.format(0)
}

func has_searched(pid int, pids []int) bool {
	var i int
	for i = 0; i < len(pids); i++ {
		if pid == pids[i] {
			return true
		}
	}
	return false
}

func find_childs(pid int, vmap map[string][]int) (retp *ChildProcs, err error) {
	var curfindpids []*ChildProcs
	var nextfindpids []*ChildProcs
	var docont bool
	var cpids []int
	var ok bool
	var insertpids *ChildProcs
	var cch *ChildProcs
	var i, j int
	var bval bool
	var searchpids []int = []int{}
	var nv string
	retp = NewChildProcs(pid)
	/*now to make curpids*/
	curfindpids = []*ChildProcs{retp}
	docont = true
	for docont {
		nextfindpids = []*ChildProcs{}
		docont = false
		for i = 0; i < len(curfindpids); i++ {
			bval = has_searched(curfindpids[i].Pid, searchpids)
			if !bval {
				nv = fmt.Sprintf("%d", curfindpids[i].Pid)
				cpids, ok = vmap[nv]
				if ok {
					docont = true
					insertpids = curfindpids[i]
					logutil.Debug("[%s] next add search %v", nv, cpids)

					for j = 0; j < len(cpids); j++ {
						if cpids[j] != curfindpids[i].Pid {
							cch = NewChildProcs(cpids[j])
							insertpids.Children = append(insertpids.Children, cch)
							nextfindpids = append(nextfindpids, cch)
						}
					}
				}
				searchpids = append(searchpids, curfindpids[i].Pid)
			}

		}
		curfindpids = nextfindpids
	}
	err = nil
	return
}
