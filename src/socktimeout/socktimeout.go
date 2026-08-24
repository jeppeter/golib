package socktimeout

import (
	"dbgutil"
	"logutil"
	"net"
	"time"
)

type SockClient struct {
	conn *net.TCPConn
}

type SockAccept struct {
	ln *net.TCPListener
}

func NewSockAccept(typestr string, bindstr string) (retp *SockAccept, err error) {
	var lner net.Listener
	var ok bool
	retp = &SockAccept{}
	lner, err = net.Listen(typestr, bindstr)
	if err != nil {
		retp = nil
		err = dbgutil.FormatError("cannot bind [%s] [%s] for  error %s", typestr, bindstr, err.Error())
		return
	}
	retp.ln, ok = lner.(*net.TCPListener)
	if !ok {
		retp = nil
		err = dbgutil.FormatError("cast to TCPConn error")
		return
	}

	return
}

func (retp *SockAccept) Close() {
	if retp.ln != nil {
		retp.ln.Close()
	}
	retp.ln = nil
	return
}

func (retp *SockAccept) AcceptTimeout(mills int) (nretp *SockClient, err error) {
	var nconn *net.TCPConn
	var tm time.Time
	var curtime time.Time
	nretp = nil
	if mills != 0 {
		tm = time.Now().Add(time.Duration(mills) * time.Millisecond)
	} else {
		tm = time.Time{}
	}

	if retp.ln == nil {
		err = dbgutil.FormatError("closed socket")
		return
	}

	err = retp.ln.SetDeadline(tm)
	if err != nil {
		err = dbgutil.FormatError("set %d timeout error %s", mills, err.Error())
		return
	}

	nconn, err = retp.ln.AcceptTCP()
	if err != nil {
		if mills != 0 {
			curtime = time.Now()
			if curtime.After(tm) {
				nretp = nil
				err = nil
				return
			}
		}
		err = dbgutil.FormatError("accept error %s", err.Error())
		return
	}

	nretp = &SockClient{}
	nretp.conn = nconn
	err = nil
	return
}

func (retp *SockAccept) AcceptTimeoutRaw(mills int) (nretp net.Conn, err error) {
	var nconn *net.TCPConn
	var tm time.Time
	var curtime time.Time
	nretp = nil
	if mills != 0 {
		tm = time.Now().Add(time.Duration(mills) * time.Millisecond)
	} else {
		tm = time.Time{}
	}

	if retp.ln == nil {
		err = dbgutil.FormatError("closed socket")
		return
	}

	err = retp.ln.SetDeadline(tm)
	if err != nil {
		err = dbgutil.FormatError("set %d timeout error %s", mills, err.Error())
		return
	}

	nconn, err = retp.ln.AcceptTCP()
	if err != nil {
		if mills != 0 {
			curtime = time.Now()
			if curtime.After(tm) {
				nretp = nil
				err = nil
				return
			}
		}
		err = dbgutil.FormatError("accept error %s", err.Error())
		return
	}

	nretp = nconn
	err = nil
	return
}

func read_tcp_buffer(conn *net.TCPConn, n int, mills int) (retb []byte, err error) {
	var tm time.Time
	var cbuf []byte
	var retn int
	var curtime time.Time
	var i int
	retb = []byte{}
	if mills != 0 {
		tm = time.Now().Add(time.Duration(mills) * time.Millisecond)
	} else {
		tm = time.Time{}
	}
	err = conn.SetReadDeadline(tm)
	if err != nil {
		err = dbgutil.FormatError("set read deadline error %s", err.Error())
		return
	}

	cbuf = make([]byte, n)
	retn, err = conn.Read(cbuf)
	if err != nil {
		if mills != 0 {
			curtime = time.Now()
			if curtime.After(tm) {
				retb = []byte{}
				err = nil
				return
			}
		}
		err = dbgutil.FormatError("Read Error %s", err.Error())
		return
	}

	if retn > 0 {
		retb = make([]byte, retn)
		for i = 0; i < retn; i += 1 {
			retb[i] = cbuf[i]
		}
		logutil.DebugBuffer(retb, "read bytes %d", n)
	}
	return
}

func write_tcp_buffer(conn *net.TCPConn, wbuf []byte, mills int) (retn int, err error) {
	var _cn int
	var tm time.Time
	retn = 0
	if mills != 0 {
		tm = time.Now().Add(time.Duration(mills) * time.Millisecond)
	} else {
		tm = time.Time{}
	}

	err = conn.SetWriteDeadline(tm)
	if err != nil {
		err = dbgutil.FormatError("set time out error %s", err.Error())
		return
	}
	_cn = 0
	_cn, err = conn.Write(wbuf)
	if err != nil {
		if mills != 0 {
			var curtime time.Time = time.Now()
			if curtime.After(tm) {
				err = nil
				retn = _cn
				return
			}
		}
		err = dbgutil.FormatError("write error %s", err.Error())
		return
	}
	err = nil
	retn = _cn
	return
}

func NewSockClient(typestr string, connstr string, mills int) (retp *SockClient, err error) {
	var conn net.Conn
	var ok bool
	retp = &SockClient{}

	if mills != 0 {
		conn, err = net.DialTimeout(typestr, connstr, time.Duration(mills)*time.Millisecond)
	} else {
		conn, err = net.Dial(typestr, connstr)
	}
	if err != nil {
		err = dbgutil.FormatError("connect %s error %s", connstr, err.Error())
		retp = nil
		return
	}

	retp.conn, ok = conn.(*net.TCPConn)
	if !ok {
		retp = nil
		err = dbgutil.FormatError("cast to TCPConn error")
		return
	}
	err = nil
	return
}

func (retp *SockClient) Close() {
	if retp.conn != nil {
		retp.conn.Close()
	}
	retp.conn = nil
	return
}

func (retp *SockClient) ReadTimeout(n int, mills int) (retb []byte, err error) {
	if retp.conn == nil {
		retb = []byte{}
		err = dbgutil.FormatError("read close socket")
		return
	}
	return read_tcp_buffer(retp.conn, n, mills)
}

func (retp *SockClient) WriteTimeout(wbuf []byte, mills int) (retn int, err error) {
	if retp.conn == nil {
		retn = 0
		err = dbgutil.FormatError("write close socket")
		return

	}
	return write_tcp_buffer(retp.conn, wbuf, mills)
}
