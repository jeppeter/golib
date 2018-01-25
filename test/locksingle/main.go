package main

import (
	"github.com/codegangsta/cli"
	"github.com/tebeka/atexit"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	var err error = nil
	app := cli.NewApp()
	AddCliFlag(app)
	defer func() {
		if err == nil {
			atexit.Exit(0)
		} else {
			atexit.Exit(5)
		}
	}()
	app.Action = func(ctx *cli.Context) error {
		SetCliFlag(ctx)
		var ptrs []*SingleLock = []*SingleLock{}
		var curptr *SingleLock = nil
		var err2 error
		var cont bool
		var i int
		var a string
		sigch := make(chan os.Signal, 1)
		if runtime.GOOS == "windows" {
			signal.Notify(sigch, os.Interrupt)
		} else {
			signal.Notify(sigch, syscall.SIGINT)
		}

		for i, a = range ctx.Args() {
			curptr, err2 = lock_single(a)
			if err2 != nil {
				Error("can not lock single [%s] error[%s]", a, err2.Error())
				return err2
			}
			Debug("lock[%d]=[%s]", i, a)
			ptrs = append(ptrs, curptr)
		}

		cont = true
		for cont {
			select {
			case <-sigch:
				cont = false
			}
		}

		for i, curptr = range ptrs {
			unlock_single(curptr)
			Debug("unlock[%d] [%s]", i, ctx.Args()[i])
		}
		ptrs = []*SingleLock{}

		cont = true
		for cont {
			select {
			case <-sigch:
				cont = false
			}
		}
		return nil
	}
	app.Run(os.Args)
	return
}
