package tcprelay

import (
	"dbgutil"
	"logutil"
	"net"
)

type RelayListen struct {
	localstr string
	svrstr   string
	mainsvr  net.Conn
	chlds    []ProxyConn
	func before_call(newconn *RelayConn) (err error)
}

type RelayConn struct {
	localchl  net.Conn
	svrchl    net.Conn
	localrdexit   chan int
	localwrexit   chan int
	rdexited chan int
	wrexited chan int

}

func (retp *RelayConn) LocalReadProc() {
	var rdata []byte
	var n int
	var err error
	var tlen int
	var wn int
	rdata = make([]byte, 2500)
outer_loop:
	for {
		n, err = retp.localchl.Read(rdata)
		if err != nil {
			err = dbgutil.FormatError("read data error [%s]", err.Error())
			break
		}
		tlen = 0
		for tlen < n {
			wn, err = retp.svrchl.Write(rdata[tlen:n])
			if err != nil {
				err = dbgutil.FormatError("write data [%d] error [%s]")
				break outer_loop
			}
			tlen += wn
		}
		select {
		case <- retp.localwrexit:
			break
		case <- time.
		}
	}

}

func (retp *RelayConn) RemoteReadProc() {
	for {
		/*we do not */
	}
}
