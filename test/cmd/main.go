package main

import (
	"flag"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"reflect"
)

type CntFlag struct {
	cli.BoolFlag
	CntValue int
}

func (c *CntFlag) Apply(f *flag.FlagSet) {
	c.ApplyWithError(f)
}

func (c *CntFlag) ApplyWithError(f *flag.FlagSet) error {
	c.CntValue += 1
	return nil
}

func (c *CntFlag) GetName() string {
	return c.Name
}

func NewCountFlag(longname string, shortname interface{}, usage interface{}) *CntFlag {
	c := &CntFlag{}
	namestr := longname
	if shortname != nil {
		namestr += fmt.Sprintf(", %s", shortname)
	}
	c.Name = namestr

	if usage != nil {
		c.Usage = fmt.Sprintf("%v", usage)
	} else {
		c.Usage = fmt.Sprintf("%s set", longname)
	}
	c.CntValue = 0
	return c
}

func addVerboseCnt(app *cli.App) {
	c := NewCountFlag("verbose", "v", "verbose mode")
	fmt.Fprintf(os.Stdout, "name [%s]\n", c.Name)
	app.Flags = append(app.Flags, c)
	return
}

func makeTestCmd() cli.Command {
	appcmd := cli.Command{
		Name:  "test",
		Usage: "test command",
	}

	appcmd.Action = func(ctx *cli.Context) error {
		flags := ctx.App.Flags
		for i := 0; i < len(flags); i++ {
			c := flags[i]
			if c.GetName() == "verbose, v" {
				fmt.Fprintf(os.Stderr, "get verbose [%d]\n", i)
			}
		}
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
		flags := ctx.Command.VisibleFlags()
		for i := 0; i < len(flags); i++ {
			c := flags[i]
			fmt.Printf("[%d]=[%s] name [%s]\n", i, reflect.TypeOf(c), c.GetName())
		}
		return nil
	}
	return appcmd
}

func main() {
	app := cli.NewApp()
	addVerboseCnt(app)

	app.Commands = append(app.Commands, makeTestCmd())
	app.Commands = append(app.Commands, makeCmdCmd())
	app.Run(os.Args)
	return
}
