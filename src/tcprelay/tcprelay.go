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
}

type RelayConn struct {
	localchl    net.Conn
	svrchl      net.Conn
	noteexited  int
	localrdexit chan int
	localwrexit chan int
	rdexited    int
	wrexited    int
}

func (retp *RelayConn) NotifyExit(bwait bool) (exited int) {
	if retp.noteexited != 0 {
		retp.localrdexit <- 1
		retp.localwrexit <- 1
		retp.noteexited = 1
	}

	for bwait == true {
		if retp.rdexited == 0 || retp.wrexited == 0 {
			/*we exited*/
			time.Sleep(time.Millisecond * 100)
		}
	}

	exited = 0
	if retp.rdexited != 0 && retp.wrexited != 0 {
		exited = 1
	}
	return
}

func (retp *RelayConn) Close() {
	retp.NotifyExit(true)
	retp.localchl.Close()
	retp.svrchl.Close()
	return
}

func (retp *RelayConn) LocalReadProc() {
	var rdata []byte
	var n int
	var err error
	var tlen int
	var wn int
	var needexited bool
	rdata = make([]byte, 2500)
outer_loop:
	for {
		n, err = retp.localchl.Read(rdata)
		if err != nil {
			err = dbgutil.FormatError("read local data error [%s]", err.Error())
			break
		}
		tlen = 0
		for tlen < n {
			wn, err = retp.svrchl.Write(rdata[tlen:n])
			if err != nil {
				err = dbgutil.FormatError("write server data [%d] error [%s]")
				break outer_loop
			}
			tlen += wn
		}

		needexited = false
		select {
		case <-retp.wrexited:
			needexited = true
		case <-retp.localrdexit:
			needexited = true
		case <-time.After(10 * time.Millisecond):
			n = n
		}

		if needexited || retp.wrexited != 0 {
			err = nil
			break
		}
	}

	/*we exit*/
	retp.rdexited = 1
}

func (retp *RelayConn) RemoteReadProc() {
	var rdata []byte
	var n int
	var err error
	var tlen int
	var wn int
	var needexited bool
	rdata = make([]byte, 2500)
outer_loop:
	for {
		n, err = retp.svrchl.Read(rdata)
		if err != nil {
			err = dbgutil.FormatError("read server data error [%s]", err.Error())
			break
		}
		tlen = 0
		for tlen < n {
			wn, err = retp.localchl.Write(rdata[tlen:n])
			if err != nil {
				err = dbgutil.FormatError("write local data [%d] error [%s]")
				break outer_loop
			}
			tlen += wn
		}

		needexited = false
		select {
		case <-retp.localwrexit:
			needexited = true
		case <-time.After(10 * time.Millisecond):
			n = n
		}

		if needexited || retp.rdexited != 0 {
			err = nil
			break
		}
	}

	/*we exit*/
	retp.wrexited = 1
}
