package sockproto

import (
	"dbgutil"
	"jsonproto"
	"socktimeout"
)

type SockChannel struct {
	cli   *socktimeout.SockClient
	cbuf  []byte
	proto *jsonproto.CommProto
}

func NewSockChannel(cli *socktimeout.SockClient) (retp *SockChannel, err error) {
	retp = &SockChannel{}
	retp.cli = cli
	retp.cbuf = []byte{}
	retp.proto, err = jsonproto.NewCommProto()
	if err != nil {
		retp = nil
		return
	}
	err = nil
	return
}

func (retp *SockChannel) ReadJson(mills int) (rets string, err error) {
	var retb []byte
	var tlen int
	var readn int
	var i int
	rets = ""
	if len(retp.cbuf) < retp.proto.HeadSize() {
		readn = retp.proto.HeadSize() - len(retp.cbuf)
		retb, err = retp.cli.ReadTimeout(readn, mills)
		if err != nil {
			return
		}
		for i = 0; i < len(retb); i++ {
			retp.cbuf = append(retp.cbuf, retb[i])
		}
		if len(retp.cbuf) < retp.proto.HeadSize() {
			err = nil
			return
		}
	}

	/**/
	tlen, err = retp.proto.Totallen(retp.cbuf)
	if tlen > len(retp.cbuf) {
		readn = tlen - len(retp.cbuf)
		retb, err = retp.cli.ReadTimeout(readn, mills)
		if err != nil {
			return
		}
		for i = 0; i < len(retb); i++ {
			retp.cbuf = append(retp.cbuf, retb[i])
		}
	}

	err = nil
	if tlen == len(retp.cbuf) {
		rets, err = retp.proto.UnPack(retp.cbuf)
		if err == nil {
			retp.cbuf = []byte{}
		}
	}
	return

}

func (retp *SockChannel) WriteJson(jsons string, mills int) (err error) {
	var wbuf []byte
	var retn int
	wbuf, err = retp.proto.Pack(jsons)
	if err != nil {
		return
	}
	retn, err = retp.cli.WriteTimeout(wbuf, mills)
	if err != nil {
		return
	}

	if retn != len(wbuf) {
		err = dbgutil.FormatError("write %d in %d ", retn, mills)
		return
	}
	err = nil
	return
}

func (retp *SockChannel) Close() {
	retp.cli.Close()
	return
}
