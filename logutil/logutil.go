package logutil

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	l4g "github.com/jeppeter/log4go"
	"github.com/tebeka/atexit"
	"os"
	"reflect"
	"runtime"
	//"strconv"
	//"unicode/utf8"
)

type Background interface {
	LogDebugOutputBackGround(s string) error

	CloseDebugOutputBackGround() error
}

var st_logger *l4g.Logger = nil
var st_logger_level int = 0
var st_background Background = nil

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
	if st_logger_level >= 4 {
		outstr += "\n"
		st_background.LogDebugOutputBackGround(outstr)
	}
	return retval
}

func PrepareLog(parser *extargsparse.ExtArgsParse) error {
	var commandline = `{
			"verbose|v" : "+",
			"log-files##set write rotate files##" : [],
			"log-appends##set append files##" : [],
			"log-nostderr##specified no stderr output##" : false
		}`
	var err error
	err = parser.LoadCommandLineString(commandline)
	return err
}

func InitLog(ns *extargsparse.NameSpaceEx) error {
	var appfiles []string
	var cfiles []string
	var vmode int
	var lglvl l4g.Level
	var deflogfmt string = "[%T %D] %M"
	var clog l4g.Logger

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
	if !ns.GetBool("log_nostderr") {
		log4writer := l4g.NewStderrLogWriter()
		log4writer.SetFormat(deflogfmt)
		st_logger.AddFilter("stderr", lglvl, log4writer)
		clog["stderr"].Level = lglvl
	}

	cfiles = ns.GetArray("log_files")
	if len(cfiles) > 0 {
		for _, f := range cfiles {
			log4writer := l4g.NewFileLogWriter(f, true)
			log4writer.SetFormat(deflogfmt)
			st_logger.AddFilter(f, lglvl, log4writer)
			clog[f].Level = lglvl
		}
	}

	appfiles = ns.GetArray("log_appends")
	if len(appfiles) > 0 {
		for _, f := range appfiles {
			log4writer := l4g.NewFileLogWriter(f, false)
			log4writer.SetFormat(deflogfmt)
			st_logger.AddFilter(f, lglvl, log4writer)
			clog[f].Level = lglvl
		}
	}
	return nil
}
