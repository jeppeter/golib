package main

import (
	"dbgutil"
	"encoding/json"
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/jeppeter/go-sqlite3dyn"
	"github.com/tebeka/atexit"
	"logutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func init() {
	Sqlfromfile_handler(nil, nil, nil)
	Sqlfromdir_handler(nil, nil, nil)
	Sqlcreate_handler(nil, nil, nil)
	Sqltime_handler(nil, nil, nil)
	Sqlcode_handler(nil, nil, nil)
}

func LoadParser(parser *extargsparse.ExtArgsParse) (err error) {
	var commandline_fmt string
	var commandline string
	var sqldll string
	var curtime string
	commandline_fmt = `{
		"sqldll" : "%s",
		"starttime|S" : "2000-01-01 00:00:00",
		"endtime|E" : "%s",
		"sqlfromfile<Sqlfromfile_handler>##dbfile tablename jsonfile.. to insert##" : {
			"$" : "+"
		},
		"sqlfromdir<Sqlfromdir_handler>##dbfile tablename dname ... to insert##" : {
			"$" : "+"
		},
		"sqlcreate<Sqlcreate_handler>##dbfile tablename to make database##" : {
			"$" : "+"
		},
		"sqltime<Sqltime_handler>##dbfile tablename [num] to list top values default 10##" : {
			"$" : "+"
		},
		"sqlcode<Sqlcode_handler>##dbfile tablename seccode ... to list seccode ##" : {
			"$" : "+"
		}

	}`

	curtime = time.Now().Format(time.DateTime)

	if runtime.GOOS == "windows" {
		sqldll = ".\\sqlite3.dll"
		sqldll = strings.Replace(sqldll, "\\", "\\\\", -1)
	} else {
		sqldll = "./libsqlite3.so"
	}

	commandline = fmt.Sprintf(commandline_fmt, sqldll, curtime)
	err = parser.LoadCommandLineString(commandline)
	return
}

func Sqlfromfile_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var fdata []byte
	var sarr []string
	var dllfile string
	var fname string
	var idx int
	var conn *sqlite3dyn.Sqlite3BaseConn = nil
	var page *HuigouPage
	var jdx int
	var tablename string
	var sqls string
	var keys, vals string
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 3 {
		err = dbgutil.FormatError("need dbfile tablename jsonfile")
		return
	}

	dllfile = ns.GetString("sqldll")
	err = sqlite3dyn.InitDll(dllfile)
	if err != nil {
		return
	}
	conn, err = sqlite3dyn.ConnSqlite3(sarr[0])
	if err != nil {
		return
	}

	tablename = sarr[1]

	for idx = 2; idx < len(sarr); idx += 1 {
		fname = sarr[idx]
		fdata, err = fileop.ReadFileBytes(fname)
		if err != nil {
			return
		}

		page = &HuigouPage{}
		err = json.Unmarshal(fdata, page)
		if err != nil {
			return
		}

		for jdx = 0; jdx < len(page.Data); jdx += 1 {
			keys, vals, err = page.Data[jdx].FormatInsert()
			if err != nil {
				return
			}
			sqls = fmt.Sprintf(`insert into %s %s values %s;`, tablename, keys, vals)
			err = conn.Exec(sqls, uintptr(0), nil)
			if err != nil {
				return
			}
		}
	}

	conn.Close()

	err = nil
	return
}

func Sqlfromdir_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var dllfile string
	var idx int
	var conn *sqlite3dyn.Sqlite3BaseConn = nil
	var jdx int
	var tablename string
	var sqls string
	var keys, vals string
	var verbmode int
	var filenum int = 0
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 3 {
		err = dbgutil.FormatError("need dbfile tablename jsonfile")
		return
	}

	verbmode = ns.GetInt("verbose")

	dllfile = ns.GetString("sqldll")
	err = sqlite3dyn.InitDll(dllfile)
	if err != nil {
		return
	}
	conn, err = sqlite3dyn.ConnSqlite3(sarr[0])
	if err != nil {
		return
	}

	tablename = sarr[1]

	for idx = 2; idx < len(sarr); idx += 1 {
		err = filepath.Walk(sarr[idx], func(curn string, finfo os.FileInfo, err2 error) (err3 error) {
			var cdata []byte
			var npage *HuigouPage
			err3 = nil
			if finfo.IsDir() {
				return
			}
			cdata, err3 = fileop.ReadFileBytes(curn)
			if err3 != nil {
				logutil.Error("%s", err3.Error())
				err3 = nil
				return
			}
			npage = &HuigouPage{}
			logutil.Debug("curn [%s]", curn)
			err3 = json.Unmarshal(cdata, npage)
			if err3 != nil {
				logutil.Error("[%s] parse error %s", curn, err3.Error())
				err3 = nil
				return
			}

			for jdx = 0; jdx < len(npage.Data); jdx += 1 {
				keys, vals, err3 = npage.Data[jdx].FormatInsert()
				if err3 != nil {
					logutil.Error("format insert error %s", err3.Error())
					err3 = nil
					return
				}
				sqls = fmt.Sprintf(`insert into %s %s values %s;`, tablename, keys, vals)
				err3 = conn.Exec(sqls, uintptr(0), nil)
				if err3 != nil {
					logutil.Error("[%s]insert error %s\n%s", curn, err3.Error(), sqls)
					err3 = nil
					return
				}
			}
			filenum += 1
			if verbmode == 0 {
				if (filenum % 50) == 0 {
					fmt.Fprintf(os.Stdout, ".")
				}

				if (filenum % 500) == 0 {
					fmt.Fprintf(os.Stdout, "\n")
				}
			}
			err3 = nil
			return
		})
	}

	conn.Close()

	err = nil
	return
}

func Sqlcreate_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var dllfile string
	var conn *sqlite3dyn.Sqlite3BaseConn = nil
	var huigou *HuigouInfo
	var tablename string
	var sqls string
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = dbgutil.FormatError("need dbfile tablename")
		return
	}

	dllfile = ns.GetString("sqldll")
	err = sqlite3dyn.InitDll(dllfile)
	if err != nil {
		return
	}
	conn, err = sqlite3dyn.ConnSqlite3(sarr[0])
	if err != nil {
		return
	}

	tablename = sarr[1]

	huigou = &HuigouInfo{}
	sqls, err = huigou.FormatCreateTable(tablename)
	if err != nil {
		return
	}
	err = conn.Exec(sqls, uintptr(0), nil)
	if err != nil {
		return
	}

	conn.Close()
	err = nil
	return
}

func get_pact_info(ptr uintptr, vals []string, cols []string) (err error) {
	var pactpage *HuigouCompactPage
	var curpact *HuigouCompact
	var info *HuigouInfo
	var ok bool

	pactpage = (*HuigouCompactPage)(unsafe.Pointer(ptr))

	info, err = NewFromSqlResult(vals, cols)
	if err != nil {
		logutil.Error("parse error %s\n%v\n%v", err.Error(), vals, cols)
		err = nil
		return
	}

	/*now we should check */
	curpact, ok = pactpage.Data[info.SecurityCode]
	if ok {
		/*now to insert into */
		curpact.ChangeShares += info.ChangeShares
		curpact.ChangeAmount += info.ChangeAmount
		curpact.ChangeDates = append(curpact.ChangeDates, info.ChangeDate)
	} else {
		curpact = &HuigouCompact{}
		curpact.SecurityCode = info.SecurityCode
		curpact.ChangeShares = info.ChangeShares
		curpact.ChangeAmount = info.ChangeAmount
		curpact.ChangeDates = []string{}
		curpact.ChangeDates = append(curpact.ChangeDates, info.ChangeDate)
	}

	if false && pactpage.Verbose >= 3 {
		var sb []byte
		sb, err = json.Marshal(info)
		if err == nil {
			logutil.Debug("info\n%s\nvals\n%v\ncols\n%v", string(sb), vals, cols)
		}

	}

	pactpage.Data[info.SecurityCode] = curpact
	err = nil
	return
}

func Sqltime_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var dllfile string
	var conn *sqlite3dyn.Sqlite3BaseConn = nil
	var tablename string
	var sqls string
	var stime, etime string
	var st, et time.Time
	var pactpage *HuigouCompactPage = nil
	var ptrval uintptr
	var num int = 10
	var huigouarr HuigouCompactArray
	var idx int
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = dbgutil.FormatError("need dbfile tablename")
		return
	}

	dllfile = ns.GetString("sqldll")
	err = sqlite3dyn.InitDll(dllfile)
	if err != nil {
		return
	}
	conn, err = sqlite3dyn.ConnSqlite3(sarr[0])
	if err != nil {
		return
	}

	tablename = sarr[1]

	if len(sarr) > 2 {
		num, err = strconv.Atoi(sarr[2])
		if err != nil {
			return
		}
	}

	stime = ns.GetString("starttime")
	etime = ns.GetString("endtime")

	st, err = get_time_from_str(stime)
	if err != nil {
		err = dbgutil.FormatError("parse [%s] time error %s", stime, err.Error())
		return
	}
	et, err = get_time_from_str(etime)
	if err != nil {
		err = dbgutil.FormatError("parse [%s] time error %s", etime, err.Error())
		return
	}

	sqls = fmt.Sprintf("select * from %s where changedate >= %d and changedate <= %d", tablename, st.Unix(), et.Unix())

	logutil.Debug("sqls\n%s", sqls)
	pactpage = &HuigouCompactPage{}
	/*we set finalizer not make go runtime to automic optimize free the data*/
	runtime.SetFinalizer(pactpage, (*HuigouCompactPage).Close)

	defer pactpage.Close()
	pactpage.Data = make(map[string]*HuigouCompact)
	pactpage.Verbose = ns.GetInt("verbose")

	ptrval = uintptr(unsafe.Pointer(pactpage))
	logutil.Debug("ptrval 0x%x", ptrval)

	err = conn.Exec(sqls, ptrval, get_pact_info)
	if err != nil {
		return
	}

	/*now we should give the matter*/
	huigouarr = []*HuigouCompact{}
	for _, v := range pactpage.Data {
		huigouarr = append(huigouarr, v)
	}

	sort.Sort(huigouarr)

	idx = 0
	for idx < num && idx < len(huigouarr) {
		fmt.Printf("[%d] %s %d %f %v\n", idx, huigouarr[idx].SecurityCode, huigouarr[idx].ChangeShares, huigouarr[idx].ChangeAmount, huigouarr[idx].ChangeDates)
		idx += 1
	}

	conn.Close()
	err = nil
	return
}

func get_code_info(ptr uintptr, vals []string, cols []string) (err error) {
	var pactpage *HuigouPage
	var info *HuigouInfo

	pactpage = (*HuigouPage)(unsafe.Pointer(ptr))

	info, err = NewFromSqlResult(vals, cols)
	if err != nil {
		logutil.Error("parse error %s\n%v\n%v", err.Error(), vals, cols)
		err = nil
		return
	}

	pactpage.Data = append(pactpage.Data, *info)
	err = nil
	return
}

func Sqlcode_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var dllfile string
	var conn *sqlite3dyn.Sqlite3BaseConn = nil
	var tablename string
	var sqls string
	var pactpage *HuigouPage = nil
	var ptrval uintptr
	var idx int
	var jdx int
	if ns == nil {
		err = nil
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = dbgutil.FormatError("need dbfile tablename")
		return
	}

	dllfile = ns.GetString("sqldll")
	err = sqlite3dyn.InitDll(dllfile)
	if err != nil {
		return
	}
	conn, err = sqlite3dyn.ConnSqlite3(sarr[0])
	if err != nil {
		return
	}

	tablename = sarr[1]

	for idx = 2; idx < len(sarr); idx += 1 {
		sqls = fmt.Sprintf(`select * from %s where securitycode = "%s"`, tablename, sarr[idx])
		logutil.Debug("sqls\n%s", sqls)
		pactpage = &HuigouPage{}

		ptrval = uintptr(unsafe.Pointer(pactpage))
		logutil.Debug("ptrval 0x%x", ptrval)

		err = conn.Exec(sqls, ptrval, get_code_info)
		if err != nil {
			return
		}

		for jdx = 0; jdx < len(pactpage.Data); jdx += 1 {

			fmt.Printf("[%d] securite [%s] date [%s] personname %s amount %f shares %d \n", jdx,
				pactpage.Data[jdx].SecurityCode,
				pactpage.Data[jdx].ChangeDate,
				pactpage.Data[jdx].PersonName,
				pactpage.Data[jdx].ChangeAmount,
				pactpage.Data[jdx].ChangeShares)
		}

	}

	conn.Close()
	err = nil
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

	err = LoadParser(parser)
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
