package npipepack

import (
	"fmt"
	"github.com/jeppeter/npipe"
	"net"
	"time"
)

type NpipeSock struct {
	clipipe    *npipe.PipeConn
	accpipe    *npipe.PipeConn
	listenpipe *npipe.PipeListener
	isaccept   bool
	isserver   bool
	data       []byte
	readmills  int
}

func BindPipe(pipename string, timemills int) (p *NpipeSock, err error) {
	p = &NpipeSock{}
	p.clipipe = nil
	p.accpipe = nil
	p.listenpipe = nil
	p.isaccept = true
	p.isserver = true
	p.data = []byte{}
	p.readmills = timemills

	p.listenpipe, err = npipe.Listen(pipename)
	if err != nil {
		p = nil
		return
	}
	err = nil
	return
}

func ConnPipe(pipename string, timemills int, conntimes int) (p *NpipeSock, err error) {
	p = &NpipeSock{}
	p.clipipe = nil
	p.accpipe = nil
	p.listenpipe = nil
	p.isaccept = false
	p.isserver = false
	p.data = []byte{}
	p.readmills = timemills

	p.clipipe, err = npipe.DialTimeout(pipename, time.Duration(conntimes)*time.Millisecond)
	if err != nil {
		p = nil
		return
	}
	err = nil
	return
}

func (p *NpipeSock) AcceptTimeout() (retp *NpipeSock, err error) {
	var cli *npipe.PipeConn
	var neterr net.Error
	var ok bool

	if !p.isaccept || p.listenpipe == nil {
		err = fmt.Errorf("not accept one")
		return
	}

	cli, err = p.listenpipe.AcceptTimeout(p.readmills)
	if err != nil {
		neterr, ok = err.(net.Error)
		if ok && neterr.Timeout() {
			retp = nil
			err = nil
		}
		return
	}

	retp = &NpipeSock{}
	retp.clipipe = nil
	retp.accpipe = cli
	retp.listenpipe = nil
	retp.isaccept = false
	retp.isserver = true
	retp.data = []byte{}
	retp.readmills = p.readmills
	err = nil
	return
}

func (p *NpipeSock) WritePacket(ndata *NpipeData) (err error) {
	var conn *npipe.PipeConn
	var nowt time.Time
	var bs []byte
	var wlen int = 0
	var n int
	if p.isaccept {
		err = fmt.Errorf("not valid pipe")
		return
	}

	if p.isserver {
		conn = p.accpipe
	} else {
		conn = p.clipipe
	}

	if conn == nil {
		err = fmt.Errorf("not opened pipe")
		return
	}

	nowt = time.Now()
	err = conn.SetWriteDeadline(nowt.Add(time.Duration(p.readmills) * time.Millisecond))
	if err != nil {
		return
	}

	bs, err = ndata.Pack()
	if err != nil {
		return
	}

	wlen = 0
	for wlen < len(bs) {
		n, err = conn.Write(bs[wlen:])
		if err != nil {
			return
		}
		wlen += n
	}
	err = nil
	return
}

func (p *NpipeSock) ReadpacketTimeout() (ndata *NpipeData, err error) {
	var conn *npipe.PipeConn
	var nowt time.Time
	var dummy *NpipeData
	var n int
	var rlen int = 0
	var rbuf []byte
	var neterr net.Error
	var ok bool
	var curlen int
	var i int

	if p.isaccept {
		err = fmt.Errorf("accept not read")
		return
	}

	if p.isserver {
		conn = p.accpipe
	} else {
		conn = p.clipipe
	}

	if conn == nil {
		err = fmt.Errorf("read pipe not opened")
		return
	}

	nowt = time.Now()
	err = conn.SetReadDeadline(nowt.Add(time.Duration(p.readmills) * time.Millisecond))
	if err != nil {
		return
	}

	//fmt.Printf("data len[%d]\n", len(p.data))
	if uint32(len(p.data)) < NPIPE_PACK_HDR_SIZE {
		curlen = int(NPIPE_PACK_HDR_SIZE) - len(p.data)
		rbuf = make([]byte, curlen)
		rlen = 0
		for rlen < curlen {
			n, err = conn.Read(rbuf[:(curlen - rlen)])
			if err != nil {
				neterr, ok = err.(net.Error)
				if ok && neterr.Timeout() {
					err = nil
					ndata = nil
				}
				return
			}
			for i = 0; i < n; i++ {
				p.data = append(p.data, rbuf[i])
			}
			rlen += n
		}
	}

	/*now to get the data*/
	dummy = NewNpipeData()
	err = dummy.Unpack(p.data)
	if err != nil {
		return
	}

	if dummy.Length() <= len(p.data) {
		ndata = dummy
		err = nil
		p.data = []byte{}
		return
	}

	if dummy.Length() > 0x100000 {
		err = fmt.Errorf("overflow big [%d] size", dummy.Length())
		return
	}

	curlen = dummy.Length() - len(p.data)
	rbuf = make([]byte, curlen)
	rlen = 0
	for rlen < curlen {
		n, err = conn.Read(rbuf[:(curlen - rlen)])
		if err != nil {
			neterr, ok = err.(net.Error)
			if ok && neterr.Timeout() {
				err = nil
				ndata = nil
			}
			rbuf = []byte{}
			return
		}
		for i = 0; i < n; i++ {
			p.data = append(p.data, rbuf[i])
		}
		rlen += n
	}

	err = dummy.Unpack(p.data)
	if err != nil {
		return
	}

	ndata = dummy
	/*for clear the buffer*/
	p.data = []byte{}
	err = nil
	return
}

func (p *NpipeSock) Close() {
	if p.listenpipe != nil {
		p.listenpipe.Close()
	}
	p.listenpipe = nil

	if p.accpipe != nil {
		p.accpipe.Close()
	}
	p.accpipe = nil

	if p.clipipe != nil {
		p.clipipe.Close()
	}
	p.clipipe = nil

	p.isserver = false
	p.isaccept = false
	p.data = []byte{}
	p.readmills = 0
	return
}
