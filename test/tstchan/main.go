package main

import (
	"dbgutil"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"logutil"
	"math/rand"
	"reflect"
	"slices"
	"strconv"
	"time"
)

func init() {
	Selectchan_handler(nil, nil, nil)
}

type ChanTest struct {
	idx     int
	wtime   int
	notechl chan int
}

func NewChanTest(idx int, wtime int) (retp *ChanTest, err error) {
	retp = &ChanTest{}
	retp.idx = idx
	retp.wtime = wtime
	retp.notechl = make(chan int, 10)
	err = nil
	return
}

func (retp *ChanTest) MainProc() {
	logutil.Debug("run in [%d] wtime [%d]", retp.idx, retp.wtime)
	time.Sleep(time.Duration(retp.wtime) * time.Millisecond)
	logutil.Debug("exit [%d]", retp.idx)
	retp.notechl <- 1
	return
}

func Selectchan_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var numchl int = 10
	var idx int
	var chls []*ChanTest = []*ChanTest{}
	var curchl *ChanTest
	var nowt time.Time
	var wtime int
	var cases []reflect.SelectCase
	var chosen int
	var ok bool

	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) > 0 {
		numchl, err = strconv.Atoi(sarr[0])
		if err != nil {
			err = dbgutil.FormatError("parse [%s] error %s", sarr[0], err.Error())
			return
		}
	}

	nowt = time.Now()
	rand.Seed(int64(nowt.Unix()))
	for idx = 0; idx < numchl; idx += 1 {
		wtime = rand.Intn(2000)
		curchl, err = NewChanTest(idx, wtime)
		go curchl.MainProc()
		chls = append(chls, curchl)
	}

	for len(chls) > 0 {
		logutil.Debug("chls [%d]", len(chls))
		cases = make([]reflect.SelectCase, len(chls))
		for idx = 0; idx < len(chls); idx += 1 {
			logutil.Debug("idx [%d] [%d]", idx, chls[idx].idx)
			cases[idx] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(chls[idx].notechl),
			}
		}
		chosen, _, ok = reflect.Select(cases)
		if ok && chosen < len(chls) {
			logutil.Debug("chosen %d", chosen)
			chls = slices.Delete(chls, chosen, chosen+1)
		}
	}

	return nil
}

func LoadRegCmdFlags(parser *extargsparse.ExtArgsParse) (err error) {
	var commandline_fmt string
	var commandline string
	commandline_fmt = `{
		"selchan<Selectchan_handler>##[num] number of channel handle default 10##" : {
			"$" : "?"
		}
	}`

	commandline = fmt.Sprintf(commandline_fmt)
	err = parser.LoadCommandLineString(commandline)
	return
}

func main() {
	var parser *extargsparse.ExtArgsParse
	var err error
	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		logutil.Error("%s", err.Error())
		atexit.Exit(5)
	}

	err = LoadRegCmdFlags(parser)
	if err != nil {
		logutil.Error("%s", err.Error())
		atexit.Exit(5)
	}

	err = logutil.PrepareLog(parser)
	if err != nil {
		logutil.Error("%s", err.Error())
		atexit.Exit(5)
	}

	_, err = parser.ParseCommandLine(nil, nil)
	if err != nil {
		logutil.Error("%s", err.Error())
		atexit.Exit(4)
	}
	atexit.Exit(0)
}
