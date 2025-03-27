package main

import (
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"logutil"
	"net"
	"os"
	"strconv"
)

func init() {
	Tcpconn_handler(nil, nil, nil)
	Tcplisten_handler(nil, nil, nil)
	Udpsend_handler(nil, nil, nil)
	Udprecv_handler(nil, nil, nil)
}

func Tcpconn_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var i int
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")

	for i, _ = range sarr {
		_, err = net.Dial("tcp", sarr[i])
		if err != nil {
			logutil.Error("connect %s error %s", sarr[i], err.Error())
			return
		}
		fmt.Printf("connect %s succ\n", sarr[i])
	}

	err = nil
	return
}

func Tcplisten_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var liststr string = ":40008"
	var ln net.Listener
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")

	if len(sarr) > 0 {
		liststr = sarr[0]
	}

	ln, err = net.Listen("tcp", liststr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not listen on [%s] [%s]\n", liststr, err.Error())
		os.Exit(4)
	}
	logutil.Debug("listen on %s", liststr)
	for {
		_, err = ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot accept %s\n", err.Error())
			continue
		}

		fmt.Fprintf(os.Stdout, "accept %s ok\n", liststr)
	}

	err = nil
	return
}

func Udpsend_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var outbytes []byte
	var inbytes []byte
	var server string = "127.0.0.1:40008"
	var infile string
	var conn net.Conn
	var tlen int
	var r int

	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}
	infile = ns.GetString("input")
	logutil.Debug("read %s", infile)
	outbytes, err = fileop.ReadFileBytes(infile)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")

	if len(sarr) > 0 {
		server = sarr[0]
	}

	conn, err = net.Dial("udp", server)
	if err != nil {
		return
	}
	logutil.Debug("dial [%s]", server)

	tlen = 0

	for tlen < len(outbytes) {
		r, err = conn.Write(outbytes[tlen:])
		if err != nil {
			return
		}
		tlen += r
	}

	tlen = 0
	inbytes = make([]byte, 100000)
	r, err = conn.Read(inbytes)
	if err != nil {
		return
	}
	logutil.DebugBuffer(inbytes[:r], "read %s bytes", server)

	err = nil
	return
}

func Udprecv_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var liststr string = ":40008"
	var lhost string
	var lport string
	var liport int
	var ludp net.UDPAddr
	var ser *net.UDPConn
	var rbuf []byte
	var r int
	var raddr *net.UDPAddr
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")

	if len(sarr) > 0 {
		liststr = sarr[0]
	}

	lhost, lport, err = net.SplitHostPort(liststr)
	if err != nil {
		logutil.Error("parse [%s] error [%s]", liststr, err.Error())
		return
	}
	liport, err = strconv.Atoi(lport)
	if err != nil {
		return
	}
	ludp = net.UDPAddr{
		Port: liport,
		IP:   net.ParseIP(lhost),
	}

	ser, err = net.ListenUDP("udp", &ludp)
	if err != nil {
		logutil.Error("listen on [%s] error [%s]", liststr, err.Error())
		return
	}

	rbuf = make([]byte, 2000)
	logutil.Debug("listen udp on %s", liststr)
	for {
		r, raddr, err = ser.ReadFromUDP(rbuf)
		if err != nil {
			logutil.Error("read error [%s]", err.Error())
			continue
		}
		logutil.DebugBuffer(rbuf[:r], "read from %v ", raddr)
		_, err = ser.WriteToUDP(rbuf[:r], raddr)
		if err != nil {
			logutil.Error("write [%v] error[%s]", raddr, err.Error())
			continue
		}
		logutil.Debug("write [%v] succ", raddr)
	}

	err = nil
	return
}

func main() {
	var commandline string
	var err error
	var parser *extargsparse.ExtArgsParse
	var ns *extargsparse.NameSpaceEx

	commandline = `{
		"timeout|t" : 500,
		"input|i" : null,
		"output|o" : null,
		"tcpconn<Tcpconn_handler>##host:port ... to connect##" : {
			"$" : "+"
		},
		"tcplisten<Tcplisten_handler>##[:port] to listen##" : {
			"$" : "?"
		},
		"udpsend<Udpsend_handler>##host:port to send from input##" : {
			"$" : 1
		},
		"udprecv<Udprecv_handler>##:port to recv from input##" : {
			"$" : 1
		}

	}`

	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not make parser err[%s]\n", err.Error())
		atexit.Exit(5)
	}

	err = logutil.PrepareLog(parser)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not set [%s]\n", err.Error())
		atexit.Exit(5)
	}
	err = parser.LoadCommandLineString(commandline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not parse %s\n", commandline)
		atexit.Exit(5)
	}

	ns, err = parser.ParseCommandLineEx(nil, nil, nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not use parse command line [%s]\n", err.Error())
		atexit.Exit(4)
	}
	if len(ns.GetString("subcommand")) == 0 {
		fmt.Fprintf(os.Stderr, "can not get subcommand\n")
		atexit.Exit(5)
	}
	fmt.Fprintf(os.Stdout, "subcommand [%s] succ\n", ns.GetString("subcommand"))
	atexit.Exit(0)
	return
}
