package main

import (
	"dbgutil"
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"logutil"
	"sockproto"
	"socktimeout"
	"tcprelay"
)

func init() {
	Sockserver_handler(nil, nil, nil)
	Sockclient_handler(nil, nil, nil)
	Tcprelay_handler(nil, nil, nil)
}

func LoadParser(parser *extargsparse.ExtArgsParse) (err error) {
	var commandline_fmt string
	var commandline string
	commandline_fmt = `{
		"timeout|t" : 5000,
		"socksvr<Sockserver_handler>##[:9090] to listen on sock to read##" : {
			"$" : "*"
		},
		"sockcli<Sockclient_handler>##127.0.0.1:9090 files ... to send files##" : {
			"$" : "+"
		},
		"tcprelay<Tcprelay_handler>##bindstr remotestr to make relay##" : {
			"$" : 2
		}
	}`

	commandline = fmt.Sprintf(commandline_fmt)
	err = parser.LoadCommandLineString(commandline)
	return
}

func server_handler(cli *sockproto.SockChannel, mills int) {
	var rets string
	var err error
	for {
		rets, err = cli.ReadJson(mills)
		if err != nil {
			logutil.Error("%s", err.Error())
			cli.Close()
			return
		}

		if len(rets) > 0 {
			logutil.Debug("read\n%s", rets)
			err = cli.WriteJson(rets, mills)
			if err != nil {
				logutil.Error("%s", err.Error())
				cli.Close()
				return
			}
		}
	}
}

func Sockserver_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var bindstr string = ":9090"
	var svr *socktimeout.SockAccept = nil
	var cli *socktimeout.SockClient = nil
	var timeout int
	var chl *sockproto.SockChannel = nil
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) > 0 {
		bindstr = sarr[0]
	}
	timeout = ns.GetInt("timeout")

	defer func() {
		if svr != nil {
			svr.Close()
		}
	}()

	svr, err = socktimeout.NewSockAccept("tcp", bindstr)
	if err != nil {
		return
	}

	logutil.Debug("listen on %s", bindstr)

	for {
		cli, err = svr.AcceptTimeout(timeout)
		if err != nil {
			return
		}
		if cli != nil {
			chl, err = sockproto.NewSockChannel(cli)
			if err == nil {
				go server_handler(chl, timeout)
				chl = nil
				cli = nil
			} else {
				cli.Close()
				chl = nil
				cli = nil
				logutil.Error("%s", err.Error())
			}
		}
	}

	return
}

func Sockclient_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var i int
	var fname string
	var fstr string
	var cli *socktimeout.SockClient = nil
	var chl *sockproto.SockChannel = nil
	var rets string
	var timeout int
	var connstr string

	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}
	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = dbgutil.FormatError("need connstr file")
		return
	}
	timeout = ns.GetInt("timeout")
	connstr = sarr[0]

	cli, err = socktimeout.NewSockClient("tcp", connstr, timeout)
	if err != nil {
		return
	}

	defer func() {
		if cli != nil {
			cli.Close()
		}
		cli = nil

		if chl != nil {
			chl.Close()
		}
		chl = nil
	}()

	chl, err = sockproto.NewSockChannel(cli)
	if err != nil {
		return
	}

	cli = nil

	for i = 1; i < len(sarr); i += 1 {
		fname = sarr[i]
		fstr, err = fileop.ReadFile(fname)
		if err != nil {
			return
		}

		err = chl.WriteJson(fstr, timeout)
		if err != nil {
			return
		}
		logutil.Debug("write [%s]", fname)
		rets, err = chl.ReadJson(timeout)
		if err != nil {
			return
		}
		if len(rets) == 0 {
			err = dbgutil.FormatError("no data return for %s", fname)
			return
		}
		logutil.Debug("read json \n%s", rets)
	}

	err = nil
	return
}

func Tcprelay_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var i int
	var fname string
	var fstr string
	var cli *socktimeout.SockClient = nil
	var chl *sockproto.SockChannel = nil
	var rets string
	var timeout int
	var connstr string

	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}
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

	err = LoadParser(parser)
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
