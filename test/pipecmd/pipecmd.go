package main

import (
	"bufio"
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
	var stdoutp io.ReadCloser = nil
	var stderrp io.ReadCloser = nil
	var stdinp io.WriteCloser = nil
	var nstdout, nstderr *bufio.Reader
	var nstdin *bufio.Writer
	var inlen int = 0
	var outlen int = 0
	var errlen int = 0
	var nret int
	var nbytes []byte
	var c byte
	var procwait *ProcWait = nil

	nstdout = nil
	nstderr = nil
	nstdin = nil

	cmd = &exec.Cmd{}
	cmd.Path = cmds[0]
	cmd.Args = cmds

	stdoutp, err = cmd.StdoutPipe()
	if err != nil {
		Error("%s", err.Error())
		return
	}
	defer stdoutp.Close()
	nstdout = bufio.NewReader(stdoutp)

	stderrp, err = cmd.StderrPipe()
	if err != nil {
		Error("%s", err.Error())
		return
	}
	defer stderrp.Close()
	nstderr = bufio.NewReader(stderrp)

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
				outlen = outlen
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

		if nstdout != nil {
			for {
				nbytes, err = nstdout.Peek(CMD_PEEK_SIZE)
				if err != nil {
					if err == io.EOF {
						break
					}
					Error("%s", err.Error())
					return
				}
				Error("stdout [%d]", len(nbytes))
				if len(nbytes) == 0 {
					break
				} else if len(nbytes) > 0 {
					for _, c = range nbytes {
						stdoutbytes = append(stdoutbytes, c)
					}
					nret, err = nstdout.Read(nbytes)
					if err != nil {
						Error("%s", err.Error())
						return
					}
					outlen += nret
				}
			}
		}

		if nstderr != nil {
			for {
				nbytes, err = nstderr.Peek(CMD_PEEK_SIZE)
				if err != nil {
					if err == io.EOF {
						break
					}
					Error("%s", err.Error())
					return
				}
				Error("stderr [%d]", len(nbytes))
				if len(nbytes) == 0 {
					break
				} else if len(nbytes) > 0 {
					for _, c = range nbytes {
						stderrbytes = append(stderrbytes, c)
					}
					nret, err = nstderr.Read(nbytes)
					if err != nil {
						Error("%s", err.Error())
						return
					}
					errlen += nret
				}
			}
		}
	}

	if nstdout != nil {
		for {
			nbytes, err = nstdout.Peek(CMD_PEEK_SIZE)
			if err != nil {
				if err == io.EOF {
					break
				}
				Error("%s", err.Error())
				return
			}
			if len(nbytes) == 0 {
				nstdout = nil
				stdoutp.Close()
				stdoutp = nil
				break
			} else if len(nbytes) > 0 {
				for _, c = range nbytes {
					stdoutbytes = append(stdoutbytes, c)
				}
				nret, err = nstderr.Read(nbytes)
				if err != nil {
					Error("%s", err.Error())
					return
				}
				outlen += nret
			}
		}
	}

	if nstderr != nil {
		for {
			nbytes, err = nstderr.Peek(CMD_PEEK_SIZE)
			if err != nil {
				if err == io.EOF {
					break
				}
				Error("%s", err.Error())
				return
			}
			if len(nbytes) == 0 {
				nstderr = nil
				stderrp.Close()
				stderrp = nil
				break
			} else if len(nbytes) > 0 {
				for _, c = range nbytes {
					stderrbytes = append(stderrbytes, c)
				}
				nret, err = nstderr.Read(nbytes)
				if err != nil {
					Error("%s", err.Error())
					return
				}
				errlen += nret
			}
		}
	}

	exitcode = procwait.GetExitcode()
	cmd.Process = nil
	cmd.ProcessState = nil

	err = nil
	return
}
