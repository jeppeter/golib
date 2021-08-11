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

func run_cmd_output(exitch chan int, cmds []string, stdinbytes []byte) (stdoutbytes []byte, stderrbytes []byte, exitcode int, err error) {
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
			Error("%s", err.Error())
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
		Error("%s", err.Error())
		return
	}

	procwait, err = NewProcWait(cmd.Process)
	if err != nil {
		Error("%s", err.Error())
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

	Error("inlen %d stdinbytes [%d]", inlen, len(stdinbytes))

	for {
		if exitch != nil {
			select {
			case <-exitch:
				err = fmt.Errorf("exitch")
				Error("%s", err.Error())
				return
			case <-time.After(time.Duration(10) * time.Millisecond):
				inlen = inlen
			}
		} else {
			time.Sleep(time.Duration(10) * time.Millisecond)
		}
		if procwait.WaitExitTimeout(50) {
			Error("wait exit")
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
					Error("%s", err.Error())
					return
				}
				inlen += nret
				Error("nret %d inlen[%d] len[%d]", nret, inlen, len(stdinbytes))
			} else {
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

	Error("err out")

	err = nil
	return
}
