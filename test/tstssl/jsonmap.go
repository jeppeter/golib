package main

import (
	"crypto/x509"
	"dbgutil"
	"fmt"
	"logutil"
	"strconv"
	"strings"
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

const TIME_DEFAULT_FOMRAT = "2020-02-02 13:20:50"

func get_time_value(times string, note string) (retv time.Time, err error) {
	var sarr []string
	var year int
	var mon int
	var mday int
	var hour int
	var min int
	var sec int
	var msarr []string
	var hsarr []string
	var fmts string

	sarr = strings.SplitN(times, " ", 2)
	if len(sarr) < 2 {
		err = dbgutil.FormatError("need %s format", TIME_DEFAULT_FOMRAT)
		return
	}

	msarr = strings.SplitN(sarr[0], "-", 3)
	if len(msarr) < 3 {
		err = dbgutil.FormatError("need %s format", TIME_DEFAULT_FOMRAT)
		return
	}

	year, err = strconv.Atoi(msarr[0])
	if err != nil {
		err = dbgutil.FormatError("need %s format year %s not valid", TIME_DEFAULT_FOMRAT, msarr[0])
		return
	}

	mon, err = strconv.Atoi(msarr[1])
	if err != nil {
		err = dbgutil.FormatError("need %s format mon %s not valid", TIME_DEFAULT_FOMRAT, msarr[1])
		return
	}

	mday, err = strconv.Atoi(msarr[2])
	if err != nil {
		err = dbgutil.FormatError("need %s format mday %s not valid", TIME_DEFAULT_FOMRAT, msarr[2])
		return
	}

	hsarr = strings.SplitN(sarr[1], ":", 3)
	if len(hsarr) < 3 {
		err = dbgutil.FormatError("need %s format", TIME_DEFAULT_FOMRAT)
		return
	}

	hour, err = strconv.Atoi(hsarr[0])
	if err != nil {
		err = dbgutil.FormatError("need %s format hour %s not valid", TIME_DEFAULT_FOMRAT, hsarr[0])
		return
	}

	min, err = strconv.Atoi(hsarr[1])
	if err != nil {
		err = dbgutil.FormatError("need %s format min %s not valid", TIME_DEFAULT_FOMRAT, hsarr[1])
		return
	}
	sec, err = strconv.Atoi(hsarr[2])
	if err != nil {
		err = dbgutil.FormatError("need %s format sec %s not valid", TIME_DEFAULT_FOMRAT, hsarr[2])
		return
	}

	fmts = fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", year, mon, mday, hour, min, sec)
	logutil.Debug("times [%s] fmts [%s]", times, fmts)
	retv, err = time.Parse(time.RFC3339, fmts)
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
