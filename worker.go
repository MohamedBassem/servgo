package servgo

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func worker(id int, server *Server, requestQueue <-chan net.Conn) {
	for {
		cl := <-requestQueue
		response := NewResponse()

		reader := bufio.NewReader(cl)
		requestLines := make([]string, 0)
		for {
			buffer, _, _ := reader.ReadLine()
			if string(buffer) == "" {
				break
			}
			requestLines = append(requestLines, string(buffer))
		}

		request, err := parseRequest(requestLines)
		if err != nil {
			response = NewErrorResponse(400, err.Error())
			if request == nil {
				response.addServerHeaders("HTTP/1.1")
			} else {
				response.addServerHeaders(request.httpVersion)
			}
			if len(requestLines) > 0 {
				fmt.Printf("Worker %v: %v [%v] : '%v' %v %v\n", id, cl.RemoteAddr(), time.Now().String(), requestLines[0], response.status, response.headers["Content-Length"])
			} else {
				fmt.Printf("Worker %v: %v [%v]: '%v' %v %v\n", id, cl.RemoteAddr(), time.Now().String(), "Error Parsing Request", response.status, response.headers["Content-Length"])
			}
		} else {
			switch request.method {
			case "GET":
				if server.GetHandler() == nil {
					response = NewErrorResponse(405, "Method Not Allowed")
				} else {
					response = server.GetHandler()(*request)
				}
				response.addServerHeaders(request.httpVersion)
			case "HEAD":
				if server.GetHandler() == nil {
					response = NewErrorResponse(405, "Method Not Allowed")
				} else {
					response = server.GetHandler()(*request)
					response.body = ""
				}
				response.addServerHeaders(request.httpVersion)
			default:
				response = NewErrorResponse(405, "Method Not Allowed")
			}
			fmt.Printf("Worker %v: %v [%v] : '%v %v' %v %v\n", id, cl.RemoteAddr(), time.Now().String(), request.method, request.path, response.status, response.headers["Content-Length"])
		}
		cl.Write(response.toByteSlice())
		cl.Close()
	}
}
