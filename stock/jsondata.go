package main

import (
	"json"
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

	ntime, err = time.Parse(time.DateTime, s)
	logutil.Debug("ntime %v s %s", ntime, s)
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

type HuigouPage struct {
	Data  []HuigouInfo `json:"data"`
	Pages int64        `json:"pages"`
	Count int64        `json:"count"`
}
