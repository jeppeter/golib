package jsonproto

import (
	"dbgutil"
)

type CommProto struct {
}

type BaseProto interface {
	HeadSize() int
	Totallen(inbuf []byte) (tlen int, err error)
	Pack(s string) (rbuf []byte, err error)
	UnPack(inbuf []byte) (s string, err error)
	PackData(inbuf []byte, code uint32) (outbuf []byte, err error)
	UnpackData(inbuf []byte) (outbuf []byte, code uint32, err error)
	UnpackType(inbuf []byte) (typename string, outbuf []byte, err error)
}

func NewCommProto() (pr *CommProto, err error) {
	pr = &CommProto{}
	err = nil
	return
}

func NewCommProtoFile() (pr *CommProto, err error) {
	return NewCommProto()
}

func (pr *CommProto) HeadSize() (ret int) {
	ret = 12
	return
}

func (pr *CommProto) pack_no_rsa(s string) (rbuf []byte, err error) {
	var cbuf []byte
	var tlen int
	var i int
	var curv int
	rbuf = make([]byte, 12)
	rbuf[0] = byte('J')
	rbuf[1] = byte('S')
	rbuf[2] = byte('O')
	rbuf[3] = byte('N')
	cbuf = []byte(s)
	/*end the end of buffer*/
	cbuf = append(cbuf, 0)
	for _, c := range cbuf {
		rbuf = append(rbuf, c)
	}
	tlen = len(rbuf)
	for i = 0; i < 8; i++ {
		curv = (tlen >> ((7 - i) * 4)) & 0xf
		if curv < 10 {
			rbuf[4+i] = byte(curv) + byte('0')
		} else {
			rbuf[4+i] = byte(curv-10) + byte('a')
		}
	}
	err = nil
	return
}

func (pr *CommProto) Totallen(inbuf []byte) (tlen int, err error) {
	var curlen int
	var i int
	if len(inbuf) < 12 {
		err = dbgutil.FormatError("len(%d) < 12", len(inbuf))
		return
	}

	if inbuf[0] != byte('J') || inbuf[1] != byte('S') ||
		inbuf[2] != byte('O') || inbuf[3] != byte('N') {
		err = dbgutil.FormatError("[:4] != JSON")
		return
	}
	tlen = 0
	for i = 0; i < 8; i++ {
		if inbuf[4+i] >= byte('0') && inbuf[4+i] <= byte('9') {
			curlen = int(inbuf[4+i] - byte('0'))
		} else if inbuf[4+i] >= byte('a') && inbuf[4+i] <= byte('f') {
			curlen = int(inbuf[4+i]-byte('a')) + 10
		} else if inbuf[4+i] >= byte('A') && inbuf[4+i] <= byte('F') {
			curlen = int(inbuf[4+i]-byte('A')) + 10
		} else {
			err = dbgutil.FormatError("[%d] not valid hex char", i+4)
			return
		}
		tlen += (curlen << (4 * (7 - i)))
	}
	err = nil
	return
}

func (pr *CommProto) Pack(s string) (rbuf []byte, err error) {
	rbuf, err = pr.pack_no_rsa(s)
	return
}

func (pr *CommProto) unpack_no_rsa(inbuf []byte) (s string, err error) {
	var tlen int
	var sbuf []byte
	var i int
	s = ""
	tlen, err = pr.Totallen(inbuf)
	if err != nil {
		return
	}

	if tlen > len(inbuf) {
		err = dbgutil.FormatError("total len [%d] > inbuf [%d]", tlen, len(inbuf))
		return
	}

	sbuf = []byte{}
	for i = 12; i < (tlen - 1); i++ {
		sbuf = append(sbuf, inbuf[i])
	}

	s = string(sbuf)
	err = nil
	return
}

func (pr *CommProto) UnPack(inbuf []byte) (s string, err error) {
	s, err = pr.unpack_no_rsa(inbuf)
	return
}

func (pr *CommProto) PackData(inbuf []byte, code uint32) (outbuf []byte, err error) {
	var outlen int = len(inbuf) + 12
	var idx int
	var curlen int
	outbuf = make([]byte, outlen)
	outbuf[0] = byte('D')
	outbuf[1] = byte('A')
	outbuf[2] = byte('T')
	outbuf[3] = byte('A')
	for idx = 0; idx < 4; idx++ {
		curlen = (outlen >> ((3 - idx) * 8)) & 0xff
		outbuf[idx+4] = byte(curlen)
	}

	for idx = 0; idx < 4; idx++ {
		curlen = int((code >> ((3 - idx) * 8)) & 0xff)
		outbuf[idx+8] = byte(curlen)
	}

	for idx = 0; idx < len(inbuf); idx++ {
		outbuf[12+idx] = inbuf[idx]
	}

	err = nil
	return
}

func (pr *CommProto) UnpackData(inbuf []byte) (outbuf []byte, code uint32, err error) {
	var tlen int
	var idx int
	var curlen int
	outbuf = make([]byte, 0)
	err = nil
	if len(inbuf) < 12 {
		err = dbgutil.FormatError("len(%d) < 12", len(inbuf))
		return
	}

	if inbuf[0] != byte('D') || inbuf[1] != byte('A') ||
		inbuf[2] != byte('T') || inbuf[3] != byte('A') {
		err = dbgutil.FormatError("%v not valid DATA", inbuf[:4])
		return
	}

	tlen = 0
	for idx = 0; idx < 4; idx++ {
		curlen = int(inbuf[idx+4]) & 0xff
		tlen += (curlen << ((3 - idx) * 8))
	}

	if tlen > len(inbuf) {
		err = dbgutil.FormatError("tlen [%d] > len(%d)", tlen, len(inbuf))
		return
	}
	code = uint32(0)
	for idx = 0; idx < 4; idx++ {
		curlen = int(inbuf[idx+8]) & 0xff
		code += uint32((curlen << ((3 - idx) & 8)))
	}

	if tlen > 12 {
		outbuf = make([]byte, (tlen - 12))
		for idx = 0; idx < (tlen - 12); idx++ {
			outbuf[idx] = inbuf[12+idx]
		}
	}
	err = nil
	return
}

func (pr *CommProto) UnpackType(inbuf []byte) (typename string, outbuf []byte, err error) {
	var s string
	var code uint32
	var cbuf []byte
	var idx int
	var ival int
	s, err = pr.UnPack(inbuf)
	if err == nil {
		typename = "JSON"
		outbuf = []byte(s)
		return
	}

	cbuf, code, err = pr.UnpackData(inbuf)
	if err == nil {
		typename = "DATA"
		outbuf = make([]byte, 4+len(cbuf))
		for idx = 0; idx < len(outbuf); idx++ {
			if idx < 4 {
				ival = int((code >> ((3 - idx) * 8)) & 0xff)
				outbuf[idx] = byte(ival)
			} else {
				outbuf[idx] = cbuf[(idx - 4)]
			}
		}
		return
	}

	return
}
