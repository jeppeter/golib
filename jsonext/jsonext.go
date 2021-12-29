package jsonext

import (
	"encoding/json"
	"fmt"
	//"log"
	"math"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const (
	EPSILON_VALUE = float64(0.000001)
)

func parseMessage(msg string) (map[string]interface{}, error) {
	var v map[string]interface{}
	err := json.Unmarshal([]byte(msg), &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func GetJsonMap(valstr string) (retv map[string]interface{}, err error) {
	retv, err = SafeParseMessage(valstr)
	if err != nil {
		return
	}
	err = nil
	return
}

func GetJsonArray(valstr string) (retv []interface{}, err error) {
	dec := json.NewDecoder(strings.NewReader(valstr))
	err = dec.Decode(&retv)
	if err != nil {
		return
	}
	err = nil
	return
}

func __GetJsonArrayItem(sarr []interface{}, idx int) (val interface{}, types string, err error) {
	if len(sarr) <= idx {
		err = fmt.Errorf("[%d] out of range", idx)
		return
	}
	val = sarr[idx]
	if val == nil {
		types = "null"
	} else {
		switch val.(type) {
		case string:
			types = "string"
		case float64:
			types = "float64"
		case []interface{}:
			types = "array"
		case map[string]interface{}:
			types = "map"
		case bool:
			types = "bool"
		default:
			err = fmt.Errorf("[%d] type unsupported [%s]", idx, reflect.TypeOf(val).String())
			return
		}
	}
	err = nil
	return
}

func GetJsonArrayItemString(sarr []interface{}, idx int) (val string, err error) {
	var vinter interface{}
	var types string
	vinter, types, err = __GetJsonArrayItem(sarr, idx)
	if err != nil {
		return
	}
	if types != "string" {
		err = fmt.Errorf("[%d] item not string [%s]", idx, types)
		return
	}
	val = vinter.(string)
	err = nil
	return
}

func GetJsonArrayItemInt(sarr []interface{}, idx int) (val int, err error) {
	var vinter interface{}
	var types string
	var varf float64
	vinter, types, err = __GetJsonArrayItem(sarr, idx)
	if err != nil {
		return
	}
	if types != "float64" {
		err = fmt.Errorf("[%d] item not string [%s]", idx, types)
		return
	}
	varf = vinter.(float64)
	val = int(varf)
	err = nil
	return
}

func GetJsonArrayItemFloat(sarr []interface{}, idx int) (val float64, err error) {
	var vinter interface{}
	var types string
	vinter, types, err = __GetJsonArrayItem(sarr, idx)
	if err != nil {
		return
	}
	if types != "float64" {
		err = fmt.Errorf("[%d] item not string [%s]", idx, types)
		return
	}
	val = vinter.(float64)
	err = nil
	return
}

func GetJsonArrayItemBool(sarr []interface{}, idx int) (val bool, err error) {
	var vinter interface{}
	var types string
	vinter, types, err = __GetJsonArrayItem(sarr, idx)
	if err != nil {
		return
	}
	if types != "bool" {
		err = fmt.Errorf("[%d] item not float [%s]", idx, types)
		return
	}
	val = vinter.(bool)
	err = nil
	return
}

func GetJsonArrayItemNull(sarr []interface{}, idx int) (err error) {
	var types string
	_, types, err = __GetJsonArrayItem(sarr, idx)
	if err != nil {
		return
	}
	if types != "null" {
		err = fmt.Errorf("[%d] item not null [%s]", idx, types)
		return
	}
	err = nil
	return
}

func GetJsonArrayItemMap(sarr []interface{}, idx int) (val map[string]interface{}, err error) {
	var vinter interface{}
	var types string
	vinter, types, err = __GetJsonArrayItem(sarr, idx)
	if err != nil {
		return
	}
	if types != "map" {
		err = fmt.Errorf("[%d] item not string [%s]", idx, types)
		return
	}
	val = vinter.(map[string]interface{})
	err = nil
	return
}

func GetJsonArrayItemArray(sarr []interface{}, idx int) (val []interface{}, err error) {
	var vinter interface{}
	var types string
	vinter, types, err = __GetJsonArrayItem(sarr, idx)
	if err != nil {
		return
	}
	if types != "array" {
		err = fmt.Errorf("[%d] item not string [%s]", idx, types)
		return
	}
	val = vinter.([]interface{})
	err = nil
	return
}

func SafeParseMessage(fmsg string) (map[string]interface{}, error) {
	v, err := parseMessage(fmsg)
	if err != nil {
		pmsg := `"` + fmsg + `"`
		//pmsg := fmsg
		cmsg, err := strconv.Unquote(pmsg)
		if err != nil {
			cmsg = fmsg
		}
		v, err = parseMessage(cmsg)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

func __FormatLevel(level int) string {
	s := ""
	for i := 0; i < level; i++ {
		s += fmt.Sprintf("  ")
	}
	return s
}

type FormatClass interface {
	Format(level int, keyname string, value interface{}) (string, error)
	SupportType() string
}

var (
	suppFormatClass []FormatClass
)

type Float64FormatClass struct {
	FormatClass
}

func (f Float64FormatClass) Format(level int, keyname string, value interface{}) (string, error) {
	s := ""
	ival := int(value.(float64))
	if math.Abs((float64(ival) - (value.(float64)))) < EPSILON_VALUE {
		s += fmt.Sprintf(" %d", ival)
	} else {
		s += fmt.Sprintf(" %f", value.(float64))
	}
	return s, nil
}

func (f Float64FormatClass) SupportType() string {
	return "float64"
}

type IntFormatClass struct {
	FormatClass
}

func (f IntFormatClass) Format(level int, keyname string, value interface{}) (string, error) {
	ival := int(value.(int))
	s := fmt.Sprintf(" %d", ival)
	return s, nil
}

func (f IntFormatClass) SupportType() string {
	return "int"
}

type Int32FormatClass struct {
	FormatClass
}

func (f Int32FormatClass) Format(level int, keyname string, value interface{}) (string, error) {
	ival := int32(value.(int32))
	s := fmt.Sprintf(" %d", ival)
	return s, nil
}

func (f Int32FormatClass) SupportType() string {
	return "int32"
}

type Int64FormatClass struct {
	FormatClass
}

func (f Int64FormatClass) Format(level int, keyname string, value interface{}) (string, error) {
	ival := int64(value.(int64))
	s := fmt.Sprintf(" %d", ival)
	return s, nil
}

func (f Int64FormatClass) SupportType() string {
	return "int64"
}

type Float32FormatClass struct {
	FormatClass
}

func (f Float32FormatClass) Format(level int, keyname string, value interface{}) (string, error) {
	s := ""
	ival := int(value.(float32))
	if math.Abs((float64(ival) - float64((value.(float32))))) < EPSILON_VALUE {
		s += fmt.Sprintf(" %d", ival)
	} else {
		s += fmt.Sprintf(" %f", value.(float32))
	}
	return s, nil
}

func (f Float32FormatClass) SupportType() string {
	return "float32"
}

type ArrayFormatClass struct {
	FormatClass
}

func (f ArrayFormatClass) Format(level int, keyname string, value interface{}) (string, error) {
	s := ""
	s += "["
	for i, aa := range value.([]interface{}) {
		if i != 0 {
			s += ","
		}
		/*we do not format any keyname*/
		arrays, err := __FormatValue(level+1, "", aa)
		if err != nil {
			return "", err
		}
		s += arrays
	}
	s += "]"
	return s, nil
}

func (f ArrayFormatClass) SupportType() string {
	return "[]interface {}"
}

type MapStringFormatClass struct {
	FormatClass
}

func (f MapStringFormatClass) Format(level int, keyname string, v interface{}) (string, error) {
	var curmap map[string]interface{}
	var sortkeys []string
	s := ""

	curmap = v.(map[string]interface{})
	if len(keyname) != 0 {
		s += __FormatName(level, keyname)
	}

	s += "{"

	for k := range curmap {
		sortkeys = append(sortkeys, k)
	}
	sort.Strings(sortkeys)

	for i, kk := range sortkeys {
		kv, ok := curmap[kk]
		if !ok {
			err := fmt.Errorf("can not find key (%s)", kk)
			return "", err
		}
		if i != 0 {
			s += ",\n"
		} else {
			s += "\n"
		}
		ks, err := __FormatJsonValue(level+1, kk, kv)
		if err != nil {
			return "", err
		}
		s += ks
	}
	s += "\n"
	s += __FormatLevel(level)
	s += "}"
	return s, nil
}

func (f MapStringFormatClass) SupportType() string {
	return "map[string]interface {}"
}

type StringFormatClass struct {
	FormatClass
}

func (f StringFormatClass) Format(level int, keyname string, value interface{}) (string, error) {
	s := ""
	s += fmt.Sprintf(" %s", strconv.Quote(value.(string)))
	return s, nil
}

func (f StringFormatClass) SupportType() string {
	return "string"
}

type BoolFormatClass struct {
	FormatClass
}

func (f BoolFormatClass) Format(level int, keyname string, value interface{}) (string, error) {
	s := ""
	if value.(bool) {
		s += "true"
	} else {
		s += "false"
	}
	return s, nil
}

func (f BoolFormatClass) SupportType() string {
	return "bool"
}

var (
	formatMap map[string]FormatClass
)

func init() {
	formatMap = make(map[string]FormatClass)
	arrcls := ArrayFormatClass{}
	formatMap[arrcls.SupportType()] = arrcls
	fl32cls := Float32FormatClass{}
	formatMap[fl32cls.SupportType()] = fl32cls
	fl64cls := Float64FormatClass{}
	formatMap[fl64cls.SupportType()] = fl64cls
	scls := StringFormatClass{}
	formatMap[scls.SupportType()] = scls
	boolcls := BoolFormatClass{}
	formatMap[boolcls.SupportType()] = boolcls
	mapstrcls := MapStringFormatClass{}
	formatMap[mapstrcls.SupportType()] = mapstrcls
	intcls := IntFormatClass{}
	formatMap[intcls.SupportType()] = intcls
	int32cls := Int32FormatClass{}
	formatMap[int32cls.SupportType()] = int32cls
	int64cls := Int64FormatClass{}
	formatMap[int64cls.SupportType()] = int64cls
}

func __FormatName(level int, keyname string) string {
	s := ""
	s += __FormatLevel(level)
	s += fmt.Sprintf("\"%s\" : ", keyname)
	return s
}

func __FormatValue(level int, keyname string, value interface{}) (string, error) {
	var err error
	var typestr string
	s := ""
	typestr = reflect.TypeOf(value).String()
	fcls, ok := formatMap[typestr]
	if !ok {
		err := fmt.Errorf("(%s) support type %s", keyname, typestr)
		return "", err
	}

	ss, err := fcls.Format(level, keyname, value)
	if err != nil {
		return "", err
	}

	s += ss
	return s, nil
}

func __FormatValueBasic(level int, keyname string, value interface{}) (string, error) {
	var s string
	var err error
	//Debug("[%d].(%s) = %q", level, keyname, value)
	s = ""
	s += __FormatName(level+1, keyname)
	sets, err := __FormatValue(level, keyname, value)
	if err != nil {
		return "", err
	}
	s += sets

	return s, nil
}

func __FormatJsonValue(level int, keyname string, value interface{}) (string, error) {
	var err error
	var s, ss string
	s = ""
	//Debug("[%d].(%s)type(%s) = %q", level, keyname, reflect.TypeOf(value).String(), value)

	switch value.(type) {
	case map[string]interface{}:
		ss, err = FormatJsonValue(level+1, keyname, value.(map[string]interface{}))
		if err != nil {
			return "", err
		}
		s += ss
	default:
		ss, err = __FormatValueBasic(level, keyname, value)
		if err != nil {
			return "", err
		}
		s += ss
	}
	//Debug("end[%d].(%s)type(%s) = %q", level, keyname, reflect.TypeOf(value).String(), value)
	return s, nil
}

func FormatJsonValue(level int, keyname string, v map[string]interface{}) (string, error) {
	var curmap map[string]interface{}
	var sortkeys []string
	s := ""

	curmap = v
	if len(keyname) != 0 {
		s += __FormatName(level, keyname)
	}

	s += "{"

	for k := range curmap {
		sortkeys = append(sortkeys, k)
	}
	sort.Strings(sortkeys)

	for i, kk := range sortkeys {
		kv, ok := curmap[kk]
		if !ok {
			err := fmt.Errorf("can not find key (%s)", kk)
			return "", err
		}
		if i != 0 {
			s += ",\n"
		} else {
			s += "\n"
		}
		ks, err := __FormatJsonValue(level, kk, kv)
		if err != nil {
			return "", err
		}
		s += ks
	}
	s += "\n"
	s += __FormatLevel(level)
	s += "}"
	return s, nil
}

func FormatJsonArray(level int, keyname string, v []interface{}) (string, error) {
	var s string
	var err error
	s, err = __FormatJsonValue(level, keyname, v)
	if err != nil {
		return "", err
	}
	return s, nil
}

func writeToFile(infile string, jsonbytes []byte) error {
	fw, err := os.OpenFile(infile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer fw.Close()
	totalw := 0
	for totalw < len(jsonbytes) {
		n, err := fw.Write(jsonbytes[totalw:])
		if err != nil {
			return err
		}
		totalw += n
	}
	return nil
}

func WriteJson(infile string, v map[string]interface{}) error {
	jsonstring, err := FormatJsonValue(0, "", v)
	if err != nil {
		return err
	}
	jsonbytes := []byte(jsonstring)
	return writeToFile(infile, jsonbytes)
}

func WriteJsonString(infile string, val string) error {
	valmap, err := SafeParseMessage(val)
	if err != nil {
		return err
	}

	return WriteJson(infile, valmap)
}

func SetJsonValue(path, typestr, value string, v map[string]interface{}) (map[string]interface{}, error) {
	var pathext []string
	var tmpext []string
	var err error
	var mapv map[string]interface{}
	var arrayv []interface{}
	var fval float64
	var ival int

	tmpext = strings.Split(path, "/")
	for _, a := range tmpext {
		if len(a) > 0 {
			/*we fil the path*/
			pathext = append(pathext, a)
		}
	}

	curmap := v
	if len(pathext) == 0 {
		if typestr != "map" {
			err = fmt.Errorf("invalid path(%s) and type(%s)", path, typestr)
			return nil, err
		}
		mapv, err = GetJsonMap(value)
		if err != nil {
			return nil, err
		}
		v = mapv
		return v, nil
	} else {
		for i, curpath := range pathext {
			if i == (len(pathext) - 1) {
				/*this is the last one ,so we set the value*/
				switch typestr {
				case "string":
					curmap[curpath] = value
				case "float64":
					fval, err = strconv.ParseFloat(value, 64)
					if err != nil {
						return nil, err
					}
					curmap[curpath] = fval
				case "int":
					ival, err = strconv.Atoi(value)
					if err != nil {
						return nil, err
					}
					curmap[curpath] = ival
				case "map":
					mapv, err = GetJsonMap(value)
					if err != nil {
						return nil, err
					}
					curmap[curpath] = mapv
				case "array":
					arrayv, err = GetJsonArray(value)
					if err != nil {
						return nil, err
					}
					curmap[curpath] = arrayv
				default:
					err = fmt.Errorf("unknown type %s", typestr)
					return nil, err
				}
				return v, nil
			}
			curval, ok := curmap[curpath]
			if !ok {
				/*we make the next use*/
				curmap[curpath] = make(map[string]interface{})
				curmap = curmap[curpath].(map[string]interface{})
			} else {
				switch curval.(type) {
				case map[string]interface{}:
					curmap = curval.(map[string]interface{})
				default:
					/*we make the map string*/
					curval = make(map[string]interface{})
					curmap[curpath] = curval.(map[string]interface{})
					curmap = curval.(map[string]interface{})
				}
			}
		}
	}

	err = fmt.Errorf("unknown path %s", path)
	return nil, err
}

func DeleteJsonValue(path string, v map[string]interface{}, force int) (map[string]interface{}, error) {
	var pathext []string
	var tmpext []string
	var err error

	tmpext = strings.Split(path, "/")
	for _, a := range tmpext {
		if len(a) > 0 {
			pathext = append(pathext, a)
		}
	}

	curmap := v
	if len(pathext) > 0 {
		for i, curpath := range pathext {
			curval, ok := curmap[curpath]
			if !ok {
				if force > 0 {
					return v, nil
				}
				err = fmt.Errorf("can not find (%s) value", path)
				return nil, err
			}
			if i == (len(pathext) - 1) {
				/*this is the last one ,so we set the value*/
				delete(curmap, curpath)
				return v, nil
			}
			switch curval.(type) {
			case map[string]interface{}:
				curmap = curval.(map[string]interface{})
			default:
				if force > 0 {
					delete(curmap, curpath)
					return v, nil
				}
				err = fmt.Errorf("can not handle path %s", curpath)
				return nil, err
			}
		}
	} else {
		/*we set the null for total delete*/
		return nil, nil
	}

	if force > 0 {
		return v, nil
	}

	err = fmt.Errorf("unknown path %s", path)
	return nil, err

}

func GetJson(infile string) (map[string]interface{}, error) {
	var v map[string]interface{}
	fp, err := os.Open(infile)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	dec := json.NewDecoder(fp)

	err = dec.Decode(&v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func __GetJsonValueInterface(path string, v map[string]interface{}) (val interface{}, types string, err error) {
	var pathext []string
	var curmap map[string]interface{}
	var tmpext []string

	val = ""
	err = nil
	tmpext = strings.Split(path, "/")
	pathext = []string{}
	for _, a := range tmpext {
		if len(a) > 0 {
			pathext = append(pathext, a)
		}
	}

	curmap = v
	if len(pathext) > 0 {

		for i, curpath := range pathext {
			var curval, cval interface{}
			var ok bool
			curval, ok = curmap[curpath]
			if !ok {
				err = fmt.Errorf("can not find (%s) in %s", curpath, path)
				return
			}

			if i == (len(pathext) - 1) {
				val = curval
				if val == nil {
					types = "null"
				} else {
					switch curval.(type) {
					case int:
						types = "int"
					case uint32:
						types = "uint32"
					case uint64:
						types = "uint64"
					case float64:
						types = "float64"
					case float32:
						types = "float32"
					case map[string]interface{}:
						types = "map"
					case []interface{}:
						types = "array"
					case bool:
						types = "bool"
					case string:
						types = "string"
					default:
						err = fmt.Errorf("[%s]unknown type [%s]", pathext[i], reflect.TypeOf(curval).String())
						return
					}
				}
				err = nil
				return
			}

			switch curval.(type) {
			case map[string]interface{}:
				cval, ok = curval.(map[string]interface{})
				if !ok {
					err = fmt.Errorf("can not parse in (%s) for path(%s)", curpath, path)
					return
				}
			case []interface{}:
				cval, ok = curval.([]interface{})
				if !ok {
					err = fmt.Errorf("can not parse in (%s) for path(%s)", curpath, path)
					return
				}
			default:
				err = fmt.Errorf("type of (%s) error", path)
				return
			}
			curmap = cval.(map[string]interface{})
		}
	} else {
		/*we format total*/
		val = v
		types = "map"
		err = nil
		return
	}

	err = fmt.Errorf("can not find (%s) all over", path)
	return
}

func GetJsonValueInterface(path string, v map[string]interface{}) (val interface{}, types string, err error) {
	val, types, err = __GetJsonValueInterface(path, v)
	return
}

func GetJsonValueNull(path string, vmap map[string]interface{}) (err error) {
	var types string
	_, types, err = __GetJsonValueInterface(path, vmap)
	if err != nil {
		return
	}
	if types != "null" {
		err = fmt.Errorf("[%s] type [%s] not null", path, types)
		return
	}
	err = nil
	return
}

func GetJsonValueBool(path string, vmap map[string]interface{}) (val bool, err error) {
	var vinter interface{}
	var types string
	vinter, types, err = __GetJsonValueInterface(path, vmap)
	if err != nil {
		return
	}
	if types != "bool" {
		err = fmt.Errorf("[%s] type [%s] not null", path, types)
		return
	}
	val = vinter.(bool)
	err = nil
	return
}

func GetJsonValueString(path string, vmap map[string]interface{}) (val string, err error) {
	var types string
	var vinter interface{}

	vinter, types, err = __GetJsonValueInterface(path, vmap)
	if err != nil {
		return
	}
	if types != "string" {
		err = fmt.Errorf("[%s] type [%s] not string", path, types)
		return
	}
	val = vinter.(string)
	err = nil
	return
}

func GetJsonValueInt(path string, vmap map[string]interface{}) (val int, err error) {
	var types string
	var vinter interface{}
	var v32 uint32
	var v64 uint64
	var f32 float32
	var f64 float64

	vinter, types, err = __GetJsonValueInterface(path, vmap)
	if err != nil {
		return
	}
	switch types {
	case "int":
		val = vinter.(int)
	case "uint32":
		v32 = vinter.(uint32)
		val = int(v32)
	case "uint64":
		v64 = vinter.(uint64)
		val = int(v64)
	case "float32":
		f32 = vinter.(float32)
		val = int(f32)
	case "float64":
		f64 = vinter.(float64)
		val = int(f64)
	default:
		err = fmt.Errorf("[%s] type %s", path, types)
		return
	}
	err = nil
	return
}

func GetJsonValueFloat(path string, vmap map[string]interface{}) (val float64, err error) {
	var types string
	var vinter interface{}
	var v32 uint32
	var v64 uint64
	var f32 float32
	var vint int

	vinter, types, err = __GetJsonValueInterface(path, vmap)
	if err != nil {
		return
	}
	switch types {
	case "int":
		vint = vinter.(int)
		val = float64(vint)
	case "uint32":
		v32 = vinter.(uint32)
		val = float64(v32)
	case "uint64":
		v64 = vinter.(uint64)
		val = float64(v64)
	case "float32":
		f32 = vinter.(float32)
		val = float64(f32)
	case "float64":
		val = vinter.(float64)
	default:
		err = fmt.Errorf("[%s] type %s", path, types)
		return
	}
	err = nil
	return
}

func GetJsonValueArray(path string, vmap map[string]interface{}) (val []interface{}, err error) {
	var types string
	var vinter interface{}

	vinter, types, err = __GetJsonValueInterface(path, vmap)
	if err != nil {
		return
	}
	if types != "array" {
		err = fmt.Errorf("[%s] not array [%s]", path, types)
		return
	}
	val = vinter.([]interface{})
	err = nil
	return
}

func GetJsonValueMap(path string, vmap map[string]interface{}) (val map[string]interface{}, err error) {
	var types string
	var vinter interface{}

	vinter, types, err = __GetJsonValueInterface(path, vmap)
	if err != nil {
		return
	}
	if types != "map" {
		err = fmt.Errorf("[%s] not map [%s]", path, types)
		return
	}
	val = vinter.(map[string]interface{})
	err = nil
	return
}

func __GetJsonValue(path string, v map[string]interface{}) (val string, err error) {
	var pathext []string
	var curmap map[string]interface{}
	var tmpext []string

	val = ""
	err = nil
	tmpext = strings.Split(path, "/")
	for _, a := range tmpext {
		if len(a) > 0 {
			pathext = append(pathext, a)
		}
	}

	curmap = v
	if len(pathext) > 0 {

		for i, curpath := range pathext {
			var curval, cval interface{}
			var ok bool
			curval, ok = curmap[curpath]
			if !ok {
				err = fmt.Errorf("can not find (%s) in %s", curpath, path)
				return
			}

			if i == (len(pathext) - 1) {
				switch curval.(type) {
				case int:
					val = fmt.Sprintf("%d", curval)
				case uint32:
					val = fmt.Sprintf("%d", curval)
				case uint64:
					val = fmt.Sprintf("%d", curval)
				case float64:
					val = fmt.Sprintf("%f", curval)
				case float32:
					val = fmt.Sprintf("%f", curval)
				case map[string]interface{}:
					val, err = FormatJsonValue(0, "", curval.(map[string]interface{}))
					if err != nil {
						return
					}
				case []interface{}:
					val, err = __FormatValue(0, "", curval)
					if err != nil {
						return
					}
				default:
					val = fmt.Sprintf("%s", curval)
				}
				err = nil
				return
			}

			switch curval.(type) {
			case map[string]interface{}:
				cval, ok = curval.(map[string]interface{})
				if !ok {
					err = fmt.Errorf("can not parse in (%s) for path(%s)", curpath, path)
					return
				}
			case []interface{}:
				cval, ok = curval.([]interface{})
				if !ok {
					err = fmt.Errorf("can not parse in (%s) for path(%s)", curpath, path)
					return
				}
			default:
				err = fmt.Errorf("type of (%s) error", path)
				return
			}
			curmap = cval.(map[string]interface{})
		}
	} else {
		/*we format total*/
		val, err = __FormatValue(0, "", v)
		if err != nil {
			return
		}
		err = nil
		return
	}

	err = fmt.Errorf("can not find (%s) all over", path)
	return
}

func GetJsonValue(path string, v map[string]interface{}) (string, error) {
	return __GetJsonValue(path, v)
}

func GetJsonStruct(valstr string, v interface{}) error {
	var err error
	dec := json.NewDecoder(strings.NewReader(valstr))
	err = dec.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func FormJsonStruct(v interface{}) (valstr string, err error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	valstr = string(b)
	err = nil
	return
}

func GetJsonValueDefault(infile string, path string, defval string) string {
	val := defval
	fp, err := os.Open(infile)
	if err != nil {
		return val
	}
	defer fp.Close()
	dec := json.NewDecoder(fp)

	for {
		var v map[string]interface{}
		err = dec.Decode(&v)
		if err != nil {
			return val
		}

		getval, err := __GetJsonValue(path, v)
		if err == nil {
			val = getval
			return val
		}
	}

	return val
}
