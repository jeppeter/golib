package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type ZZ struct {
	innc int32
}

type CC struct {
	inner   string
	Outer   string
	inner16 uint16
	Outer16 uint16
	innerzz *ZZ
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

func find_interface_field(a interface{}, key string) int {
	p := reflect.ValueOf(a).Elem()
	founded := -1
	maxfld := p.NumField()
	typeofT := p.Type()
	for i := 0; i < maxfld; i++ {
		if typeofT.Field(i).Name == key {
			founded = i
			break
		}
	}
	return founded
}

func set_member(a interface{}, key string, val interface{}) error {
	founded := find_interface_field(a, key)
	if founded < 0 {
		return fmt.Errorf("can not found [%s]", key)
	}

	rs := reflect.ValueOf(a).Elem()
	rf := rs.Field(founded)
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	fmt.Printf("rf type [%s]\n", rf.Type().Name())
	rf.Set(reflect.ValueOf(val))
	return nil
}

func is_member(a interface{}, key string) bool {
	founed := find_interface_field(a, key)
	if founed < 0 {
		return false
	}
	return true
}

func get_value(a interface{}, key string) interface{} {
	founded := find_interface_field(a, key)

	if founded < 0 {
		s := fmt.Sprintf("can not found [%s]", key)
		panic(s)
	}

	rs := reflect.ValueOf(a).Elem()
	rf := rs.Field(founded)
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	return rf.Interface()

}

func main() {
	//var err error
	cc := &CC{inner: "inner", Outer: "outer", inner16: uint16(16), Outer16: uint16(16)}

	set_member(cc, "Outer", "annother Outer")
	set_member(cc, "Outer16", uint16(32))
	set_member(cc, "inner", "annoer inner")
	set_member(cc, "inner16", uint16(32))

	fmt.Printf("Outer [%s]\n", cc.Outer)
	fmt.Printf("Outer16 [%d]\n", cc.Outer16)
	fmt.Printf("inner [%s]\n", get_value(cc, "inner").(string))
	fmt.Printf("inner16 [%d]\n", get_value(cc, "inner16").(uint16))
	return
}
