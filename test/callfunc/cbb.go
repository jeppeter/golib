package main

import (
	"fmt"
	"os"
)

type ExtArgsOptions struct {
}

type ExtArgsParse struct {
	loadPriority []int
}

var parser_priority_args = []int{1, 2, 3, 4}

func NewExtArgsParse(options *ExtArgsOptions, priority interface{}) (self *ExtArgsParse, err error) {
	var pr []int
	if priority == nil {
		pr = parser_priority_args
	} else {
		switch priority.(type) {
		case []int:
			pr = make([]int, 0)
			for _, iv := range priority.([]int) {
				pr = append(pr, iv)
			}
		default:
			return nil, fmt.Errorf("unknown type [%v]", priority)
		}
	}
	err = nil
	self = &ExtArgsParse{}
	self.loadPriority = pr
	return
}

func main() {
	var p *ExtArgsParse
	var e error
	p, e = NewExtArgsParse(nil, []int{1, 2, 3})
	if e != nil {
		fmt.Fprintf(os.Stderr, "1 %s\n", e.Error())
		return
	}
	fmt.Fprintf(os.Stdout, "p [%v]\n", p)
	p, e = NewExtArgsParse(nil, nil)
	if e != nil {
		fmt.Fprintf(os.Stderr, "2 %s\n", e.Error())
		return
	}
	fmt.Fprintf(os.Stdout, "p [%v]\n", p)

	p, e = NewExtArgsParse(nil, []string{"hello", "bb"})
	if e != nil {
		fmt.Fprintf(os.Stderr, "2 %s\n", e.Error())
		return
	}
	fmt.Fprintf(os.Stdout, "p [%v]\n", p)

	fmt.Fprintf(os.Stdout, "succ\n")
}
