package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func read_file_bytes(fname string) (rbytes []byte, err error) {
	var fp *os.File = os.Stdin
	err = nil
	rbytes = []byte{}
	defer func() {
		if fp != nil && fp != os.Stdin {
			fp.Close()
		}
		fp = nil
	}()

	if fname != "" {
		fp, err = os.Open(fname)
		if err != nil {
			err = fmt.Errorf("open [%s] error[%s]", fname, err.Error())
			return
		}
	}
	rbytes, err = ioutil.ReadAll(fp)
	return
}

func get_current_process_exec_info(pid int) (startaddr uintptr, endaddr uintptr, err error) {
	var exename string
	var exebyte []byte
	var exefile string
	var mapname string
	var s string
	var lines []string
	var curline string
	var i int
	startaddr = 0
	endaddr = 0
	err = fmt.Errorf("not valid")
	exename = fmt.Sprintf("/proc/%d/exe", pid)
	cmd := exec.Command("readlink", "-f", exename)
	exebyte, err = cmd.Output()
	if err != nil {
		return
	}
	exefile = string(exebyte)
	mapname = fmt.Sprintf("/proc/%d/maps", pid)
	mapbytes, err = read_file_bytes(mapname)
	if err != nil {
		return
	}
	s = string(mapbytes)
	lines = strings.Split(s, "\n")
	for i, curline = range lines {
		curline = strings.TrimRight(curline, "\r")
		if len(curline) == 0 {
			continue
		}

	}

	return
}
