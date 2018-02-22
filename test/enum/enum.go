package main

import (
	"fmt"
	"reflect"
)

type CC struct {
	inner   string
	Outer   string
	inner16 uint16
	Outer16 uint16
}

func list_member(a interface{}) []string {
	names := make([]string, 0)
	s := reflect.ValueOf(a).Elem()
	typeofT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		cur := typeofT.Field(i)
		names = append(names, cur.Name)
	}

	return names
}

func set_member(a interface{}, key string, val interface{}) error {
	names := list_member(a)
	founded := false
	for _, n := range names {
		if n == key {
			founded = true
			break
		}
	}

	if !founded {
		return fmt.Errorf("can not found [%s]", key)
	}

	s := reflect.ValueOf(a).Elem()
	sv := s.FieldByName(key)
	sv.Set(reflect.ValueOf(val))
	return nil
}

func main() {
	cc := &CC{inner: "inner", Outer: "outer", inner16: uint16(16), Outer16: uint16(16)}
	names := list_member(cc)
	for _, n := range names {
		fmt.Printf("%s\n", n)
	}

	set_member(cc, "Outer", "annother Outer")
	set_member(cc, "Outer16", uint16(32))

	fmt.Printf("inner [%s]\n", cc.Outer)
	fmt.Printf("inner16 [%d]\n", cc.Outer16)
	return
}
