package fileop

import (
	"dbgutil"
	"errors"
	"io/ioutil"
	"os"
)

func ReadFileBytes(fname string) (rbytes []byte, err error) {
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
			err = dbgutil.FormatError("open [%s] error[%s]", fname, err.Error())
			return
		}
	}
	rbytes, err = ioutil.ReadAll(fp)
	return
}

func WriteFileBytes(fname string, obytes []byte) (nret int, err error) {
	var fp *os.File = os.Stdout
	err = nil
	nret = 0

	defer func() {
		if fp != nil && fp != os.Stdout {
			fp.Close()
		}
		fp = nil
	}()

	if fname != "" {
		fp, err = os.Create(fname)
		if err != nil {
			err = dbgutil.FormatError("create [%s] error[%s]", fname, err.Error())
			return
		}
	}
	nret, err = fp.Write(obytes)

	return
}

func ReadFile(fname string) (s string, err error) {
	var ob []byte
	s = ""
	ob, err = ReadFileBytes(fname)
	if err != nil {
		return
	}
	s = string(ob)
	return
}

func WriteFile(fname string, ostring string) (nret int, err error) {
	nret, err = WriteFileBytes(fname, []byte(ostring))
	return
}

func DeleteFile(fname string) (err error) {
	var err2 error
	err = os.Remove(fname)
	if err != nil {
		_, err2 = os.Stat(fname)
		if err2 != nil {
			if errors.Is(err2, os.ErrNotExist) {
				err = nil
			}
		}
	}
	return
}
