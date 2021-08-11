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

func run_cmd_output(exitch *chan int, cmds []string, stdinbytes []byte) (stdoutbytes []byte, stderrbytes []byte, exitcode int, err error) {
	var cmd *exec.Cmd = nil
	var stdinp io.WriteCloser = nil
	var nstdin *bufio.Writer
	var inlen int = 0
	var nret int
	var procwait *ProcWait = nil
	var outb bytes.Buffer
	var errb bytes.Buffer

	nstdin = nil

	cmd = &exec.Cmd{}
	cmd.Path = cmds[0]
	cmd.Args = cmds
	cmd.Stdout = &outb
	cmd.Stderr = &errb
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

	for {
		if exitch != nil {
			select {
			case <-*exitch:
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
			break
		}

		if nstdin != nil {
			if nstdin.Size() > 0 {
				err = nstdin.Flush()
				if err != nil {
					Error("%s", err.Error())
					return
				}
			} else {
				if inlen < len(stdinbytes) {
					nret, err = nstdin.Write(stdinbytes[inlen:])
					if err != nil {
						Error("%s", err.Error())
						return
					}
					inlen += nret
				} else {
					nstdin = nil
					stdinp.Close()
					stdinp = nil
				}
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
