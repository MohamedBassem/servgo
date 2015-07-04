package servgo

import (
	"fmt"
	"net"
)

type Server struct {
	addr                  string
	numWorkers, maxQueued int
	getHandler            func(Request) Response
}

func NewServer(addr string, numWorkers, maxQueued int) Server {
	server := Server{
		addr:       addr,
		numWorkers: numWorkers,
		maxQueued:  maxQueued,
	}
	return server
}

func (s *Server) SetGetHandler(f func(Request) Response) {
	s.getHandler = f
}

func (s *Server) GetHandler() func(Request) Response {
	return s.getHandler
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	var requestQueue = make(chan net.Conn, s.maxQueued)
	for i := 0; i < s.numWorkers; i++ {
		go worker(i, s, requestQueue)
	}

	for {
		cl, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		requestQueue <- cl
	}
	return nil
}
