package logutil

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	l4g "github.com/jeppeter/log4go"
	"github.com/tebeka/atexit"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	//"strconv"
	//"unicode/utf8"
)

type Background interface {
	LogDebugOutputBackGround(s string) error

	CloseDebugOutputBackGround() error
}

const (
	DEFAULT_LOG_MAX_LINES  = 100000000
	DEFAULT_LOG_ROTATE_NUM = 2
)

var st_logger *l4g.Logger = nil
var st_logger_level int = 0
var st_background Background = nil
var st_logger_nostderr_output bool = false

func exithandler() {
	st_background.CloseDebugOutputBackGround()
	if st_logger != nil {
		st_logger.Close()
	}
	st_logger = nil
	st_logger_level = 0
}

func init() {
	atexit.Register(exithandler)
	st_background = nativeGround()
}

func format_out_stack(level int) string {
	_, f, l, _ := runtime.Caller(level)
	return fmt.Sprintf("[%s:%d]", f, l)
}

func format_out_string_total(level int, fmtstr string, a ...interface{}) string {
	outstr := format_out_stack((level + 1))
	outstr += fmt.Sprintf(fmtstr, a...)
	return outstr
}

func format_out_string_singal(level int, fmtstr string) string {
	outstr := format_out_stack((level + 1))
	outstr += fmt.Sprintf(fmtstr)
	return outstr
}

func format_out_stack_data(data []byte, fmtstr string, a ...interface{}) string {
	var i, lasti int
	//var r rune
	//var p []byte
	var j int
	outstr := ""
	if fmtstr != "" {
		outstr += fmt.Sprintf(fmtstr, a...)
	}
	lasti = 0
	//p = make([]byte, 1)
	for i = 0; i < len(data); i++ {
		if (i % 16) == 0 {
			if i > 0 {
				outstr += "    "
				for i != lasti {
					//p[0] = data[lasti]
					//r, j = utf8.DecodeRune(p)
					if data[lasti] < ' ' || data[lasti] > '~' {
						outstr += "."
					} else {
						outstr += fmt.Sprintf("%c", data[lasti])
					}
					lasti++
				}
			}
			outstr += fmt.Sprintf("\n[0x%08x]", i)
		}
		outstr += fmt.Sprintf(" 0x%02x", data[i])
	}

	if lasti != i {
		j = i
		for (j % 16) != 0 {
			outstr += "     "
			j++
		}

		outstr += "    "
		for lasti != i {
			//p[0] = data[lasti]
			//r, j = utf8.DecodeRune(p)
			//if j != 1 || !strconv.IsPrint(r) {
			if data[lasti] < ' ' || data[lasti] > '~' {
				outstr += "."
			} else {
				outstr += fmt.Sprintf("%c", data[lasti])
			}
			lasti++
		}
		outstr += "\n"
	}

	return outstr
}

const (
	def_stacklevel_added = 3
)

func format_out_string_cap(a ...interface{}) string {
	var stacklevel int = def_stacklevel_added
	var vaargs []interface{}
	var fmtstr string = ""
	var ct string
	if len(a) > 0 {
		switch v := a[0].(type) {
		case int:
			stacklevel = a[0].(int)
			stacklevel += def_stacklevel_added
			if len(a) > 2 {
				vaargs = a[2:]
			}
			if len(a) > 1 {
				ct = reflect.TypeOf(a[1]).Name()
				if ct == "string" {
					fmtstr = a[1].(string)
				} else {
					fmtstr = "unknown type string"
				}
			}
		case string:
			if len(a) > 1 {
				vaargs = a[1:]
			}
			fmtstr = a[0].(string)
		default:
			fmtstr = fmt.Sprintf("unknown type [%s]", v)
		}
	}

	outstr := format_out_stack(stacklevel)
	if len(vaargs) == 0 {
		outstr += fmt.Sprintf(fmtstr)
	} else {
		outstr += fmt.Sprintf(fmtstr, vaargs...)
	}
	return outstr
}

func format_out_data_cap(a ...interface{}) string {
	var stacklevel int = def_stacklevel_added
	var vaargs []interface{}
	var data []byte = []byte{}
	var fmtstr string = ""
	var ct string
	if len(a) > 0 {
		switch v := a[0].(type) {
		case int:
			stacklevel = a[0].(int)
			stacklevel += def_stacklevel_added
			if len(a) > 3 {
				vaargs = a[3:]
			}
			if len(a) > 1 {
				switch a[1].(type) {
				case []byte:
					data = a[1].([]byte)
				}
			}

			if len(a) > 2 {
				ct = reflect.TypeOf(a[2]).Name()
				if ct == "string" {
					fmtstr = a[2].(string)
				} else {
					fmtstr = fmt.Sprintf("unknown type string [%s]", ct)
				}
			}
		case []byte:
			if len(a) > 2 {
				vaargs = a[2:]
			}
			data = a[0].([]byte)
			if len(a) > 1 {
				ct = reflect.TypeOf(a[1]).Name()
				if ct == "string" {
					fmtstr = a[1].(string)
				} else {
					fmtstr = fmt.Sprintf("unknown type [%s]", ct)
				}
			}

		default:
			fmtstr = fmt.Sprintf("unknown type [%s]", v)
			fmt.Printf("%s\n", fmtstr)
		}
	}

	outstr := format_out_stack(stacklevel)
	if len(vaargs) == 0 {
		outstr += format_out_stack_data(data, fmtstr)
	} else {
		outstr += format_out_stack_data(data, fmtstr, vaargs...)
	}
	return outstr
}

func output_stderr(level int, outs string) {
	if level <= st_logger_level && !st_logger_nostderr_output {
		fmt.Fprintf(os.Stderr, "%s\n", outs)
		os.Stderr.Sync()
	}
}

func Error(a ...interface{}) int {
	var retval int = 0
	outstr := "<ERROR>"
	outstr += format_out_string_cap(a...)
	retval = len(outstr)
	if st_logger != nil {
		st_logger.Error(outstr)
	} else {
		fmt.Fprintf(os.Stderr, "no out %s", outstr)
	}
	output_stderr(0, outstr)
	if st_logger_level >= 0 {
		outstr += "\n"
		st_background.LogDebugOutputBackGround(outstr)
	}
	return retval
}

func ErrorBuffer(a ...interface{}) int {
	var retval int = 0
	outstr := "<ERROR>"
	outstr += format_out_data_cap(a...)
	retval = len(outstr)
	if st_logger != nil {
		st_logger.Error(outstr)
	} else {
		fmt.Fprintf(os.Stderr, "no out %s", outstr)
	}
	output_stderr(0, outstr)
	if st_logger_level >= 0 {
		outstr += "\n"
		st_background.LogDebugOutputBackGround(outstr)
	}
	return retval
}

func Warn(a ...interface{}) int {
	var retval int = 0
	outstr := "<WARN>"
	outstr += format_out_string_cap(a...)
	retval = len(outstr)
	if st_logger != nil {
		st_logger.Warn(outstr)
	} else {
		fmt.Fprintf(os.Stderr, "no out %s", outstr)
	}
	output_stderr(1, outstr)
	if st_logger_level >= 1 {
		outstr += "\n"
		st_background.LogDebugOutputBackGround(outstr)
	}
	return retval
}

func WarnBuffer(a ...interface{}) int {
	var retval int = 0
	outstr := "<WARN>"
	outstr += format_out_data_cap(a...)
	retval = len(outstr)
	if st_logger != nil {
		st_logger.Warn(outstr)
	} else {
		fmt.Fprintf(os.Stderr, "no out %s", outstr)
	}
	output_stderr(1, outstr)
	if st_logger_level >= 1 {
		outstr += "\n"
		st_background.LogDebugOutputBackGround(outstr)
	}
	return retval
}

func Info(a ...interface{}) int {
	var retval int = 0
	outstr := "<INFO>"
	outstr += format_out_string_cap(a...)
	retval = len(outstr)
	if st_logger != nil {
		st_logger.Info(outstr)
	} else {
		fmt.Fprintf(os.Stderr, "no out %s", outstr)
	}
	output_stderr(2, outstr)
	if st_logger_level >= 2 {
		outstr += "\n"
		st_background.LogDebugOutputBackGround(outstr)
	}
	return retval
}

func InfoBuffer(a ...interface{}) int {
	var retval int = 0
	outstr := "<INFO>"
	outstr += format_out_data_cap(a...)
	retval = len(outstr)
	if st_logger != nil {
		st_logger.Info(outstr)
	} else {
		fmt.Fprintf(os.Stderr, "no out %s", outstr)
	}
	output_stderr(2, outstr)
	if st_logger_level >= 2 {
		outstr += "\n"
		st_background.LogDebugOutputBackGround(outstr)
	}
	return retval
}

func Debug(a ...interface{}) int {
	var retval int = 0
	outstr := "<DEBUG>"
	outstr += format_out_string_cap(a...)
	retval = len(outstr)
	if st_logger != nil {
		st_logger.Debug(outstr)
	}
	output_stderr(3, outstr)
	if st_logger_level >= 3 {
		outstr += "\n"
		st_background.LogDebugOutputBackGround(outstr)
	}
	return retval
}

func DebugBuffer(a ...interface{}) int {
	var retval int = 0
	outstr := "<DEBUG>"
	outstr += format_out_data_cap(a...)
	retval = len(outstr)
	if st_logger != nil {
		st_logger.Debug(outstr)
	}
	output_stderr(3, outstr)
	if st_logger_level >= 3 {
		outstr += "\n"
		st_background.LogDebugOutputBackGround(outstr)
	}
	return retval
}

func Trace(a ...interface{}) int {
	var retval int = 0
	outstr := "<TRACE>"
	outstr += format_out_string_cap(a...)
	retval = len(outstr)
	if st_logger != nil {
		st_logger.Trace(outstr)
	}
	output_stderr(4, outstr)
	if st_logger_level >= 4 {
		outstr += "\n"
		st_background.LogDebugOutputBackGround(outstr)
	}
	return retval
}

func TraceBuffer(a ...interface{}) int {
	var retval int = 0
	outstr := "<TRACE>"
	outstr += format_out_data_cap(a...)
	retval = len(outstr)
	if st_logger != nil {
		st_logger.Trace(outstr)
	}
	output_stderr(4, outstr)
	if st_logger_level >= 4 {
		outstr += "\n"
		st_background.LogDebugOutputBackGround(outstr)
	}
	return retval
}

func PrepareLog(parser *extargsparse.ExtArgsParse) error {
	var commandline = `{
			"verbose|v" : "+",
			"log-files##set write rotate files format name,maxlines##" : [],
			"log-appends##set append files name,maxlines,maxbackups##" : [],
			"log-nostderr##specified no stderr output##" : false
		}`
	var err error
	err = parser.LoadCommandLineString(commandline)
	return err
}

func filter_maxfiles(fname string, maxbackup int) (err error) {
	var dname string
	var bname string
	var finfos []os.FileInfo
	var fh *os.File
	var finfo os.FileInfo
	var i int
	var numexpr *regexp.Regexp
	var exprstr string
	var matchstrings []string
	var ival int
	var curname string
	err = nil
	dname = filepath.Dir(fname)
	bname = filepath.Base(fname)
	fh, err = os.Open(dname)
	if err != nil {
		fmt.Printf("open [%s]error [%s]\n", fname, err.Error())
		err = nil
		return
	}
	defer fh.Close()

	exprstr = fmt.Sprintf("%s\\.([0-9]+)", bname)
	numexpr, err = regexp.Compile(exprstr)
	if err != nil {
		err = fmt.Errorf("%s", format_out_string_total(1, "can not compile [%s] error[%s]", exprstr, err.Error()))
		return
	}

	for {
		finfos, err = fh.Readdir(10)
		if err != nil {
			err = nil
			return
		}
		for i, finfo = range finfos {
			if !finfo.IsDir() {
				matchstrings = numexpr.FindStringSubmatch(finfo.Name())
				if len(matchstrings) > 1 {
					ival, err = strconv.Atoi(matchstrings[1])
					if err == nil {
						if ival >= maxbackup {
							curname = filepath.Join(dname, finfo.Name())
							err = os.Remove(curname)
							if err != nil {
								err = fmt.Errorf("%s", format_out_string_total(1, "[%d]can not remove [%s] error[%s]", i, curname, err.Error()))
								return
							}
						}
					}
				}
			}

		}
	}

	err = nil
	return

}

func InitLog(ns *extargsparse.NameSpaceEx) error {
	var appfiles []string
	var cfiles []string
	var vmode int
	var lglvl l4g.Level
	var deflogfmt string = "[%T %D] %M"
	var clog l4g.Logger
	var sarr []string
	var err error
	var maxlines int
	var maxbackups int
	var fname string

	if st_logger != nil {
		st_logger.Close()
	}
	st_logger = nil

	vmode = ns.GetInt("verbose")
	if vmode <= 0 {
		lglvl = l4g.ERROR
		st_logger_level = 0
	} else if vmode == 1 {
		lglvl = l4g.WARNING
		st_logger_level = 1
	} else if vmode == 2 {
		lglvl = l4g.INFO
		st_logger_level = 2
	} else if vmode == 3 {
		lglvl = l4g.DEBUG
		st_logger_level = 3
	} else if vmode == 4 {
		lglvl = l4g.TRACE
		st_logger_level = 4
	} else if vmode >= 5 {
		lglvl = l4g.FINEST
		st_logger_level = 5
	}

	clog = l4g.NewLogger()
	st_logger = &clog
	st_logger_nostderr_output = ns.GetBool("log_nostderr")

	cfiles = ns.GetArray("log_files")
	if len(cfiles) > 0 {
		for _, f := range cfiles {
			sarr = strings.Split(f, ",")
			fname = sarr[0]
			maxlines = DEFAULT_LOG_MAX_LINES
			if len(sarr) > 1 {
				maxlines, err = strconv.Atoi(sarr[1])
				if err != nil {
					maxlines = DEFAULT_LOG_MAX_LINES
				}
			}

			maxbackups = DEFAULT_LOG_ROTATE_NUM
			if len(sarr) > 2 {
				maxbackups, err = strconv.Atoi(sarr[2])
				if err != nil {
					maxbackups = DEFAULT_LOG_ROTATE_NUM
				}
			}
			if maxbackups < DEFAULT_LOG_ROTATE_NUM {
				maxbackups = DEFAULT_LOG_ROTATE_NUM
			}

			filter_maxfiles(fname, maxbackups)

			log4writer := l4g.NewFileLogWriter(fname, true)
			log4writer.SetFormat(deflogfmt)
			log4writer.SetRotateLines(maxlines)

			log4writer.SetRotateMaxBackup(maxbackups)
			st_logger.AddFilter(f, lglvl, log4writer)
			clog[f].Level = lglvl
		}
	}

	appfiles = ns.GetArray("log_appends")
	if len(appfiles) > 0 {
		for _, f := range appfiles {
			sarr = strings.Split(f, ",")
			fname = sarr[0]
			maxlines = DEFAULT_LOG_MAX_LINES
			if len(sarr) > 1 {
				maxlines, err = strconv.Atoi(sarr[1])
				if err != nil {
					maxlines = DEFAULT_LOG_MAX_LINES
				}
			}

			maxbackups = DEFAULT_LOG_ROTATE_NUM
			if len(sarr) > 2 {
				maxbackups, err = strconv.Atoi(sarr[2])
				if err != nil {
					maxbackups = DEFAULT_LOG_ROTATE_NUM
				}
			}

			if maxbackups < DEFAULT_LOG_ROTATE_NUM {
				maxbackups = DEFAULT_LOG_ROTATE_NUM
			}

			filter_maxfiles(fname, maxbackups)

			log4writer := l4g.NewFileLogWriter(fname, true)
			log4writer.SetFormat(deflogfmt)
			log4writer.SetRotateLines(maxlines)

			log4writer.SetRotateMaxBackup(maxbackups)
			st_logger.AddFilter(f, lglvl, log4writer)
			clog[f].Level = lglvl
		}
	}
	return nil
}
