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
	GGEid              string  `json:"GGEID"`
	BeginHoldNum       int64   `json:"BEGIN_HOLD_NUM"`
	EndHoldNum         int64   `json:"END_HOLD_NUM"`
}

type HuigouPage struct {
	Data  []HuigouInfo `json:"data"`
	Pages int64        `json:"pages"`
	Count int64        `json:"count"`
}

func (hg *HuigouInfo) FormatInsert(dbname string) (outs string, err error) {
	var keys string
	var vals string
	outs = fmt.Sprintf("insert into %s ", dbname)

	keys = "("
	vals = "("

	keys += "changedate"
	vals += fmt.Sprintf(`"%s"`, hg.ChangeDate)

	keys += ",securitycode"
	vals += fmt.Sprintf(`,"%s"`, hg.SecurityCode)

	keys += ",derivesecurity"
	vals += fmt.Sprintf(`,"%s"`, hg.DeriveSecurity)

	keys += ",changeshares"
	vals += fmt.Sprintf(`, %d`, hg.ChangeShares)

	keys += ")"
	vals += ")"

	outs += keys
	outs += " values "
	outs += vals
	outs += ";"
	err = nil
	return
}
