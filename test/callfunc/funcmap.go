package main

import (
	"fmt"
	"reflect"
)

type funcMap struct {
	mapfunc map[string]reflect.Value
}

func (self *funcMap) bind(funcname string, fn interface{}) error {
	v := reflect.ValueOf(fn)
	self.mapfunc[funcname] = v
	return nil
}

func newFuncMap() (*funcMap, error) {
	self := &funcMap{}
	self.mapfunc = make(map[string]reflect.Value)
	self.bind("call1", self.call_1)
	self.bind("call2", self.call_2)
	return self, nil
}

func (self *funcMap) call_1(k string) {
	fmt.Printf("call_1 [%s]\n", k)
	return
}

func (self *funcMap) call_2(k string) {
	fmt.Printf("call_2 [%s]\n", k)
	return
}

func (self *funcMap) CallFunc(funcname string, a ...interface{}) []reflect.Value {
	in := make([]reflect.Value, len(a))
	for k, p := range a {
		in[k] = reflect.ValueOf(p)
	}
	return self.mapfunc[funcname].Call(in)
}

func main() {
	p, _ := newFuncMap()
	p.CallFunc("call1", "go")
	p.CallFunc("call2", "cc")
	return
}
