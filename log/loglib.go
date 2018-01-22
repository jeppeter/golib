package main

import (
	"github.com/codegangsta/cli"
	l4g "github.com/jeppeter/log4go"
	"fmt"
	"github.com/tebaka/atexit"
	"flag"
)

var st_logger *l4g.Logger = nil

func exithandler() {
	if st_logger != nil {
		st_logger.Close()
	}
	st_logger = nil
}

func init() {
	atexit.Register(exithandler)
}

func Debug(a ...interface{}) int{
	
}

func Error(a ...interface{}) int{

}


func SetCliFlag(cliapp cli.App) {
	var vflag *CntFlag
	var wflag *cli.StringSliceFlag
	var aflag *cli.StringSliceFlag
	vflag = &cli.IntFlag{
		Name: "verbose, V"
		Value: 0
	}
	st_verboseflag = vflag
	cliapp.Flags = append(cliapp.Flags,vflag)
	wflag = &cli.StringSliceFlag{
		Name: "log-files",
		Usage: "set write rotate files",
	}
	cliapp.Flags = append(cliapp.Flags, wflag)
	aflag = &cli.StringSliceFlag{
		Name: "log-appends",
		Usage: "set append files"
	}
	cliapp.Flags = append(cliapp.Flags, aflag)
	nostderr := &cli.BoolFlag {
		Name: "log-nostderr"
		Usage: "specified no stdout"
	}
	cliapp.Flags = append(cliapp.Flags, nostderr)	
	return
}

func SetCliFlag(ctx *cli.Context)  error{
	var appfiles []string
	var cfiles []string
	var vmode int
	var lglvl l4g.Level

	if st_logger != nil {
		st_logger.Close()
	}
	st_logger = nil



	vmode = ctx.GlobalInt("verbose")
	if vmode <= 0 {
		lglvl = l4g.ERROR
	} else if vmode == 1 {
		lglvl = l4g.WARNING
	} else if vmode == 2 {
		lglvl = l4g.INFO
	} else if vmode == 3 {
		lglvl = l4g.DEBUG
	} else if vmode == 4 {
		lglvl = l4g.TRACE
	} else if vmode >= 5 {
		lglvl = l4g.FINEST
	}

	st_logger = l4g.NewLogger()
	if ! ctx.GlobalIsSet("log-nostderr") {
		st_logger.AddFilter("stderr", lglvl, l4g.NewStderrLogWriter())
	}

	cfiles = ctx.GlobalStringSlice("log-files")
	if len(cfiles) > 0 {
		for _, f := range cfiles {
			st_logger.AddFilter(f, lglvl, l4g.NewFileLogWriter(f, false))
		}
	}

	appfiles = ctx.GlobalFlagNames("log-appends")
	if len(appfiles) > 0 {
		for _, f := range appfiles {
			st_logger.AddFilter(f, lglvl, l4g.NewFileLogWriter(f, true))
		}
	}
	return nil
}