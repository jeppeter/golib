package tcprelay

import (
	"dbgutil"
	"errors"
	"logutil"
	"net"
	"os"
	"reflect"
	"slices"
	"socktimeout"
	"time"
)

type RemoteConn interface {
	Open() (conn net.Conn, err error)
	ReadHandle(inbyte []byte) (outbytes []byte, err error)
	WriteHandle(inbyte []byte) (outbytes []byte, err error)
	Close()
}

type CreateConn interface {
	Create(conn net.Conn) (retconn RemoteConn, err error)
	Close()
}

type DefaultRemoteConn struct {
	bindstr string
}

func NewDefaultRemoteConn(bindstr string) (retp *DefaultRemoteConn, err error) {
	retp = &DefaultRemoteConn{}
	retp.bindstr = bindstr
	err = nil
	return
}

func (retp *DefaultRemoteConn) Open() (conn net.Conn, err error) {
	return net.Dial("tcp", retp.bindstr)
}

func (retp *DefaultRemoteConn) ReadHandle(inbyte []byte) (outbytes []byte, err error) {
	outbytes = inbyte
	err = nil
	return
}

func (retp *DefaultRemoteConn) WriteHandle(inbyte []byte) (outbytes []byte, err error) {
	outbytes = inbyte
	err = nil
	return
}

func (retp *DefaultRemoteConn) Close() {
	return
}

type DefaultCreateConn struct {
	remotestr string
}

func NewDefaultCreateConn(remotestr string) (retp *DefaultCreateConn, err error) {
	retp = &DefaultCreateConn{}
	retp.remotestr = remotestr
	err = nil
	return
}

func (retp *DefaultCreateConn) Create(conn net.Conn) (retconn RemoteConn, err error) {
	var dret *DefaultRemoteConn
	dret, err = NewDefaultRemoteConn(retp.remotestr)
	if err != nil {
		return
	}
	retconn = dret
	err = nil
	return
}

func (retp *DefaultCreateConn) Close() {
	return
}

type RelayConn struct {
	localchl    net.Conn
	svrchl      net.Conn
	remotehdl   RemoteConn
	noteexited  int
	localrdexit chan int
	localwrexit chan int
	outnotechl  chan int
	rdexited    int
	wrexited    int
}

func NewRelayConn(localchl net.Conn, remotehdl RemoteConn) (retp *RelayConn, err error) {
	retp = &RelayConn{}
	retp.localchl = localchl
	retp.svrchl = nil
	retp.remotehdl = remotehdl
	retp.noteexited = 1
	/*we make 10 size for not hanging when send the remote*/
	retp.localwrexit = make(chan int, 10)
	retp.localrdexit = make(chan int, 10)
	retp.outnotechl = make(chan int, 10)
	retp.rdexited = 1
	retp.wrexited = 1

	retp.svrchl, err = retp.remotehdl.Open()
	if err != nil {
		retp.Close()
		retp = nil
		return
	}

	/*now we pretend to go for ok*/
	retp.noteexited = 0
	retp.wrexited = 0
	retp.rdexited = 0

	/*now to make sure */
	go retp.LocalReadProc()
	go retp.RemoteReadProc()

	err = nil
	return

}

func (retp *RelayConn) NotifyExit(bwait bool) (exited int) {
	if retp.noteexited == 0 {
		retp.localrdexit <- 1
		retp.localwrexit <- 1
		retp.noteexited = 1
	}

	for bwait {
		if retp.rdexited == 0 || retp.wrexited == 0 {
			/*we exited*/
			time.Sleep(time.Millisecond * 10)
		} else if retp.rdexited != 0 && retp.wrexited != 0 {
			break
		}
	}

	exited = 0
	if retp.rdexited != 0 && retp.wrexited != 0 {
		exited = 1
	}
	return
}

func (retp *RelayConn) Close() {
	for {
		if retp.NotifyExit(true) != 0 {
			break
		}
	}
	retp.localchl.Close()
	retp.svrchl.Close()
	retp.remotehdl.Close()
	return
}

func (retp *RelayConn) Exited() bool {
	if retp.wrexited != 0 && retp.rdexited != 0 {
		return true
	}
	return false
}

func (retp *RelayConn) LocalReadProc() {
	var rdata []byte
	var outbytes []byte
	var n int
	var err error
	var tlen int
	var wn int
	var needexited bool
	var tm time.Time
	rdata = make([]byte, 2048)
outer_loop:
	for {
		tm = time.Now().Add(time.Duration(500) * time.Millisecond)
		err = retp.localchl.SetReadDeadline(tm)
		if err != nil {
			logutil.Error("localchl SetReadDeadline error %s", err.Error())
			break
		}

		n, err = retp.localchl.Read(rdata)
		if err != nil {
			if !errors.Is(err, os.ErrDeadlineExceeded) {
				logutil.Error("read local data error [%s]", err.Error())
				break
			}
			n = 0
		}
		if n > 0 {
			outbytes, err = retp.remotehdl.WriteHandle(rdata[:n])
			if err != nil {
				break
			}

			tlen = 0
			for tlen < len(outbytes) {
				wn, err = retp.svrchl.Write(outbytes[tlen:])
				if err != nil {
					err = dbgutil.FormatError("write server data [%d] error [%s]", tlen, err.Error())
					break outer_loop
				}
				tlen += wn
			}
		}

		needexited = false
		select {
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

	logutil.Debug("RelayConn Read Exit")
	retp.outnotechl <- 1

	/*we exit*/
	retp.rdexited = 1
}

func (retp *RelayConn) RemoteReadProc() {
	var rdata []byte
	var n int
	var err error
	var tlen int
	var wn int
	var outbytes []byte
	var needexited bool
	var tm time.Time
	rdata = make([]byte, 2500)
outer_loop:
	for {
		/*first to */
		tm = time.Now().Add(time.Duration(500) * time.Millisecond)
		err = retp.svrchl.SetReadDeadline(tm)
		if err != nil {
			logutil.Error("svrhdl SetReadDeadline error %s", err.Error())
			break
		}
		n, err = retp.svrchl.Read(rdata)
		if err != nil {
			if !errors.Is(err, os.ErrDeadlineExceeded) {
				logutil.Error("read server data error [%s]", err.Error())
				break
			}
			/*we make sure not read any data*/
			n = 0
		}

		if n > 0 {
			outbytes, err = retp.remotehdl.ReadHandle(rdata[:n])
			if err != nil {
				break
			}

			tlen = 0
			for tlen < len(outbytes) {
				wn, err = retp.localchl.Write(outbytes[tlen:])
				if err != nil {
					err = dbgutil.FormatError("write local data [%d] error [%s]", tlen, err.Error())
					break outer_loop
				}
				tlen += wn
			}
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

	logutil.Debug("RelayConn Write Exit")
	retp.outnotechl <- 1
	/*we exit*/
	retp.wrexited = 1
}

type RelayListen struct {
	localstr     string
	remoteconn   CreateConn
	mainsvr      *socktimeout.SockAccept
	chlds        []*RelayConn
	noteexited   int
	noteexitchl  chan int
	exited       int
	tickexited   int
	tickexitchan chan int
	tickchan     chan int
}

func (retp *RelayListen) TimeTick(mills int) {
	var exited int = 0
	for {
		select {
		case <-retp.tickexitchan:
			exited = 1
		case <-time.After(time.Duration(mills) * time.Millisecond):
			retp = retp
		}

		if exited != 0 {
			break
		}

		/*we send this to make sure */
		retp.tickchan <- 1
	}

	retp.tickexited = 1
}

func (retp *RelayListen) MainProc() {
	var nconn RemoteConn = nil
	var netconn net.Conn = nil
	var cconn *RelayConn = nil
	var selcases []reflect.SelectCase = []reflect.SelectCase{}
	var basesel reflect.SelectCase = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(retp.tickchan),
	}
	var exitreflect reflect.SelectCase = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(retp.noteexitchl),
	}
	var selc reflect.SelectCase
	var chosen int
	var ok bool
	var err error
	var idx int

	/*we make sure 500 */
	retp.tickexited = 0
	go retp.TimeTick(500)
	selcases = append(selcases, basesel)
	selcases = append(selcases, exitreflect)

	for {
		if retp.mainsvr == nil {
			retp.mainsvr, err = socktimeout.NewSockAccept("tcp", retp.localstr)
			if err != nil {
				retp.mainsvr = nil
				logutil.Error("NewSockAccept error %s", err.Error())
			}
		}
		err = nil
		netconn = nil

		if retp.mainsvr != nil {
			netconn, err = retp.mainsvr.AcceptTimeoutRaw(300)
		}
		if err == nil {
			if netconn != nil {
				if retp.remoteconn != nil {
					nconn, err = retp.remoteconn.Create(netconn)
					if err == nil {
						/*now we should */
						cconn, err = NewRelayConn(netconn, nconn)
						if err == nil {
							retp.chlds = append(retp.chlds, cconn)
							/*now to update selcases*/
							selc = reflect.SelectCase{
								Dir:  reflect.SelectRecv,
								Chan: reflect.ValueOf(cconn.outnotechl),
							}
							selcases = append(selcases, selc)
						} else {
							nconn.Close()
							netconn.Close()
							nconn = nil
							netconn = nil
							logutil.Error("NewRelayConn error %s", err.Error())
						}

					} else {
						logutil.Error("call_remote error %s", err.Error())
						netconn.Close()
						netconn = nil
					}
				} else {
					netconn.Close()
					logutil.Error("no call_remote")
				}
			}
		} else {
			logutil.Error("accept Error %s", err.Error())
			retp.mainsvr.Close()
			retp.mainsvr = nil
		}

		/*now we should give the chlds wait exited*/
		chosen, _, ok = reflect.Select(selcases)
		if ok && chosen > 1 && chosen < len(selcases) {
			/*ok we should remove */
			logutil.Debug("chosen %d children to remove", chosen-2)
			cconn = retp.chlds[chosen-2]
			retp.chlds = slices.Delete(retp.chlds, chosen-2, chosen-1)
			cconn.Close()
			cconn = nil
			selcases = slices.Delete(selcases, chosen, chosen+1)
		} else if ok && chosen == 1 {
			logutil.Debug("NotifyExit MainProc")
			break
		}
	}

	for idx = 0; idx < len(retp.chlds); idx += 1 {
		/*now to notified to remove so we do this*/
		retp.chlds[idx].NotifyExit(false)
	}

	/*now to make sure tick routine to remove first*/
	retp.tickexitchan <- 1
	for {
		if retp.tickexited != 0 {
			break
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
	}

	for len(retp.chlds) > 0 {
		for {
			if retp.chlds[0].NotifyExit(false) != 0 {
				break
			}
			time.Sleep(time.Duration(10) * time.Millisecond)
		}
		/*remove the first*/
		retp.chlds = slices.Delete(retp.chlds, 0, 1)
	}

	retp.exited = 1
}

func NewRelayListen(localstr string, rconn CreateConn) (retp *RelayListen, err error) {
	retp = &RelayListen{}
	retp.localstr = localstr
	retp.remoteconn = rconn
	retp.chlds = []*RelayConn{}
	retp.noteexited = 0
	retp.noteexitchl = make(chan int, 10)
	retp.exited = 1
	retp.tickexited = 1
	retp.tickchan = make(chan int, 10)
	retp.tickexitchan = make(chan int, 10)
	retp.mainsvr, err = socktimeout.NewSockAccept("tcp", localstr)
	if err != nil {
		retp.Close()
		retp = nil
		err = dbgutil.FormatError("listen on [%s] error [%s]", localstr, err.Error())
		return
	}
	err = nil
	return
}

func (retp *RelayListen) Start() (err error) {
	retp.exited = 0
	go retp.MainProc()
	err = nil
	return
}

func (retp *RelayListen) NotifyExit(bwait bool) (exited int) {
	if retp.exited == 0 {
		if retp.noteexited == 0 {
			retp.noteexitchl <- 1
			retp.noteexited = 1
		}
	}

	for bwait {
		if retp.exited != 0 {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}

	exited = 0
	if retp.exited != 0 {
		exited = 1
	}
	return
}

func (retp *RelayListen) Close() {

	for {
		if retp.NotifyExit(true) != 0 {
			break
		}
	}

	if len(retp.chlds) != 0 {
		panic("not free all chlds")
	}
	if retp.mainsvr != nil {
		retp.mainsvr.Close()
	}
	return
}
