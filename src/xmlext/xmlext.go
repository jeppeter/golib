package xmlext

import (
	"dbgutil"
	"encoding/xml"
	"fmt"
)

type XmlExt struct {
	inner map[string]string
}

func ParseXmlExt(ins string) (retv *XmlExt, err error) {
	retv = &XmlExt{}
	retv.inner = make(map[string]string)
	err = xml.Unmarshal([]byte(ins), &retv.inner)
	if err != nil {
		retv = nil
		err = dbgutil.FormatError("can not parse %s\n%s", err.Error(), ins)
		return
	}
	return
}

func (retv *XmlExt) String() (s string) {
	var err error
	var outb []byte
	outb, err = xml.Marshal(retv.inner)
	if err != nil {
		var errs string
		errs = fmt.Sprintf("marshal %s", err.Error())
		panic(errs)
	}
	s = string(outb)
	return
}
