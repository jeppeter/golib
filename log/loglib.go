package main

import (
	"github.com/codegangsta/cli"
	l4g "github.com/jeppeter/log4go"
	"fmt"
	"flag"
)

type CntFlag struct {
	cli.BoolFlag
}

var st_verboseflag *CntFlag
st_verboseflag = nil

func NewCntFlag(longname string,usage interface,shortname interface) *CntFlag {
	var name string
	var usagestr string
	name = longname
	switch shortname.type {
	case string:
		name += fmt.Sprintf(",%s", shortname)
	default:
		name += ""
	}

	switch usage.type {
	case string:
		usagestr = usage
	default:
		usagestr = fmt.Sprintf("set %s count value",longname)
	}

	flag := &CntFlag{
		Name: name,
		Usage: usagestr,
		Value: 0
	}
	return flag
}

func (f *CntFlag) Apply(set *flag.FlagSet) {
	f.ApplyWithError(set)
}

func (f *CntFlag) ApplyWithError(set *flag.FlagSet) error{
	f.Value += 1
	return nil
}

func (f *CntFlag) String() string {
	return fmt.Sprintf("%d",f.Value)
}

func (f *CntFlag) GetName() string {
	return fmt.Sprintf("%s", f.Name)
}

func SetCliFlag(cliapp cli.App) {
	var vflag *CntFlag
	var wflag *cli.StringSliceFlag
	var aflag *cli.StringSliceFlag
	vflag = NewCntFlag("verbose","verbose mode set", "V")
	st_verboseflag = vflag
	cliapp.Flags = append(cliapp.Flags,vflag)
	wflag = &cli.StringFlag{
		Name: "log-files",
		Usage: "set write rotate files",
	}
	cliapp.Flags = append(cliapp.Flags, wflag)
	aflag = &cli.StringSliceFlag{
		Name: "log-appends",
		Usage: "set append files"
	}
	cliapp.Flags = append(cliapp.Flags, aflag)
	return
}

func SetCliFlag(ctx *cli.Context) {
	
}