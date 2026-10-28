package main

import (
	"context"
	"dbgutil"
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"gopkg.in/yaml.v2"
	"logutil"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"randfile"
	"sort"
	"strconv"
	"time"
)

func init() {
	Tempfile_handler(nil, nil, nil)
	Tempdir_handler(nil, nil, nil)
	Walkdir_handler(nil, nil, nil)
	Sig_handler(nil, nil, nil)
	Yaml_handler(nil, nil, nil)
	Rand_handler(nil, nil, nil)
}

func Tempfile_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var idx int
	var pattern string
	var realfile string
	var dirn string
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	dirn = ns.GetString("directory")

	sarr = ns.GetArray("subnargs")

	for idx = 0; idx < len(sarr); idx += 1 {
		pattern = sarr[idx]
		realfile, err = fileop.Mktempfile(dirn, pattern)
		if err != nil {
			return
		}
		fmt.Printf("[%s].[%s] => [%s]\n", dirn, pattern, realfile)
	}
	err = nil
	return

}

func Tempdir_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var idx int
	var pattern string
	var realfile string
	var dirn string
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	dirn = ns.GetString("directory")

	sarr = ns.GetArray("subnargs")

	for idx = 0; idx < len(sarr); idx += 1 {
		pattern = sarr[idx]
		realfile, err = fileop.Mktempdir(dirn, pattern)
		if err != nil {
			return
		}
		fmt.Printf("[%s].[%s] => [%s]\n", dirn, pattern, realfile)
	}
	err = nil
	return
}

func Walkdir_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var idx int
	var curdir string
	var totaln int = 0
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")

	for idx = 0; idx < len(sarr); idx += 1 {
		curdir = sarr[idx]
		err = filepath.Walk(curdir, func(curn string, info os.FileInfo, err2 error) error {
			if !info.IsDir() {
				totaln += 1
			}
			fmt.Printf("curdir [%s] curn %s\n", curdir, curn)
			return nil
		})
	}
	fmt.Printf("totaln %d\n", totaln)
	err = nil
	return
}

func Sig_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx1 interface{}) (err error) {
	var ctx context.Context
	var stop context.CancelFunc
	var cnt int = 0
	var exited int
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	ctx, stop = signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	for {
		exited = 0
		select {
		case <-ctx.Done():
			logutil.Debug("get ctx Done")
			exited = 1
		case <-time.After(time.Second * 1):
			logutil.Debug("cnt %d", cnt)
		}
		if exited != 0 {
			break
		}

		cnt += 1
	}

	logutil.Debug("exit cnt %d", cnt)
	return
}

type StringArray []string

func (a StringArray) Len() int {
	return len(a)
}

func (a StringArray) Less(i, j int) bool {
	var bval bool = false
	if a[i] < a[j] {
		bval = true
	}
	return bval
}

func (a StringArray) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func Yaml_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var yamlfile string
	var vmap map[string]interface{}
	var outb []byte
	var keys StringArray
	var output string
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")

	yamlfile = sarr[0]
	outb, err = fileop.ReadFileBytes(yamlfile)
	if err != nil {
		return
	}

	vmap = make(map[string]interface{})
	err = yaml.Unmarshal(outb, &vmap)
	if err != nil {
		err = dbgutil.FormatError("can not parse [%s] error %s", yamlfile, err.Error())
		return
	}

	keys = []string{}
	for k, _ := range vmap {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	for _, k := range keys {
		fmt.Printf("[%s]=%v\n", k, vmap[k])
	}

	output = ns.GetString("output")
	if len(output) > 0 {
		outb, err = yaml.Marshal(vmap)
		if err != nil {
			err = dbgutil.FormatError("can not marshal %s error %s", yamlfile, err.Error())
			return
		}
		_, err = fileop.WriteFileBytes(output, outb)
		if err != nil {
			return
		}
	}
	err = nil
	return

}

func Rand_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var randf *randfile.RandFile
	var input string
	var orand *rand.Rand
	var number int = 10
	var inb []byte
	var idx int
	err = nil

	if ns == nil {
		return
	}

	err = logutil.InitLog(ns)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	input = ns.GetString("input")
	if len(input) > 0 {
		randf, err = randfile.NewRandFile(input)
		if err != nil {
			return
		}
	} else {
		randf, err = randfile.NewRandFp(os.Stdin, "stdin")
		if err != nil {
			return
		}
	}

	defer randf.Close()
	orand = rand.New(randf)

	if len(sarr) > 0 {
		number, err = strconv.Atoi(sarr[0])
		if err != nil {
			err = dbgutil.FormatError("[%s] not number", sarr[0])
			return
		}
	}
	inb = make([]byte, 1)
	for idx = 0; idx < number; idx += 1 {

		_, err = orand.Read(inb)
		if err != nil {
			return
		}
		fmt.Printf("[%d] value 0x%x\n", idx, inb[0])

	}

	err = nil
	return

}

func main() {
	var commandline string
	var err error
	var parser *extargsparse.ExtArgsParse
	var ns *extargsparse.NameSpaceEx

	commandline = `{
		"input|i" : null,
		"output|o" : null,
		"directory|D" : null,
		"tempfile<Tempfile_handler>##pattern to create file##" : {
			"$" : "+"
		},
		"tempdir<Tempdir_handler>##pattern to create file##" : {
			"$" : "+"
		},
		"walkdir<Walkdir_handler>##dir... to scan dir##" : {
			"$" : "+"
		},
		"sig<Sig_handler>##to demonstrate the signal##" : {
			"$" : 0
		},
		"yaml<Yaml_handler>##file to decode yaml##" : {
			"$" : 1
		},
		"rand<Rand_handler>##num ... to make rand number##" : {
			"$" : "*"
		}

	}`

	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not make parser err[%s]\n", err.Error())
		atexit.Exit(5)
	}

	err = logutil.PrepareLog(parser)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not set [%s]\n", err.Error())
		atexit.Exit(5)
	}
	err = parser.LoadCommandLineString(commandline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not parse %s\n", commandline)
		atexit.Exit(5)
	}

	ns, err = parser.ParseCommandLineEx(nil, nil, nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not use parse command line [%s]\n", err.Error())
		atexit.Exit(4)
	}
	if len(ns.GetString("subcommand")) == 0 {
		fmt.Fprintf(os.Stderr, "can not get subcommand\n")
		atexit.Exit(5)
	}
	//fmt.Fprintf(os.Stdout, "subcommand [%s] succ\n", ns.GetString("subcommand"))
	atexit.Exit(0)
	return
}
