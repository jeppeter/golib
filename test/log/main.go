package main

import (
	"github.com/codegangsta/cli"
	"github.com/tebeka/atexit"
	"os"
	"reflect"
)

func makeTestCmd() cli.Command {
	appcmd := cli.Command{
		Name:  "test",
		Usage: "test command",
	}

	appcmd.Action = func(ctx *cli.Context) error {
		SetCliFlag(ctx)
		flags := ctx.App.Flags
		for i := 0; i < len(flags); i++ {
			c := flags[i]
			if c.GetName() == "verbose, v" {
				Error("get verbose [%d]", i)
			}
		}
		atexit.Exit(0)
		return nil
	}
	return appcmd
}

func makeCmdCmd() cli.Command {
	appcmd := cli.Command{
		Name:  "cmd",
		Usage: "cmd command",
	}

	appcmd.Action = func(ctx *cli.Context) error {
		SetCliFlag(ctx)
		flags := ctx.Command.VisibleFlags()
		for i := 0; i < len(flags); i++ {
			c := flags[i]
			Debug("[%d]=[%s] name [%s]", i, reflect.TypeOf(c), c.GetName())
		}
		atexit.Exit(0)
		return nil
	}
	return appcmd
}

func main() {
	cliapp := cli.NewApp()
	AddCliFlag(cliapp)
	cliapp.Commands = append(cliapp.Commands, makeTestCmd())
	cliapp.Commands = append(cliapp.Commands, makeCmdCmd())
	cliapp.Run(os.Args)
}
