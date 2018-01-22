package main

import (
	"fmt"
	"os"
)

func TestA(a ...interface{}) error {
	var icnt int
	breakout := false
	for icnt = 0; icnt < len(a); icnt += 1 {
		switch v := a[icnt].(type) {
		case int:
			fmt.Fprintf(os.Stdout, "[%d] int [%v]\n", icnt, a[icnt])
		case string:
			breakout = true
		default:
			fmt.Fprintf(os.Stdout, "[%d] %s [%v]\n", icnt, v, a[icnt])
		}
		if breakout {
			break
		}
	}
	if icnt < (len(a) - 1) {
		fmt.Fprintf(os.Stdout, a[icnt].(string), a[(icnt+1):]...)
	} else {
		fmt.Fprintf(os.Stdout, a[icnt].(string))
	}

	return nil
}

func main() {
	TestA(10, "newok [%d] [%d]\n", 20, 300)
	TestA(20, "hello world\n")
}
