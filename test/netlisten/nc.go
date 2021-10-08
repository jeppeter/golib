package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	var ln net.Listener
	var err error
	var liststr string = ":40008"
	var conn net.Conn
	if len(os.Args) > 1 {
		liststr = os.Args[1]
	}
	fmt.Printf("listen on %s\n", liststr)
	ln, err = net.Listen("tcp", liststr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not listen on [%s] [%s]\n", liststr, err.Error())
		os.Exit(4)
	}
	for {
		conn, err = ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "accept [%s] error[%s]\n", liststr, err.Error())
			os.Exit(4)
		}
		fmt.Fprintf(os.Stdout, "accept [%s]\n", liststr)
		conn.Close()
	}
}
