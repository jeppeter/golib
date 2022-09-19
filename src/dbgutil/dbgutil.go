package dbgutil

import (
	"fmt"
	"reflect"
	"runtime"
)

func format_out_stack(level int) string {
	_, f, l, _ := runtime.Caller(level)
	return fmt.Sprintf("[%s:%d]", f, l)
}

func format_out_string_total(level int, fmtstr string, a ...interface{}) string {
	outstr := format_out_stack((level + 1))
	outstr += fmt.Sprintf(fmtstr, a...)
	return outstr
}

const (
	def_stacklevel_added = 3
)

func FormatError(a ...interface{}) (err error) {
	var stacklevel int = def_stacklevel_added
	var vaargs []interface{}
	var fmtstr string
	var ct string

	if len(a) > 0 {
		switch v := a[0].(type) {
		case int:
			stacklevel = a[0].(int)
			stacklevel += def_stacklevel_added
			if len(a) > 2 {
				vaargs = a[2:]
			}
			if len(a) > 1 {
				ct = reflect.TypeOf(a[1]).Name()
				if ct == "string" {
					fmtstr = a[1].(string)
				} else {
					fmtstr = "unknown type string"
				}
			}
		case string:
			if len(a) > 1 {
				vaargs = a[1:]
			}
			if len(a) > 0 {
				fmtstr = a[0].(string)
			}
		default:
			fmtstr = fmt.Sprintf("unknown type [%s]", v)
		}
	}
	outstr := format_out_stack(stacklevel)
	if len(vaargs) == 0 {
		outstr += fmt.Sprintf(fmtstr)
	} else {
		outstr += fmt.Sprintf(fmtstr, vaargs...)
	}
	return fmt.Errorf("%s", outstr)
}
