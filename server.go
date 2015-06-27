package main

import (
	"flag"
	"fmt"
	"net"
)

type Server struct {
	port int
}

func main() {
	var addr = flag.String("addr", ":8080", "The server's port")

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Println(err)
	}

	var requestQueue = make(chan net.Conn, 1024)
	for i := 0; i < 8; i++ {
		go worker(i, requestQueue)
	}

	for {
		cl, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		requestQueue <- cl
	}
}
