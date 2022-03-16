package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"net/http"
)

func init() {
	Server_handler(nil, nil, nil)
}

type HttpSvr struct {
	content string
	refsvr  *http.Server
}

func (ht *HttpSvr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Debug("path %s", r.URL.Path)
	fmt.Fprintf(w, ht.content)
}

func Server_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var ht *HttpSvr
	var svr *http.Server
	var ls string
	if ns == nil {
		return nil
	}
	err = InitLog(ns)
	if err != nil {
		Error("can not Initlog err[%s]", err.Error())
		return err
	}
	ls = fmt.Sprintf(":%s", ns.GetArray("subnargs")[0])
	svr = &http.Server{
		Addr: ls,
	}

	ht = &HttpSvr{content: "hello", refsvr: svr}

	mux := http.NewServeMux()
	mux.Handle("/", ht)
	svr.Handler = mux

	Debug("listen on %s", ls)
	svr.ListenAndServe()
	err = nil
	return

}
