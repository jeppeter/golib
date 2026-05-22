package main

import (
	"dbgutil"
	"fmt"
	"logutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type HuigouInfo struct {
	SecurityCode       string  `json:"SECURITY_CODE"`
	DeriveSecurity     string  `json:"DERIVE_SECURITY_CODE"`
	SecurityName       string  `json:"SECURITY_NAME"`
	ChangeDate         string  `json:"CHANGE_DATE"`
	PersonName         string  `json:"PERSON_NAME"`
	ChangeShares       int64   `json:"CHANGE_SHARES"`
	AveragePrice       float64 `json:"AVERAGE_PRICE"`
	ChangeAmount       float64 `json:"CHANGE_AMOUNT"`
	ChangeReason       string  `json:"CHANGE_REASON"`
	ChangeRatio        float64 `json:"CHANGE_RATIO"`
	ChangeAfterHoldNum int64   `json:"CHANGE_AFTER_HOLDNUM"`
	HoldType           string  `json:"HOLD_TYPE"`
	DsePersonName      string  `json:"DSE_PERSON_NAME"`
	PositionName       string  `json:"POSITION_NAME"`
	PersonDseRelation  string  `json:"PERSON_DSE_RELATION"`
	OrgCode            string  `json:"ORG_CODE"`
	GGEid              int64   `json:"GGEID"`
	BeginHoldNum       int64   `json:"BEGIN_HOLD_NUM"`
	EndHoldNum         int64   `json:"END_HOLD_NUM"`
}

func (hg *HuigouInfo) get_epoch(s string) (rval int64, err error) {
	var ntime time.Time
	rval = 0

	ntime, err = get_time_from_str(s)
	//logutil.Debug("ntime %v s %s", ntime, s)
	rval = ntime.Unix()

	err = nil
	return
}

func (hg *HuigouInfo) FormatCreateTable(table string) (outs string, err error) {
	var kname string
	var pt reflect.Type
	var rf reflect.Value
	var i int
	var types string
	var reg *regexp.Regexp
	var bmatch bool

	reg, err = regexp.Compile(".*date.*")

	err = nil
	outs = ""
	outs += fmt.Sprintf("create table %s (", table)

	rf = reflect.ValueOf(hg).Elem()
	pt = rf.Type()
	for i = 0; i < rf.NumField(); i += 1 {
		kname = strings.ToLower(pt.Field(i).Name)
		types = pt.Field(i).Type.String()

		if i > 0 {
			outs += ","
		}

		if types == "int64" {
			outs += fmt.Sprintf("%s int64", kname)
		} else if types == "string" {
			/**/
			bmatch = reg.MatchString(kname)
			if bmatch {
				outs += fmt.Sprintf("%s int64", kname)
			} else {
				outs += fmt.Sprintf("%s text", kname)
			}

		} else if types == "float64" {
			outs += fmt.Sprintf("%s double", kname)
		}

	}
	outs += ")"

	err = nil
	return
}

func (hg *HuigouInfo) FormatInsert() (keys string, vals string, err error) {
	var kname string
	var pt reflect.Type
	var rf, vrf reflect.Value
	var i int
	var types string
	var reg *regexp.Regexp
	var bmatch bool
	var iv64 int64
	var fv64 float64
	var s string
	var ok bool

	reg, err = regexp.Compile(".*date.*")

	err = nil
	keys = "("
	vals = "("

	rf = reflect.ValueOf(hg).Elem()
	pt = rf.Type()
	for i = 0; i < rf.NumField(); i += 1 {
		kname = strings.ToLower(pt.Field(i).Name)
		types = pt.Field(i).Type.String()

		if i > 0 {
			keys += ","
			vals += ","
		}

		vrf = rf.Field(i)
		vrf = reflect.NewAt(vrf.Type(), unsafe.Pointer(vrf.UnsafeAddr())).Elem()

		if types == "int64" {
			iv64, ok = vrf.Interface().(int64)
			if !ok {
				err = dbgutil.FormatError("%s not int64", kname)
				return
			}
			keys += kname
			vals += fmt.Sprintf("%d", iv64)
		} else if types == "string" {
			/**/
			s, ok = vrf.Interface().(string)
			if !ok {
				err = dbgutil.FormatError("%s not string", kname)
				return
			}
			bmatch = reg.MatchString(kname)
			if bmatch {
				iv64, err = hg.get_epoch(s)
				if err != nil {
					return
				}
				keys += kname
				vals += fmt.Sprintf("%d", iv64)
			} else {
				s = strings.Replace(s, "\"", "", -1)
				keys += kname
				vals += fmt.Sprintf(`"%s"`, s)
			}

		} else if types == "float64" {
			fv64, ok = vrf.Interface().(float64)
			if !ok {
				err = dbgutil.FormatError("%s not float64", kname)
				return
			}
			keys += kname
			vals += fmt.Sprintf("%f", fv64)
		}

	}

	keys += ")"
	vals += ")"
	return
}

func NewFromSqlResult(vals []string, cols []string) (retp *HuigouInfo, err error) {
	var idx int = 0
	var jdx int = 0
	var matched bool
	var pt reflect.Type
	var rf, vrf reflect.Value
	var colname, kname string
	var types string
	var iv64 int64
	var fv64 float64

	retp = &HuigouInfo{}

	rf = reflect.ValueOf(retp).Elem()
	pt = rf.Type()

	for idx = 0; idx < len(cols); idx += 1 {
		colname = cols[idx]
		matched = false
		for jdx = 0; jdx < rf.NumField(); jdx += 1 {
			kname = strings.ToLower(pt.Field(jdx).Name)
			if kname == colname {
				/**/
				types = pt.Field(jdx).Type.String()

				vrf = rf.Field(jdx)
				vrf = reflect.NewAt(vrf.Type(), unsafe.Pointer(vrf.UnsafeAddr())).Elem()

				if types == "int64" {
					/*now we should give the int*/
					iv64, err = strconv.ParseInt(vals[idx], 10, 64)
					if err != nil {
						logutil.Error("parse [%s] %s error %s", kname, vals[idx], err.Error())
						return
					}
					vrf.Set(reflect.ValueOf(iv64))
				} else if types == "float64" {
					fv64, err = strconv.ParseFloat(vals[idx], 64)
					if err != nil {
						logutil.Error("parse [%s] %s error %s", kname, vals[idx], err.Error())
						return
					}
					vrf.Set(reflect.ValueOf(fv64))
				} else if types == "string" {
					if kname == "changedate" {
						iv64, err = strconv.ParseInt(vals[idx], 10, 64)
						if err != nil {
							logutil.Error("parse [%s] %s error %s", kname, vals[idx], err.Error())
							return
						}

						/*now to parse*/
						var curtime time.Time
						curtime = time.Unix(iv64, 0)
						retp.ChangeDate = curtime.Format(time.DateTime)
					} else {
						vrf.Set(reflect.ValueOf(vals[idx]))
					}
				}

				matched = true
				break
			}
		}

		if !matched {
			logutil.Debug("not matched %s", colname)
		}
	}

	err = nil
	return

}

type HuigouPage struct {
	Data  []HuigouInfo `json:"data"`
	Pages int64        `json:"pages"`
	Count int64        `json:"count"`
}

type HuigouCompact struct {
	SecurityCode string   `json:"SECURITY_CODE"`
	ChangeShares int64    `json:"CHANGE_SHARES"`
	ChangeAmount float64  `json:"CHANGE_AMOUNT"`
	ChangeDates  []string `json:"CHANGE_DATES"`
}

type HuigouCompactPage struct {
	Data    map[string]*HuigouCompact `json:"data"`
	Verbose int
}

func (hp *HuigouCompactPage) Close() {
	hp.Data = make(map[string]*HuigouCompact)
	return
}

type HuigouCompactArray []*HuigouCompact

func (a HuigouCompactArray) Len() int {
	return len(a)
}

func (a HuigouCompactArray) Less(i, j int) bool {
	return a[i].ChangeAmount > a[j].ChangeAmount
}

func (a HuigouCompactArray) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
