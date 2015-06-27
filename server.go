package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

type Server struct {
	addr    string
	rootDir string
}

func isValidRootDir(dir string) bool {
	fi, err := os.Stat(dir)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return fi.IsDir()
}

func main() {
	var addr = flag.String("addr", ":8080", "The server's port")
	var maxQueuedConnecions = flag.Int("max-queued", 1024, "[Optional] The maximum number of connections that can be queued in the server")
	var numberOfWorkers = flag.Int("num-workers", 2, "[Optional] Number of workers serving the requests")
	var rootDir = flag.String("root-dir", "nil", "[Required] The root dir for serving the files")
	flag.Parse()

	if *rootDir == "nil" || !isValidRootDir(*rootDir) {
		fmt.Println("A --root-dir must be given and valid.")
		os.Exit(1)
	}

	if string(*rootDir)[len(*rootDir)-1:] != "/" {
		*rootDir += "/"
	}

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Println(err)
	}
	server := Server{*addr, *rootDir}
	fmt.Printf("Server is listening to %v\n", *addr)

	var requestQueue = make(chan net.Conn, *maxQueuedConnecions)
	for i := 0; i < *numberOfWorkers; i++ {
		go worker(i, &server, requestQueue)
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
