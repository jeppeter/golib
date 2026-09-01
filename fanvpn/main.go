package main

import (
	"context"
	"dbgutil"
	"fileop"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/tebeka/atexit"
	"logutil"
	"os"
	"os/signal"
	"tcprelay"
	"time"
)

const PEM_DEFAULT = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCUHMdMJFSzOnuO
OIRY3UZcWTXfCexZ15UOXjZB8fvgLxyvTN5gTh+Wo6SNb48ojuIw/p3fltMom3ZW
eyeKD5QB/sL+LgWHs10GrCE73TYwsmhdXUCyZhZX7VK8meqsFBBf7SDTPcJ6lfup
XsFstxQWV9nq2FQAQLSzUVU46KGjADKV7cEwamQ+C+0ix74EfDYvaZohfUtHEMie
MB5xYePJnxPS1qqP4Ftsgco5/hZfscsf9lR4+LSDAj+krNRo2j6Kto030Kxrwymi
7KbfVvlXujYsoh+LlhEXEAxHYUTbllKvcKRWvEiMNkW64+lFMclVIkHs0KjRPrOY
0nfy0iVJAgMBAAECggEAQdrP5HOM84nv0OshMW/lbn89/Dcpz0KTJGnQXx7seqAH
9YvMnm5uDikhq79sHED3oog7guRJbBc/lTE6AeFuUjrH0YN98vnVxXc4aakwhJN2
4vhpIUlR6vN7I5+eH7fmFfjV7QbbV20jkgmvIBsBA/Q40Pox01Dx538k0OJiqBnr
h+jdL99lURPYkG3/mfT3R2pG+vIP2PYydW0pi87f6AK4pxFZnAF82MuGcCM+Byr0
ga8/pSV4wmkU+kezK4P6RfbSi4tuUszvaBWcczAgqEN3yltYwN0v/JSk2oHcBSx4
phJV/RGGGg/IRHt4d0zi3PRfTrx2YI7AEgc1TCVSxwKBgQDMTdr6j64IMtWM45TG
QMcwtZLDutdCpvkVBGtpzdj3DWP2ZYAWQoO1ko1T2IOKkNWIXO3PEmX2fnTvyERk
XPQeiLf085TNnfnzKvCG8DrQwHZqdb2Wt1M2FJkYvxd14CALK2l6yHUEAyNveUoD
VcmtOapYvCndo9d7JpmoiCzelwKBgQC5lwKqyhwLEvbKegMMSZKrL5AsdKbp/r/r
Vovgw2wBhtmj7Gr8bnkS9nAmI3mETCFolxmbDHPK2Yt6D/+LtIcOdalI6Ox93nHx
IV2kPQff5czWa78IiiGPeJYrp/ZBZK33egWJWPAsQsZAms0GXAQ9vSomUUBi8rUT
Lkc03c13HwKBgQCsSqv0yd5WA6ib3ADHADH7HeTbM2H9T5qW4tdCrtnd3mkCja5r
F0TDhwewQdMMs/+fs97I1hcuvI4Y+KbUjJ9CcMHRzOkcTbFQJFIbOdQf3279cLWl
uIxv+wbxG5XJTm03fjDB3vLvo0Xq6DpGfb5KW2sQ0f3scBN0Q6Upv003mQKBgFPk
oG8Fx6F15BtpBiGyzFsXuAtwe9dAsg6246opjJQwGgfQohgT9CUPQ2jqFk8oft2h
mBCPk3Q53KPDwZesdnSh2XE84VKQkF8Y3xSUBhA+99ZhhExe7IbHUtLPLTEoSr+Y
6BHLI15OnQGtOErMo5oo/XmutvVDk3jlLYkHTo6vAoGBAKaT2qIDOStdCrwRbvD1
SF/pcEytM0rQhiJYmBXKeayUsICTxnSdixb42BSRDTL14F6Jzv2GcGRh80Jx1DVL
6Dmv27MEXx3OnCiHmTCHi3CxqKXhOvJGQCbtLLjluP6pAvQCZ7s3KB6/zS4v/fIv
zygLJrETnjWa1iAMPLnIB9lB
-----END PRIVATE KEY-----`

var DEFAULT_CONFIG_URLS []string = []string{"https://gitlab.com/zhifan999/fq/-/raw/main/config.json", "https://www.githubip.xyz/config.json", "https://d23lye95wfkvbk.cloudfront.net/config.json"}

func Decodejson_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var pemstr string
	var datastr string
	var pemfile string
	var parser *VPNParse = nil
	var cfg *VPNConfig = nil
	var idx int
	err = nil
	if ns == nil {
		return nil
	}
	err = logutil.InitLog(ns)
	if err != nil {
		logutil.Error("can not Initlog err[%s]", err.Error())
		return err
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 1 {
		err = dbgutil.FormatError("need jsonfile ...")
		return
	}

	pemfile = ns.GetString("pemfile")
	if len(pemfile) == 0 {
		pemstr = PEM_DEFAULT
	} else {
		pemstr, err = fileop.ReadFile(pemfile)
		if err != nil {
			return
		}
	}

	parser, err = NewVPNParse(pemstr)
	if err != nil {
		return
	}

	for idx = 0; idx < len(sarr); idx += 1 {
		datastr, err = fileop.ReadFile(sarr[idx])
		if err != nil {
			return
		}
		cfg, err = parser.GetConfig(datastr)
		if err != nil {
			return
		}
		cfg = cfg
	}

	err = nil
	return
}

func Getcfg_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var sarr []string
	var pemstr string
	var parser *VPNParse = nil
	var cfg *VPNConfig = nil
	var urls []string
	var u string
	var pemfile string
	var totalerror error = nil
	var jsons string
	err = nil
	if ns == nil {
		return nil
	}
	err = logutil.InitLog(ns)
	if err != nil {
		logutil.Error("can not Initlog err[%s]", err.Error())
		return err
	}

	pemfile = ns.GetString("pemfile")
	if len(pemfile) == 0 {
		pemstr = PEM_DEFAULT
	} else {
		pemstr, err = fileop.ReadFile(pemfile)
		if err != nil {
			return
		}
	}

	parser, err = NewVPNParse(pemstr)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) == 0 {
		urls = DEFAULT_CONFIG_URLS
	} else {
		urls = sarr
	}

	totalerror = nil
	for _, u = range urls {
		jsons, err = GetHttp(u, 3000, false)
		if err != nil {
			logutil.Error("get [%s] error %s", u, err.Error())
			totalerror = err
			continue
		}
		cfg, err = parser.GetConfig(jsons)
		if err != nil {
			logutil.Error("parse error [%s]\n%s", err.Error(), jsons)
			totalerror = err
			continue
		}

		for _, node := range cfg.Nodes {
			fmt.Printf("server [%s] flag[%s] [%s:%d]\n", node.Name, node.Flag, node.Server, node.Port)
		}
	}
	err = totalerror
	return
}

func Proxy_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx1 interface{}) (err error) {
	var sarr []string
	var localstr string
	var remotestr string
	var lister *tcprelay.RelayListen
	var defremote *tcprelay.DefaultCreateConn
	var ctx context.Context
	var stop context.CancelFunc
	var exited int = 0
	err = nil
	if ns == nil {
		return nil
	}
	err = logutil.InitLog(ns)
	if err != nil {
		logutil.Error("can not Initlog err[%s]", err.Error())
		return err
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = dbgutil.FormatError("need localstr remotestr")
		return
	}
	localstr = sarr[0]
	remotestr = sarr[1]

	defremote, err = tcprelay.NewDefaultCreateConn(remotestr)
	if err != nil {
		return
	}

	lister, err = tcprelay.NewRelayListen(localstr, defremote)
	if err != nil {
		return
	}

	logutil.Debug("listen on [%s]", localstr)

	err = lister.Start()
	if err != nil {
		return
	}
	defer lister.Close()

	ctx, stop = signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for exited == 0 {
		select {
		case <-ctx.Done():
			exited = 1
		case <-time.After(time.Second * 1):
			exited = exited
		}
	}
	err = nil
	return
}

func init() {
	Decodejson_handler(nil, nil, nil)
	Getcfg_handler(nil, nil, nil)
	Proxy_handler(nil, nil, nil)
}
func main() {
	var commandline string
	var err error
	var parser *extargsparse.ExtArgsParse
	var ns *extargsparse.NameSpaceEx

	commandline = `{
		"input|i" : null,
		"output|o" : null,
		"pemfile" : null,
		"decjson<Decodejson_handler>##jsonfile  with decode json for config##" : {
			"$" : "+"
		},
		"getcfg<Getcfg_handler>##url ... to get config##" : {
			"$" : "*"
		},
		"proxy<Proxy_handler>##localstr remotestr to make proxy##" : {
			"$" : 2
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
