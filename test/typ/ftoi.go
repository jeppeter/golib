package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
)

func is_float_to_int(fv float64) (iv int, bret bool) {
	var mzeros *regexp.Regexp
	var cv *regexp.Regexp
	var vstr string
	var err error
	var dotstring string
	var numstr string
	vstr = fmt.Sprintf("%v", fv)
	mzeros = regexp.MustCompile(`^[0-9]+$`)
	cv = regexp.MustCompile(`^[0-9]\.([0-9]+)[eE][\-\+]([0-9]+)$`)
	iv = 0
	bret = false
	if mzeros.MatchString(vstr) {
		iv, err = strconv.Atoi(vstr)
		if err == nil {
			bret = true
		}
		return
	}
	if cv.MatchString(vstr) {

	}

	return
}

func get_json(k string) error {
	var vstr string
	var vmap map[string]interface{}
	var err error
	var fv float64
	var v interface{}
	var iv int
	var bret bool

	vstr = fmt.Sprintf(`{"code": %s}`, k)
	err = json.Unmarshal([]byte(vstr), &vmap)
	if err != nil {
		return err
	}
	v = vmap["code"]
	switch v.(type) {
	case uint32:
		fv = float64(v.(uint32))
	case uint16:
		fv = float64(v.(uint16))
	case uint64:
		fv = float64(v.(uint64))
	case int16:
		fv = float64(v.(int16))
	case int32:
		fv = float64(v.(int32))
	case int64:
		fv = float64(v.(int64))
	case float64:
		fv = v.(float64)
	case float32:
		fv = float64(v.(float32))
	default:
		return fmt.Errorf("unknown type [%s]", reflect.ValueOf(v).Type().String())
	}
	iv, bret = is_float_to_int(fv)
	if !bret {
		fmt.Fprintf(os.Stdout, "[%s] not int\n", k)
		return nil
	}
	fmt.Fprintf(os.Stdout, "[%s]=[%d]\n", k, iv)
	return nil
}

func main() {
	for _, c := range os.Args[1:] {
		get_json(c)
	}
	return
}
