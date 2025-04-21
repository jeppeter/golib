package npipepack

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	NPIPE_PACK_HDR_SIZE = uint32(4)
)

type NpipeData struct {
	Size uint32
	Str  string
}

func NewNpipeData() (ptr *NpipeData) {
	ptr = &NpipeData{}
	ptr.Size = NPIPE_PACK_HDR_SIZE
	ptr.Str = ""
	return
}

func (p *NpipeData) String() string {
	return fmt.Sprintf("[%d][%s]", p.Size, p.Str)
}

func (p *NpipeData) Unpack(bs []byte) (err error) {
	var rb *bytes.Reader
	var nb []byte
	var i int
	if uint32(len(bs)) < NPIPE_PACK_HDR_SIZE {
		err = fmt.Errorf("[%d] < [%d] size", len(bs), NPIPE_PACK_HDR_SIZE)
		return
	}

	rb = bytes.NewReader(bs)
	err = binary.Read(rb, binary.LittleEndian, &p.Size)
	if err != nil {
		return
	}

	if p.Size <= uint32(len(bs)) {
		//fmt.Println(bs[NPIPE_PACK_HDR_SIZE:p.Size])
		nb = []byte{}
		for i = int(NPIPE_PACK_HDR_SIZE); i < int(p.Size); i++ {
			if bs[i] == byte(0) {
				break
			}
			nb = append(nb, bs[i])
		}

		p.Str = string(nb)
	}
	err = nil
	return
}

func (p *NpipeData) Pack() (bs []byte, err error) {
	var sb []byte
	var buf *bytes.Buffer
	var i int
	sb = []byte(p.Str)
	p.Size = uint32(len(sb)) + 1 + NPIPE_PACK_HDR_SIZE

	bs = []byte{}
	buf = new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, p.Size)
	if err != nil {
		return
	}

	bs = buf.Bytes()
	for i = 0; i < len(sb); i++ {
		bs = append(bs, sb[i])
	}

	bs = append(bs, byte(0))
	err = nil
	return
}

func (p *NpipeData) SetStr(s string) {
	p.Str = s
	return
}

func (p *NpipeData) Length() int {
	return int(p.Size)
}

func (p *NpipeData) GetStr() string {
	return p.Str
}
