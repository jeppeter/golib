package main

import (
	"crypto/x509"
	"dbgutil"
	"time"
)

func trans_inter_to_string(valarr []interface{}, note string) (arrs []string, err error) {
	arrs = []string{}
	var idx int
	var s string
	var ok bool
	for idx = 0; idx < len(valarr); idx += 1 {
		s, ok = valarr[idx].(string)
		if !ok {
			err = dbgutil.FormatError("[%s].[%s] not string", note, idx)
			return
		}
		arrs = append(arrs, s)
	}
	err = nil
	return
}

func get_time_value(times string, note string) (retv time.Time, err error) {
	retv, err = time.Parse("2020-02-02 13:20:50", times)
	if err != nil {
		err = dbgutil.FormatError("[%s] [%s] parse error %s please use [2020-02-02 13:20:50] format", note, times, err.Error())
		return
	}
	return
}

func get_key_ext_usage(arrs []string, note string) (retv []x509.ExtKeyUsage, err error) {
	retv = []x509.ExtKeyUsage{}
	var idx, jdx int
	var matched bool
	for idx = 0; idx < len(arrs); idx += 1 {
		matched = false
		for jdx = 0; jdx < len(extKeyUsageValue); jdx += 1 {
			if extKeyUsageValue[jdx].key == arrs[idx] {
				matched = true
				retv = append(retv, extKeyUsageValue[jdx].extKeyUsage)
				break
			}
		}

		if !matched {
			err = dbgutil.FormatError("[%s].[%d][%s] not supported", note, idx, arrs[idx])
			return
		}
	}
	err = nil
	return
}

func get_int_value(val interface{}, note string) (ival int, err error) {
	var vali int
	var valf float64
	var ok bool
	vali, ok = val.(int)
	if ok {
		ival = vali
	} else {
		valf, ok = val.(float64)
		if ok {
			ival = int(valf)
		} else {
			err = dbgutil.FormatError("[%s] not int value", note)
			return
		}
	}
	err = nil
	return
}
