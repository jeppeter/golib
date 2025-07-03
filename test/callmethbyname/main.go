package main

import (
	"fmt"
	"os"
	"reflect"
)

type CN struct {
	val  int
	name string
	maps map[string]string
	arr  []string
}

func NewCN() (retv *CN) {
	retv = &CN{}
	retv.val = 0
	retv.name = ""
	retv.maps = make(map[string]string)
	retv.arr = []string{}
	return
}

func (retv *CN) SetValue(val int) (oldval int) {
	oldval = retv.val
	retv.val = val
	return
}

func (retv *CN) AddArray(n string) (retn []string) {
	retn = retv.arr
	retv.arr = append(retv.arr, n)
	return
}

func (retv *CN) NN(val int) (oldval int) {
	oldval = retv.val
	retv.val = val
	return
}

func (retv *CN) SetName(val string) (oldval string) {
	oldval = retv.name
	retv.name = val
	return
}

func (retv *CN) SetMap(k, v string) (oldval string) {
	var ok bool
	oldval, ok = retv.maps[k]
	if !ok {
		oldval = ""
	}
	retv.maps[k] = v
	return
}

func (retv *CN) GetMapValue(k string) (retn string, err error) {
	var ok bool
	retn, ok = retv.maps[k]
	if !ok {
		err = fmt.Errorf("no [%s] key", k)
		return
	}
	err = nil
	return
}

func call_method(name string, ptr *CN, a ...interface{}) (retv []reflect.Value, err error) {
	var curval reflect.Value
	var methval reflect.Value
	var args []reflect.Value
	var i int
	curval = reflect.ValueOf(ptr)
	methval = curval.MethodByName(name)
	retv = []reflect.Value{}
	if !methval.IsValid() {
		err = fmt.Errorf("can not find [%s]", name)
		return
	}

	args = []reflect.Value{}
	for i = 0; i < len(a); i += 1 {
		args = append(args, reflect.ValueOf(a[i]))
	}

	retv = methval.Call(args)
	err = nil
	return
}

func main() {
	var c *CN
	var args []reflect.Value
	var err error
	c = NewCN()
	args, err = call_method("SetValue", c, 77)
	if err != nil {
		fmt.Fprintf(os.Stderr, "get SetValue error %s\n", err.Error())
	} else {
		if len(args) != 1 || !args[0].CanInt() {
			fmt.Fprintf(os.Stderr, "SetValue retargs [%d]\n", len(args))
		} else {
			fmt.Fprintf(os.Stdout, "SetValue retval %d\n", args[0].Int())
		}
	}

	args, err = call_method("SetValue", c, 99)
	if err != nil {
		fmt.Fprintf(os.Stderr, "get SetValue error %s\n", err.Error())
	} else {
		if len(args) != 1 || !args[0].CanInt() {
			fmt.Fprintf(os.Stderr, "SetValue retargs [%d]\n", len(args))
		} else {
			fmt.Fprintf(os.Stdout, "SetValue retval %d\n", args[0].Int())
		}
	}

	args, err = call_method("NN", c, 88)
	if err != nil {
		fmt.Fprintf(os.Stderr, "get nn error %s\n", err.Error())
	} else {
		if len(args) != 1 || !args[0].CanInt() {
			fmt.Fprintf(os.Stderr, "nn retargs [%d]\n", len(args))
		} else {
			fmt.Fprintf(os.Stdout, "nn retval %d\n", args[0].Int())
		}
	}

	args, err = call_method("SetMap", c, "99k", "002")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get SetMap error %s\n", err.Error())
	} else {
		if len(args) != 1 || args[0].Kind() != reflect.String {
			fmt.Fprintf(os.Stderr, "SetMap retargs [%d]\n", len(args))
		} else {
			fmt.Fprintf(os.Stdout, "SetMap retval [%s]\n", args[0].String())
		}
	}

	args, err = call_method("SetMap", c, "99k", "087")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get SetMap error %s\n", err.Error())
	} else {
		if len(args) != 1 || args[0].Kind() != reflect.String {
			fmt.Fprintf(os.Stderr, "SetMap retargs [%d]\n", len(args))
		} else {
			fmt.Fprintf(os.Stdout, "SetMap retval [%s]\n", args[0].String())
		}
	}

	args, err = call_method("AddArray", c, "77k")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get AddArray error %s\n", err.Error())
	} else {
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "AddArray retargs [%d] \n", len(args))
		} else if args[0].Kind() != reflect.Array && args[0].Kind() != reflect.Slice {
			fmt.Fprintf(os.Stderr, "AddArray kind %s", args[0].Kind())
		} else {
			v := args[0].Interface()
			fmt.Fprintf(os.Stdout, "AddArray retval %v\n", v.([]string))
		}
	}

	args, err = call_method("AddArray", c, "99k")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get AddArray error %s\n", err.Error())
	} else {
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "AddArray retargs [%d] \n", len(args))
		} else if args[0].Kind() != reflect.Array && args[0].Kind() != reflect.Slice {
			fmt.Fprintf(os.Stderr, "AddArray kind %s", args[0].Kind())
		} else {
			v := args[0].Interface()
			fmt.Fprintf(os.Stdout, "AddArray retval %v\n", v.([]string))
		}
	}

	args, err = call_method("GetMapValue", c, "99k")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get GetMapValue error %s\n", err.Error())
	} else {
		if len(args) != 2 || args[0].Kind() != reflect.String {
			fmt.Fprintf(os.Stderr, "GetMapValue retargs [%d]\n", len(args))
		} else {
			ival := args[1].Interface()
			if ival == nil {
				fmt.Fprintf(os.Stdout, "GetMapValue value [%s]\n", args[0].String())
			} else {
				switch ival.(type) {
				case error:
					errval := ival.(error)
					fmt.Fprintf(os.Stderr, "GetMapValue error %s\n", errval.Error())
				default:
					fmt.Fprintf(os.Stderr, "GetMapValue args[1] value %v", ival)
				}
			}
		}
	}

	args, err = call_method("GetMapValue", c, "993k")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get GetMapValue error %s\n", err.Error())
	} else {
		if len(args) != 2 || args[0].Kind() != reflect.String {
			fmt.Fprintf(os.Stderr, "GetMapValue retargs [%d]\n", len(args))
		} else {
			ival := args[1].Interface()
			if ival == nil {
				fmt.Fprintf(os.Stdout, "GetMapValue value [%s]\n", args[0].String())
			} else {
				switch ival.(type) {
				case error:
					errval := ival.(error)
					fmt.Fprintf(os.Stderr, "GetMapValue error %s\n", errval.Error())
				default:
					fmt.Fprintf(os.Stderr, "GetMapValue args[1] value %v", ival)
				}
			}
		}
	}
}
