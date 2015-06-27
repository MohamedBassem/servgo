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
	var maxQueuedConnecions = flag.Int("max-queued", 1024, "The maximum number of connections that can be queued in the server")
	var numberOfWorkers = flag.Int("num-workers", 2, "Number of workers serving the requests")
	flag.Parse()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Server is listening to %v\n", *addr)

	var requestQueue = make(chan net.Conn, *maxQueuedConnecions)
	for i := 0; i < *numberOfWorkers; i++ {
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
