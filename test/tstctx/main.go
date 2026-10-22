package main

import (
	"context"
	"dbgutil"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"logutil"
	"net"
	"strconv"
	"time"
)

func init() {
	Mergectx_handler(nil, nil, nil)
}

func merge_ctx(port int, mills int) {
	readFromConn := func(ctx context.Context, conn net.Conn, b []byte) (n int, err error) {
		stopc := make(chan struct{})
		stop := context.AfterFunc(ctx, func() {
			logutil.Debug("after func")
			conn.SetReadDeadline(time.Now())
			close(stopc)
		})
		n, err = conn.Read(b)
		if !stop() {
			// The AfterFunc was started.
			// Wait for it to complete, and reset the Conn's deadline.
			logutil.Debug("stop false")
			<-stopc
			conn.SetReadDeadline(time.Time{})
			return n, ctx.Err()
		}
		logutil.Debug("stopped n %d err %s", n, err.Error())
		return n, err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	conn, err := net.Dial(listener.Addr().Network(), listener.Addr().String())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mills)*time.Millisecond)
	defer cancel()

	b := make([]byte, 1024)
	_, err = readFromConn(ctx, conn, b)
	fmt.Println(err)

	// Output:
	// context deadline exceeded

}

func Mergectx_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var args []string
	var port int = 9090
	var mills int = 1
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	args = ns.GetArray("subnargs")
	if len(args) > 0 {
		port, err = strconv.Atoi(args[0])
		if err != nil {
			err = dbgutil.FormatError("%s not port %s", args[0], err.Error())
			return
		}
	}

	if len(args) > 1 {
		mills, err = strconv.Atoi(args[1])
		if err != nil {
			if err != nil {
				err = dbgutil.FormatError("%s not mills %s", args[1], err.Error())
				return
			}
		}
	}

	merge_ctx(port, mills)
	return nil
}

func LoadRegCmdFlags(parser *extargsparse.ExtArgsParse) (err error) {
	var commandline_fmt string
	var commandline string
	commandline_fmt = `{
		"mergectx<Mergectx_handler>##[port] [timemills] to make merge ctx default port 9090 timemills 1##" : {
			"$" : "*"
		}

	}`

	commandline = fmt.Sprintf(commandline_fmt)
	err = parser.LoadCommandLineString(commandline)
	return
}

func main() {
	var parser *extargsparse.ExtArgsParse
	var err error
	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		logutil.Error("%s", err.Error())
		atexit.Exit(5)
	}

	err = LoadRegCmdFlags(parser)
	if err != nil {
		logutil.Error("%s", err.Error())
		atexit.Exit(5)
	}

	err = logutil.PrepareLog(parser)
	if err != nil {
		logutil.Error("%s", err.Error())
		atexit.Exit(5)
	}

	_, err = parser.ParseCommandLine(nil, nil)
	if err != nil {
		logutil.Error("%s", err.Error())
		atexit.Exit(4)
	}
	atexit.Exit(0)
}
