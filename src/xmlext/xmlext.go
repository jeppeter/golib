package xmlext

import (
	"bytes"
	"dbgutil"
	"encoding/xml"
	"fmt"
	"io"
	"logutil"
	"reflect"
	"strings"
)

type xmlmap struct {
	attrs   map[string]string
	chlds   map[string]*xmlmap
	val     string
	key     string
	started bool
}

func new_xmlmap() (retv *xmlmap) {
	retv = &xmlmap{}
	retv.attrs = make(map[string]string)
	retv.chlds = make(map[string]*xmlmap)
	retv.val = ""
	retv.key = ""
	retv.started = false
	return
}

type XmlExt struct {
	inner *xmlmap
}

func xml_fmt_name(n xml.Name) (rets string) {
	if len(n.Space) > 0 {
		rets = fmt.Sprintf("%s:%s", n.Space, n.Local)
	} else {
		rets = fmt.Sprintf("%s", n.Local)
	}
	return
}

func ParseXmlExt(ins string) (retv *XmlExt, err error) {
	retv = &XmlExt{}
	retv.inner = new_xmlmap()
	xmldec := xml.NewDecoder(bytes.NewReader([]byte(ins)))
	var tks xml.Token
	var parentmaps []*xmlmap = []*xmlmap{}
	var curmap *xmlmap = nil
	var idx int
	var key string
	var val string
	var sizev int
	var started bool
	var insertbytes []byte
	var insertlen int
	curmap = retv.inner
	for {
		tks, err = xmldec.RawToken()
		if err != nil {
			if err != io.EOF {
				err = dbgutil.FormatError("parse error %s", err.Error())
				retv = nil
				return
			}

			if curmap != nil {
				err = dbgutil.FormatError("can not parse [%s]", curmap.key)
				return
			}
			break
		}
		switch v := tks.(type) {
		case xml.CharData:
			logutil.Debug("CharData %v", v)
			if curmap != nil && curmap.started {
				/*now to check for the chardata*/
				insertbytes = []byte{}
				started = false
				for idx = 0; idx < len(v); idx += 1 {
					if !started {
						if v[idx] != byte('\n') && v[idx] != byte(' ') && v[idx] != byte('\t') && v[idx] != byte('\r') {
							insertbytes = append(insertbytes, v[idx])
							started = true
						}
					} else {
						insertbytes = append(insertbytes, v[idx])
					}
				}
				if len(insertbytes) > 0 {
					insertlen = len(insertbytes)
					for idx = len(insertbytes) - 1; idx >= 0; idx -= 1 {
						if insertbytes[idx] == byte('\n') || insertbytes[idx] == byte(' ') || insertbytes[idx] == byte('\t') && insertbytes[idx] == byte('\r') {
							insertlen -= 1
						} else {
							break
						}
					}
					if insertlen > 0 {
						val = string(insertbytes[:insertlen])
						logutil.Debug("add val [%s]", val)
						curmap.val += val
					}

				}
			}
		case xml.Comment:
			logutil.Debug("Comment %v", v)
		case xml.Directive:
			logutil.Debug("Directive %v", v)
		case xml.ProcInst:
			logutil.Debug("ProcInst %v", v)
		case xml.StartElement:
			key = xml_fmt_name(v.Name)
			logutil.Debug("StartElement name [%s]", key)
			if curmap.started {
				parentmaps = append(parentmaps, curmap)
				curmap = new_xmlmap()
			}
			if !curmap.started {
				curmap.started = true
				curmap.key = key

				for _, attr := range v.Attr {
					logutil.Debug("[%s:%s] -- [%s]", attr.Name.Space, attr.Name.Local, attr.Value)
					key = xml_fmt_name(attr.Name)
					val = fmt.Sprintf("%s", attr.Value)
					curmap.attrs[key] = val
				}

			}
		case xml.EndElement:
			if curmap == nil {
				err = dbgutil.FormatError("no current context")
				return
			}

			key = xml_fmt_name(v.Name)
			if key != curmap.key {
				err = dbgutil.FormatError("EndElement [%s] not match [%s]", key, curmap.key)
				return
			}
			logutil.Debug("EndElement [%s]", key)

			sizev = len(parentmaps)
			if sizev > 0 {
				parentmaps[sizev-1].chlds[key] = curmap
				curmap = nil
				curmap = parentmaps[sizev-1]
				parentmaps = parentmaps[:(sizev - 1)]
			} else {
				curmap = nil
			}
		default:
			logutil.Debug("type %s %v", reflect.TypeOf(v), v)
		}
	}

	err = nil
	return
}

func (cmap *xmlmap) String() (outs string) {
	var k, v string
	outs = ""
	if len(cmap.val) > 0 || len(cmap.chlds) > 0 {
		outs += fmt.Sprintf("<%s", cmap.key)
		if len(cmap.attrs) > 0 {
			for k, v = range cmap.attrs {
				logutil.Info("attr [%s]=[%s]", k, v)
				outs += fmt.Sprintf(" ")
				outs += fmt.Sprintf("%s=\"%s\"", k, v)
			}
		}
		outs += ">"
		if len(cmap.val) > 0 {
			outs += fmt.Sprintf("%s", cmap.val)
		}
		for _, chld := range cmap.chlds {
			outs += chld.String()
		}
		outs += fmt.Sprintf("</%s>", cmap.key)
	} else {
		outs += fmt.Sprintf("< %s", cmap.key)
		if len(cmap.attrs) > 0 {
			for k, v = range cmap.attrs {
				outs += fmt.Sprintf(" ")
				outs += fmt.Sprintf("%s=\"%s\"", k, v)
			}
		}
		outs += "/>"
	}
	return
}

func format_tab_xml(tab int) (s string) {
	s = ""
	for i := 0; i < tab; i += 1 {
		s += fmt.Sprintf("    ")
	}
	return
}

func (cmap *xmlmap) Ident(tabs int) (outs string) {
	var k, v string
	outs = ""
	if len(cmap.val) > 0 || len(cmap.chlds) > 0 {
		outs += format_tab_xml(tabs)
		outs += fmt.Sprintf("<%s", cmap.key)
		if len(cmap.attrs) > 0 {
			for k, v = range cmap.attrs {
				logutil.Info("attr [%s]=[%s]", k, v)
				outs += fmt.Sprintf(" ")
				outs += fmt.Sprintf("%s=\"%s\"", k, v)
			}
		}
		outs += ">\n"
		if len(cmap.val) > 0 {
			outs += format_tab_xml(tabs + 1)
			outs += fmt.Sprintf("%s\n", cmap.val)
		}
		for _, chld := range cmap.chlds {
			outs += chld.Ident(tabs + 1)
		}
		outs += format_tab_xml(tabs)
		outs += fmt.Sprintf("</%s>\n", cmap.key)
	} else {
		outs += format_tab_xml(tabs)
		outs += fmt.Sprintf("<%s", cmap.key)
		if len(cmap.attrs) > 0 {
			for k, v = range cmap.attrs {
				outs += fmt.Sprintf(" ")
				outs += fmt.Sprintf("%s=\"%s\"", k, v)
			}
		}
		outs += "/>\n"
	}
	return
}

func (retv *XmlExt) String() (outs string) {
	return retv.inner.String()
}

func (retv *XmlExt) Ident(tabs int) (outs string) {
	return retv.inner.Ident(tabs)
}

func (cmap *xmlmap) getattrs(path string) (retv map[string]string, err error) {
	var sarr []string
	var k string
	var nk string
	var nextchld *xmlmap
	var ok bool
	var idx int

	if path == "" {
		retv = cmap.attrs
		err = nil
		return
	}
	sarr = strings.Split(path, "/")
	if len(sarr) < 1 {
		return
	}

	for idx = 0; idx < len(sarr); idx += 1 {
		if len(sarr[idx]) == 0 {
			continue
		}
		k = sarr[idx]
		nextchld, ok = cmap.chlds[k]
		if !ok {
			err = dbgutil.FormatError("can not get [%s]", k)
			return
		}
		if len(sarr) > idx {
			nk = strings.Join(sarr[idx+1:], "/")
		} else {
			nk = ""
		}

		return nextchld.getattrs(nk)
	}
	retv = cmap.attrs
	err = nil
	return

}

func (retv *XmlExt) GetAttrs(path string) (map[string]string, error) {
	return retv.inner.getattrs(path)
}

func (retv *XmlExt) GetAttrsMust(path string) (retn map[string]string) {
	var err error
	retn, err = retv.GetAttrs(path)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (cmap *xmlmap) getvalue(path string) (retv string, err error) {
	var sarr []string
	var k string
	var nk string
	var nextchld *xmlmap
	var ok bool
	var idx int

	if path == "" {
		retv = cmap.val
		err = nil
		return
	}
	sarr = strings.Split(path, "/")
	if len(sarr) < 1 {
		return
	}

	for idx = 0; idx < len(sarr); idx += 1 {
		if len(sarr[idx]) == 0 {
			continue
		}
		k = sarr[idx]
		nextchld, ok = cmap.chlds[k]
		if !ok {
			err = dbgutil.FormatError("can not get [%s]", k)
			return
		}
		if len(sarr) > idx {
			nk = strings.Join(sarr[idx+1:], "/")
		} else {
			nk = ""
		}

		return nextchld.getvalue(nk)
	}
	retv = cmap.val
	err = nil
	return
}

func (retv *XmlExt) GetValue(path string) (retn string, err error) {
	retn, err = retv.inner.getvalue(path)
	return
}

func (retv *XmlExt) GetValueMust(path string) (retn string) {
	var err error
	retn, err = retv.inner.getvalue(path)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (cmap *xmlmap) _inner_chld_keys() []string {
	var retv []string = []string{}
	for k, _ := range cmap.chlds {
		retv = append(retv, k)
	}
	return retv
}

func (cmap *xmlmap) getchilds(path string) (retv []string, err error) {
	var sarr []string
	var k string
	var nk string
	var nextchld *xmlmap
	var ok bool
	var idx int

	if path == "" {
		retv = cmap._inner_chld_keys()
		err = nil
		return
	}
	sarr = strings.Split(path, "/")
	if len(sarr) < 1 {
		return
	}

	for idx = 0; idx < len(sarr); idx += 1 {
		if len(sarr[idx]) == 0 {
			continue
		}
		k = sarr[idx]
		nextchld, ok = cmap.chlds[k]
		if !ok {
			err = dbgutil.FormatError("can not get [%s]", k)
			return
		}
		if len(sarr) > idx {
			nk = strings.Join(sarr[idx+1:], "/")
		} else {
			nk = ""
		}

		return nextchld.getchilds(nk)
	}
	retv = cmap._inner_chld_keys()
	err = nil
	return
}

func (retv *XmlExt) GetChilds(path string) (retn []string, err error) {
	retn, err = retv.inner.getchilds(path)
	return
}

func (retv *XmlExt) GetChildsMust(path string) (retn []string) {
	var err error
	retn, err = retv.inner.getchilds(path)
	if err != nil {
		panic(err.Error())
	}
	return
}
