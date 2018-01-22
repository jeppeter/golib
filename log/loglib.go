package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	l4g "github.com/jeppeter/log4go"
	"github.com/tebeka/atexit"
	"reflect"
	"runtime"
)

var st_logger *l4g.Logger = nil
var st_logger_level int = 0

func exithandler() {
	CloseDebugOutputBackGround()
	if st_logger != nil {
		st_logger.Close()
	}
	st_logger = nil
}

func init() {
	atexit.Register(exithandler)
}

func format_out_stack(level int) string {
	_, f, l, _ := runtime.Caller(level)
	return fmt.Sprintf("[%s:%d]", f, l)
}

func format_out_string_total(level int, fmtstr string, a ...interface{}) string {
	outstr := format_out_stack((level + 1))
	outstr += fmt.Sprintf(fmtstr, a...)
	outstr += "\n"
	return outstr
}

func format_out_string_singal(level int, fmtstr string) string {
	outstr := format_out_stack((level + 1))
	outstr += fmt.Sprintf(fmtstr)
	outstr += "\n"
	return outstr
}

func format_out_string_cap(a ...interface{}) string {
	var stacklevel int = 2
	var vaargs []interface{}
	var fmtstr string = ""
	var ct string
	if len(a) > 0 {
		switch v := a[0].(type) {
		case int:
			stacklevel = a[0].(int)
			stacklevel += 2
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
	outstr += "\n"
	return outstr
}

func Debug(a ...interface{}) int {
	var retval int = 0
	outstr := format_out_string_cap(a...)
	retval = len(outstr)
	if st_logger != nil {
		st_logger.Debug(outstr)
	}
	if st_logger_level >= 3 {
		LogDebugOutputBackGround(outstr)
	}
	return retval
}

func Error(a ...interface{}) int {
	var retval int = 0
	outstr := format_out_string_cap(a...)
	retval = len(outstr)
	if st_logger != nil {
		st_logger.Debug(outstr)
	}
	if st_logger_level >= 0 {
		LogDebugOutputBackGround(outstr)
	}
	return retval
}

func AddCliFlag(cliapp *cli.App) {
	var vflag *cli.IntFlag
	var wflag *cli.StringSliceFlag
	var aflag *cli.StringSliceFlag
	vflag = &cli.IntFlag{
		Name:  "verbose, V",
		Value: 0,
	}
	cliapp.Flags = append(cliapp.Flags, vflag)
	wflag = &cli.StringSliceFlag{
		Name:  "log-files",
		Usage: "set write rotate files",
	}
	cliapp.Flags = append(cliapp.Flags, wflag)
	aflag = &cli.StringSliceFlag{
		Name:  "log-appends",
		Usage: "set append files",
	}
	cliapp.Flags = append(cliapp.Flags, aflag)
	nostderr := &cli.BoolFlag{
		Name:  "log-nostderr",
		Usage: "specified no stdout",
	}
	cliapp.Flags = append(cliapp.Flags, nostderr)
	return
}

func SetCliFlag(ctx *cli.Context) error {
	var appfiles []string
	var cfiles []string
	var vmode int
	var lglvl l4g.Level
	var deflogfmt string = "[%T %D] (%S) %M"
	var clog l4g.Logger

	if st_logger != nil {
		st_logger.Close()
	}
	st_logger = nil

	vmode = ctx.GlobalInt("verbose")
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
	fmt.Printf("st_logger_level [%d]\n", st_logger_level)

	clog = l4g.NewLogger()
	st_logger = &clog
	if !ctx.GlobalIsSet("log-nostderr") {
		log4writer := l4g.NewStderrLogWriter()
		log4writer.SetFormat(deflogfmt)
		st_logger.AddFilter("stderr", lglvl, log4writer)
	}

	cfiles = ctx.GlobalStringSlice("log-files")
	if len(cfiles) > 0 {
		for _, f := range cfiles {
			log4writer := l4g.NewFileLogWriter(f, false)
			log4writer.SetFormat(deflogfmt)
			st_logger.AddFilter(f, lglvl, log4writer)
		}
	}

	appfiles = ctx.GlobalStringSlice("log-appends")
	if len(appfiles) > 0 {
		for _, f := range appfiles {
			log4writer := l4g.NewFileLogWriter(f, true)
			log4writer.SetFormat(deflogfmt)
			st_logger.AddFilter(f, lglvl, log4writer)
		}
	}
	return nil
}
