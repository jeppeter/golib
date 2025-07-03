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
	curmap.started = true
	for {
		tks, err = xmldec.RawToken()
		if err != nil {
			if err != io.EOF {
				err = dbgutil.FormatError("parse error %s", err.Error())
				retv = nil
				return
			}

			if curmap != nil && curmap != retv.inner {
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
		outs += fmt.Sprintf("<%s", cmap.key)
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
	var curmap *xmlmap
	outs = ""
	for _, curmap = range retv.inner.chlds {
		outs += curmap.String()
	}
	return outs
}

func (retv *XmlExt) Ident(tabs int) (outs string) {
	var curmap *xmlmap
	outs = ""
	for _, curmap = range retv.inner.chlds {
		outs += curmap.Ident(tabs)
	}
	return outs
}

func (cmap *xmlmap) call_func(methname string, a ...interface{}) (retv []reflect.Value, err error) {
	var curval reflect.Value
	var methval reflect.Value
	var args []reflect.Value
	var idx int
	curval = reflect.ValueOf(cmap)
	methval = curval.MethodByName(methname)
	retv = []reflect.Value{}
	if !methval.IsValid() {
		err = dbgutil.FormatError("can not find [%s]", methname)
		return
	}
	args = []reflect.Value{}
	for idx = 0; idx < len(a); idx += 1 {
		args = append(args, reflect.ValueOf(a[idx]))
	}
	retv = methval.Call(args)
	err = nil
	return
}

func (cmap *xmlmap) find_function_callback(methname string, path string, a ...interface{}) (retv []reflect.Value, err error) {
	var sarr []string
	var k string
	var nk string
	var nextchld *xmlmap
	var ok bool
	var idx int
	retv = []reflect.Value{}

	if path == "" {
		return cmap.call_func(methname, a...)
	}
	sarr = strings.Split(path, "/")

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

		return nextchld.find_function_callback(methname, nk, a...)
	}
	return cmap.call_func(methname, a...)
}

func (cmap *xmlmap) call_func_set(methname string, a ...interface{}) (retv []reflect.Value, err error) {
	var curval reflect.Value
	var methval reflect.Value
	var args []reflect.Value
	var idx int
	curval = reflect.ValueOf(cmap)
	methval = curval.MethodByName(methname)
	retv = []reflect.Value{}
	if !methval.IsValid() {
		err = dbgutil.FormatError("can not find [%s]", methname)
		return
	}
	args = []reflect.Value{}
	for idx = 0; idx < len(a); idx += 1 {
		args = append(args, reflect.ValueOf(a[idx]))
	}
	retv = methval.Call(args)
	err = nil
	return
}

func (cmap *xmlmap) find_function_callback_set(methname string, path string, a ...interface{}) (retv []reflect.Value, err error) {
	var sarr []string
	var k string
	var nk string
	var nextchld *xmlmap
	var ok bool
	var idx int
	retv = []reflect.Value{}

	if path == "" {
		return cmap.call_func_set(methname, a...)
	}
	sarr = strings.Split(path, "/")

	for idx = 0; idx < len(sarr); idx += 1 {
		if len(sarr[idx]) == 0 {
			continue
		}
		k = sarr[idx]
		nextchld, ok = cmap.chlds[k]
		if !ok {
			/*not in the */
			nextchld = new_xmlmap()
			nextchld.started = true
			nextchld.key = k
			nextchld.val = ""
			cmap.chlds[k] = nextchld
		}
		if len(sarr) > idx {
			nk = strings.Join(sarr[idx+1:], "/")
		} else {
			nk = ""
		}

		return nextchld.find_function_callback_set(methname, nk, a...)
	}
	return cmap.call_func_set(methname, a...)
}

func (cmap *xmlmap) GetAttrs() (retv map[string]string) {
	return cmap.attrs
}

func get_first_path(path string) (curpath string, leftpath string) {
	var sarr []string
	var idx int
	sarr = strings.Split(path, "/")
	curpath = ""
	leftpath = ""
	for idx = 0; idx < len(sarr); idx += 1 {
		if len(sarr[idx]) == 0 {
			continue
		}
		curpath = sarr[idx]
		break
	}
	if idx < len(sarr) {
		leftpath = strings.Join(sarr[idx+1:], "/")
	}
	logutil.Debug("curpath [%s] leftpath[%s]", curpath, leftpath)
	return
}

func (ptr *XmlExt) get_root_xmlmap(path string) (curmap *xmlmap, leftpath string) {
	var ok bool
	var curpath string
	curmap = nil

	curpath, leftpath = get_first_path(path)
	if curpath == "" {
		if len(ptr.inner.chlds) == 1 {
			for _, curmap = range ptr.inner.chlds {
				break
			}
		}
	} else {
		curmap, ok = ptr.inner.chlds[curpath]
		if !ok {
			curmap = nil
		}
	}
	return
}

func (ptr *XmlExt) get_root_xmlmap_set(path string) (curmap *xmlmap, leftpath string) {
	var ok bool
	var curpath string

	curmap = nil

	curpath, leftpath = get_first_path(path)
	if curpath == "" {
		if len(ptr.inner.chlds) == 1 {
			for _, curmap = range ptr.inner.chlds {
				break
			}
		}
	} else {
		curmap, ok = ptr.inner.chlds[curpath]
		if !ok {
			curmap = new_xmlmap()
			curmap.started = true
			curmap.key = curpath
			ptr.inner.chlds[curpath] = curmap
		}
	}
	return
}

func (ptr *XmlExt) GetAttrs(path string) (retv map[string]string, err error) {
	var cretv []reflect.Value
	var curmap *xmlmap = nil
	var leftpath string
	curmap, leftpath = ptr.get_root_xmlmap(path)

	if curmap == nil {
		err = dbgutil.FormatError("can not get path [%s]", path)
		return
	}

	retv = make(map[string]string)
	cretv, err = curmap.find_function_callback("GetAttrs", leftpath)
	if err != nil {
		return
	}

	if len(cretv) != 1 || cretv[0].Kind() != reflect.Map {
		err = dbgutil.FormatError("return value not map")
		return
	}

	i := cretv[0].Interface()
	retv = i.(map[string]string)
	err = nil
	return
}

func (retv *XmlExt) GetAttrsMust(path string) (retn map[string]string) {
	var err error
	retn, err = retv.GetAttrs(path)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (cmap *xmlmap) GetAttrValue(k string) (retn string, err error) {
	var ok bool
	retn, ok = cmap.attrs[k]
	if !ok {
		err = dbgutil.FormatError("no [%s] attr", k)
		return
	}
	err = nil
	return
}

func (retv *XmlExt) GetAttrValue(path, k string) (retn string, err error) {
	var cretv []reflect.Value
	var curmap *xmlmap = nil
	var leftpath string

	curmap, leftpath = retv.get_root_xmlmap(path)
	if curmap == nil {
		err = dbgutil.FormatError("can not get path [%s]", path)
		return
	}

	var ival interface{}
	retn = ""
	cretv, err = curmap.find_function_callback("GetAttrValue", leftpath, k)
	if err != nil {
		return
	}

	if len(cretv) != 2 {
		err = dbgutil.FormatError("GetAttrValue len(%d)", len(cretv))
		return
	}
	if cretv[0].Kind() != reflect.String {
		err = dbgutil.FormatError("GetAttrValue [0] type != String")
		return
	}

	ival = cretv[1].Interface()
	if ival != nil {
		switch ival.(type) {
		case error:
			err = ival.(error)
		default:
			err = dbgutil.FormatError("GetAttrValue [1] type != Error")
		}
		return
	}

	if ival != nil {
		err = ival.(error)
		return
	}
	retn = cretv[0].String()
	err = nil
	return
}

func (retv *XmlExt) GetAttrValueMust(path, k string) (retn string) {
	var err error
	retn, err = retv.GetAttrValue(path, k)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (cmap *xmlmap) GetValue() string {
	return cmap.val
}

func (retv *XmlExt) GetValue(path string) (retn string, err error) {
	var cretv []reflect.Value
	var curmap *xmlmap = nil
	var leftpath string
	curmap, leftpath = retv.get_root_xmlmap(path)
	if curmap == nil {
		err = dbgutil.FormatError("can not get [%s]", path)
		return
	}
	retn = ""
	cretv, err = curmap.find_function_callback("GetValue", leftpath)
	if err != nil {
		return
	}

	if len(cretv) != 1 || cretv[0].Kind() != reflect.String {
		err = dbgutil.FormatError("GetValue reflect.Value not string")
		return
	}
	retn = cretv[0].String()
	err = nil
	return
}

func (retv *XmlExt) GetValueMust(path string) (retn string) {
	var err error
	retn, err = retv.GetValue(path)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (cmap *xmlmap) GetChilds() (retv []string) {
	var k string
	retv = []string{}
	for k, _ = range cmap.chlds {
		retv = append(retv, k)
	}
	return
}

func (retv *XmlExt) GetChilds(path string) (retn []string, err error) {
	var cretv []reflect.Value
	var curmap *xmlmap = nil
	var leftpath string
	curmap, leftpath = retv.get_root_xmlmap(path)
	if curmap == nil {
		err = dbgutil.FormatError("can not get [%s]", path)
		return
	}
	retn = []string{}
	cretv, err = curmap.find_function_callback("GetChilds", leftpath)
	if err != nil {
		return
	}
	if len(cretv) != 0 || (cretv[0].Kind() != reflect.Array && cretv[0].Kind() != reflect.Slice) {
		err = dbgutil.FormatError("not return array")
		return
	}

	v := cretv[0].Interface()
	retn = v.([]string)
	err = nil
	return
}

func (retv *XmlExt) GetChildsMust(path string) (retn []string) {
	var err error
	retn, err = retv.GetChilds(path)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (cmap *xmlmap) SetAttr(k, v string) (retn string) {
	var ok bool
	retn, ok = cmap.attrs[k]
	if !ok {
		retn = ""
	}
	cmap.attrs[k] = v
	return
}

func (retv *XmlExt) SetAttr(path, k, v string) (retn string, err error) {
	var cretv []reflect.Value
	var curmap *xmlmap = nil
	var leftpath string
	curmap, leftpath = retv.get_root_xmlmap_set(path)
	if curmap == nil {
		err = dbgutil.FormatError("can not get set [%s]", path)
		return
	}

	retn = ""
	cretv, err = curmap.find_function_callback_set("SetAttr", leftpath, k, v)
	if err != nil {
		return
	}
	if len(cretv) != 1 || cretv[0].Kind() != reflect.String {
		err = dbgutil.FormatError("SetAttr return not Stirng")
		return
	}
	retn = cretv[0].String()
	err = nil
	return
}

func (retv *XmlExt) SetAttrMust(path, k, v string) (retn string) {
	var err error
	retn, err = retv.SetAttr(path, k, v)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (cmap *xmlmap) SetValue(v string) (retn string) {
	retn = cmap.val
	cmap.val = v
	return
}

func (retv *XmlExt) SetValue(path, v string) (retn string, err error) {
	var cretv []reflect.Value
	var curmap *xmlmap = nil
	var leftpath string
	curmap, leftpath = retv.get_root_xmlmap_set(path)
	if curmap == nil {
		err = dbgutil.FormatError("can not get set [%s]", path)
		return
	}
	retn = ""
	cretv, err = curmap.find_function_callback_set("SetValue", leftpath, v)
	if err != nil {
		return
	}
	if len(cretv) != 1 || cretv[0].Kind() != reflect.String {
		err = dbgutil.FormatError("SetValue return not Stirng")
		return
	}
	retn = cretv[0].String()
	err = nil
	return
}

func (retv *XmlExt) SetValueMust(path, v string) (retn string) {
	var err error
	retn, err = retv.SetValue(path, v)
	if err != nil {
		panic(err.Error())
	}
	return
}
