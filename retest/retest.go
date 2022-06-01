package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
	"regexp"
)

func Findall_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var restr string
	var idx int
	var sarr []string
	var reg *regexp.Regexp
	var matchstrings []string
	var j int
	var s string
	err = nil
	if ns == nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = fmt.Errorf("findall need at least 2 args")
		return
	}
	restr = sarr[0]
	reg, err = regexp.Compile(restr)
	if err != nil {
		err = fmt.Errorf("compile [%s] error[%s]", restr, err.Error())
		return
	}

	for idx = 1; idx < len(sarr); idx++ {
		matchstrings = reg.FindStringSubmatch(sarr[idx])
		if len(matchstrings) > 0 {
			fmt.Fprintf(os.Stdout, "[%s] find all in [%s]\n", restr, sarr[idx])
			for j, s = range matchstrings {
				fmt.Fprintf(os.Stdout, "\t[%d] [%s]\n", j, s)
			}
		} else {
			fmt.Fprintf(os.Stdout, "[%s] not find all in [%s]\n", restr, sarr[idx])
		}
	}
	err = nil
	return
}

func Ifindall_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var restr string
	var idx int
	var sarr []string
	var reg *regexp.Regexp
	var matchstrings []string
	var j int
	var s string
	err = nil
	if ns == nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = fmt.Errorf("findall need at least 2 args")
		return
	}
	restr = sarr[0]
	reg, err = regexp.Compile("(?i)" + restr)
	if err != nil {
		err = fmt.Errorf("compile [%s] error[%s]", restr, err.Error())
		return
	}

	for idx = 1; idx < len(sarr); idx++ {
		matchstrings = reg.FindStringSubmatch(sarr[idx])
		if len(matchstrings) > 0 {
			fmt.Fprintf(os.Stdout, "[%s] find all in [%s]\n", restr, sarr[idx])
			for j, s = range matchstrings {
				fmt.Fprintf(os.Stdout, "\t[%d] [%s]\n", j, s)
			}
		} else {
			fmt.Fprintf(os.Stdout, "[%s] not find all in [%s]\n", restr, sarr[idx])
		}
	}
	err = nil
	return
}

func Findindex_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var restr string
	var reg *regexp.Regexp
	var indexes []int
	var idx int
	err = nil
	if ns == nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = fmt.Errorf("findindex need at least 2 args")
		return
	}

	restr = sarr[0]
	reg, err = regexp.Compile(restr)
	if err != nil {
		err = fmt.Errorf("compile [%s] error[%s]", restr, err.Error())
		return
	}

	for idx = 1; idx < len(sarr); idx++ {
		indexes = reg.FindStringIndex(sarr[idx])
		if len(indexes) > 0 {
			fmt.Fprintf(os.Stdout, "[%s] find all in [%s]\n", restr, sarr[idx])
			fmt.Fprintf(os.Stdout, "\t[%d] [%d] [%s]\n", indexes[0], indexes[1], sarr[idx][indexes[0]:indexes[1]])
		} else {
			fmt.Fprintf(os.Stdout, "[%s] not find all in [%s]\n", restr, sarr[idx])
		}
	}

	err = nil
	return
}

func Ifindindex_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var restr string
	var reg *regexp.Regexp
	var indexes []int
	var idx int
	err = nil
	if ns == nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = fmt.Errorf("findindex need at least 2 args")
		return
	}

	restr = sarr[0]
	reg, err = regexp.Compile("(?i)" + restr)
	if err != nil {
		err = fmt.Errorf("compile [%s] error[%s]", restr, err.Error())
		return
	}

	for idx = 1; idx < len(sarr); idx++ {
		indexes = reg.FindStringIndex(sarr[idx])
		if len(indexes) > 0 {
			fmt.Fprintf(os.Stdout, "[%s] find all in [%s]\n", restr, sarr[idx])
			fmt.Fprintf(os.Stdout, "\t[%d] [%d] [%s]\n", indexes[0], indexes[1], sarr[idx][indexes[0]:indexes[1]])
		} else {
			fmt.Fprintf(os.Stdout, "[%s] not find all in [%s]\n", restr, sarr[idx])
		}
	}

	err = nil
	return
}

func Match_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var idx int
	var restr string
	var reg *regexp.Regexp
	err = nil
	if ns == nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = fmt.Errorf("match need at least 2 args")
		return
	}

	restr = sarr[0]
	reg, err = regexp.Compile(restr)
	if err != nil {
		err = fmt.Errorf("compile [%s] error[%s]", restr, err.Error())
		return
	}

	for idx = 1; idx < len(sarr); idx++ {
		if reg.MatchString(sarr[idx]) {
			fmt.Fprintf(os.Stdout, "[%s] match [%s]\n", sarr[0], sarr[idx])
		} else {
			fmt.Fprintf(os.Stdout, "[%s] not match [%s]\n", sarr[0], sarr[idx])
		}
	}

	return
}

func Imatch_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var idx int
	var restr string
	var reg *regexp.Regexp
	err = nil
	if ns == nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = fmt.Errorf("match need at least 2 args")
		return
	}

	restr = sarr[0]
	reg, err = regexp.Compile("(?i)" + restr)
	if err != nil {
		err = fmt.Errorf("compile [%s] error[%s]", restr, err.Error())
		return
	}

	for idx = 1; idx < len(sarr); idx++ {
		if reg.MatchString(sarr[idx]) {
			fmt.Fprintf(os.Stdout, "[%s] match [%s]\n", sarr[0], sarr[idx])
		} else {
			fmt.Fprintf(os.Stdout, "[%s] not match [%s]\n", sarr[0], sarr[idx])
		}
	}

	return
}

func Split_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var idx int
	var restr string
	var reg *regexp.Regexp
	var splitstrings []string
	var j int
	var s string
	err = nil
	if ns == nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = fmt.Errorf("split need at least 2 args")
		return
	}

	restr = sarr[0]
	reg, err = regexp.Compile(restr)
	if err != nil {
		err = fmt.Errorf("compile [%s] error[%s]", restr, err.Error())
		return
	}

	for idx = 1; idx < len(sarr); idx++ {
		splitstrings = reg.Split(sarr[idx], -1)
		if len(splitstrings) > 0 {
			fmt.Fprintf(os.Stdout, "[%s] split with [%s]\n", sarr[idx], sarr[0])
			for j, s = range splitstrings {
				fmt.Fprintf(os.Stdout, "\t[%d] [%s]\n", j, s)
			}
		} else {
			fmt.Fprintf(os.Stdout, "[%s] nothing with split [%s]\n", sarr[idx], sarr[0])
		}
	}

	return
}

func Replace_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var rplstr string
	var instr string
	var restr string
	var result string
	var bs []byte
	var re *regexp.Regexp
	var idx int

	err = nil
	if ns == nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 3 {
		err = fmt.Errorf("replace need at least 3 args")
		return
	}

	restr = sarr[0]
	rplstr = sarr[1]
	idx = 2
	re, err = regexp.Compile(restr)
	if err != nil {
		err = fmt.Errorf("compile [%s] error[%s]", restr, err.Error())
		return
	}

	for idx = 2; idx < len(sarr); idx++ {
		instr = sarr[idx]
		bs = re.ReplaceAll([]byte(instr), []byte(rplstr))
		result = string(bs)
		fmt.Fprintf(os.Stdout, "replace [%s]  [%s]([%s]) => [%s]\n", instr, restr, rplstr, result)
	}
	err = nil
	return
}

func init() {
	Match_handler(nil, nil, nil)
	Imatch_handler(nil, nil, nil)
	Findall_handler(nil, nil, nil)
	Ifindall_handler(nil, nil, nil)
	Findindex_handler(nil, nil, nil)
	Ifindindex_handler(nil, nil, nil)
	Split_handler(nil, nil, nil)
	Replace_handler(nil, nil, nil)
}

func main() {
	var commandline = `
	{
		"match<Match_handler>## restr instr ... to find match##" : {
			"$" : "+"
		},
		"imatch<Imatch_handler>## restr instr ... to find match##" : {
			"$" : "+"
		},
		"findall<Findall_handler>## restr instr ... to find all matches##" : {
			"$" : "+"
		},
		"findindex<Findindex_handler>## restr instr ... to find index matches##" : {
			"$" : "+"
		},
		"ifindall<Ifindall_handler>## restr instr ... to find all matches in ignore case##" : {
			"$" : "+"
		},
		"ifindindex<Ifindindex_handler>## restr instr ... to find index matches in ignore case##" : {
			"$" : "+"
		},
		"split<Split_handler>## restr instr ... to split string with regexp##" : {
			"$" : "+"
		},
		"replace<Replace_handler>## restr rplstr instr ... to replace match with rplstr ##" : {
			"$" : "+"
		}
	}
	`
	var parser *extargsparse.ExtArgsParse
	var err error

	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not make parser err[%s]\n", err.Error())
		os.Exit(5)
	}
	err = parser.LoadCommandLineString(commandline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not parse %s\n", commandline)
		os.Exit(5)
	}

	_, err = parser.ParseCommandLineEx(nil, nil, nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not use parse command line [%s]\n", err.Error())
		os.Exit(4)
	}
	return
}
