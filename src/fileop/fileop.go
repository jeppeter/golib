package fileop

import (
	"dbgutil"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
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

func GetExeFull() (fullname string, err error) {
	var paths []string
	var envpath string
	var curpath string
	var wholepath string
	var finfo os.FileInfo
	err = nil
	fullname = ""
	fullname, err = filepath.Abs(os.Args[0])
	if err == nil {
		finfo, err = os.Stat(fullname)
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
			fullname = wholepath
			err = nil
			return
		}
	}
	fullname = ""
	err = dbgutil.FormatError("can not get exe from path")
	return
}
