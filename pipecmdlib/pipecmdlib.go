package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"
)

const (
	CMD_PEEK_SIZE = 4096
)

type ProcWait struct {
	pid      int
	handle   uintptr // handle is accessed atomically on Windows
	isdone   int
	exitcode int
}

func run_cmd_output_evt(exitch chan int, cmds []string, stdinbytes []byte) (stdoutbytes []byte, stderrbytes []byte, exitcode int, err error) {
	var cmd *exec.Cmd = nil
	var stdinp io.WriteCloser = nil
	var nstdin *bufio.Writer
	var inlen int = 0
	var nret int
	var procwait *ProcWait = nil
	var outb bytes.Buffer
	var errb bytes.Buffer
	var curlen int

	nstdin = nil

	cmd = &exec.Cmd{}
	cmd.Path = cmds[0]
	cmd.Args = cmds
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	stdoutbytes = []byte{}
	stderrbytes = []byte{}
	defer cmd.Wait()

	if len(stdinbytes) > 0 {
		stdinp, err = cmd.StdinPipe()
		if err != nil {
			return
		}
		defer func() {
			if stdinp != nil {
				stdinp.Close()
			}
			stdinp = nil
		}()
		nstdin = bufio.NewWriter(stdinp)
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	procwait, err = NewProcWait(cmd.Process)
	if err != nil {
		return
	}
	defer procwait.Close()

	defer func() {
		if cmd != nil {
			if cmd.Process != nil {
				cmd.Process.Kill()
				cmd.Process = nil
			}
		}
	}()

	for {
		if exitch != nil {
			select {
			case <-exitch:
				err = fmt.Errorf("exit notified")
				return
			case <-time.After(time.Duration(10) * time.Millisecond):
				inlen = inlen
			}
		} else {
			time.Sleep(time.Duration(10) * time.Millisecond)
		}
		if procwait.WaitExitTimeout(50) {
			procwait.Close()
			break
		}

		if nstdin != nil {
			if inlen < len(stdinbytes) {
				curlen = nstdin.Size()
				if (len(stdinbytes) - inlen) < curlen {
					curlen = len(stdinbytes) - inlen
				}
				nret, err = nstdin.Write(stdinbytes[inlen:(inlen + curlen)])
				if err != nil {
					return
				}
				inlen += nret
			} else {
				err = nstdin.Flush()
				if err != nil {
					return
				}
				nstdin = nil
				stdinp.Close()
				stdinp = nil
			}
		}

	}

	cmd.Wait()

	exitcode = procwait.GetExitcode()
	cmd.Process = nil
	cmd.ProcessState = nil
	stdoutbytes = outb.Bytes()
	stderrbytes = errb.Bytes()

	err = nil
	return
}

type ProcCommandHandle struct {
	exitch      chan int
	exited      int
	exitcode    int
	cmds        []string
	errstr      string
	inputbytes  []byte
	outputbytes []byte
	errbytes    []byte
}

func NewProcCommandHandle(cmds []string, inputbytes []byte) (retp *ProcCommandHandle, err error) {
	if len(cmds) == 0 {
		retp = nil
		err = fmt.Errorf("len(cmds) == 0")
		return
	}
	retp = &ProcCommandHandle{}
	retp.exitch = make(chan int, 10)
	retp.exited = 1
	retp.exitcode = 1
	retp.errstr = ""
	retp.inputbytes = inputbytes
	retp.outputbytes = []byte{}
	retp.errbytes = []byte{}
	retp.cmds = cmds
	err = nil
	return
}

func run_proc_handle(retp *ProcCommandHandle) {
	var err error = nil
	var exitcode int = 1
	var outbytes, errbytes []byte

	defer func() {
		if err != nil {
			retp.errstr = fmt.Sprintf("%s", err.Error())
		} else {
			retp.errstr = ""
		}
		retp.exitcode = exitcode
		retp.exited = 1
	}()

	outbytes, errbytes, exitcode, err = run_cmd_output_evt(retp.exitch, retp.cmds, retp.inputbytes)
	if err != nil {
		return
	}
	retp.outputbytes = outbytes
	retp.errbytes = errbytes
	return
}

func (retp *ProcCommandHandle) Start() (err error) {
	if retp.exited == 0 {
		err = fmt.Errorf("%v already started", retp.cmds)
		return
	}

	retp.exited = 0
	retp.exitcode = 1
	retp.errstr = ""
	go run_proc_handle(retp)
	err = nil
	return
}

func (retp *ProcCommandHandle) Exited() bool {
	if retp.exited > 0 {
		return true
	}
	return false
}

func (retp *ProcCommandHandle) Stop() bool {
	if retp.exited > 0 {
		return true
	}
	retp.exitch <- 1
	for {
		if retp.Exited() {
			return true
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
	return true
}

func (retp *ProcCommandHandle) GetOutput() (outbytes []byte, errbytes []byte, err error) {
	if retp.exited == 0 {
		err = fmt.Errorf("not exited")
		return
	}
	outbytes = retp.outputbytes
	errbytes = retp.errbytes
	err = nil
	return
}

func (retp *ProcCommandHandle) GetExitcode() (exitcode int, err error) {
	if retp.exited == 0 {
		err = fmt.Errorf("not exited")
		return
	}
	err = nil
	exitcode = retp.exitcode
	return
}

func (retp *ProcCommandHandle) Close() {
	retp.Stop()
	retp.outputbytes = []byte{}
	retp.errbytes = []byte{}
	retp.inputbytes = []byte{}
	retp.cmds = []string{}
	return
}
