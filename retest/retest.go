package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"regexp"
)

func makeFindAllCommand() cli.Command {
	cmd := cli.Command{}
	cmd.Name = "findall"
	cmd.ShortName = "fa"
	cmd.Usage = fmt.Sprintf(" restr instrs...")
	cmd.Action = func(c *cli.Context) {
		if len(c.Args()) < 2 {
			fmt.Fprintf(os.Stderr, "findall %s", cmd.Usage)
			os.Exit(4)
		}
		reg, err := regexp.Compile(c.Args()[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "can not compile[%s] %s", c.Args()[0], err.Error())
			os.Exit(5)
		}
		for i := 1; i < len(c.Args()); i++ {
			matchstrings := reg.FindStringSubmatch(c.Args()[i])
			if len(matchstrings) > 0 {
				fmt.Fprintf(os.Stdout, "[%s] find all in [%s]\n", c.Args()[0], c.Args()[i])
				for j, s := range matchstrings {
					fmt.Fprintf(os.Stdout, "\t[%d] [%s]\n", j, s)
				}
			} else {
				fmt.Fprintf(os.Stdout, "[%s] not find all in [%s]\n", c.Args()[0], c.Args()[i])
			}
		}
		os.Exit(0)
	}

	return cmd
}

func main() {
	app := cli.NewApp()
	app.Version = "1.0.2"
	app.Commands = append(app.Commands, makeFindAllCommand())
	app.Run(os.Args)
}
