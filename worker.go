package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"strings"
	"time"
)

func handleGetRequest(req *Request, server *Server) Response {
	res := NewResponse()
	path := req.path[1:]
	if path == "" {
		path = "index.html"
	}
	fullPath := server.rootDir + path
	absPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absPath, server.rootDir) {
		return handleErrorResponse(&ForbiddenError{"403 Forbidden"})
	}

	f, err := ioutil.ReadFile(server.rootDir + path)
	if err != nil {
		return handleErrorResponse(&NotFoundError{"File not found"})
	}
	res.setBody(string(f))
	res.setStatusCode(200)
	return res
}

func handleErrorResponse(err error) Response {
	res := NewResponse()
	res.setBody(err.Error())
	switch err.(type) {
	case *UnparsableRequestError:
		res.setStatusCode(400)
	case *NotAllowedMethodError:
		res.setStatusCode(405)
	case *NotFoundError:
		res.setStatusCode(404)
	case *ForbiddenError:
		res.setStatusCode(403)
	}
	return res
}

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

		request, err := ParseRequest(requestLines)
		if err != nil {
			response = handleErrorResponse(err)
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
				response = handleGetRequest(request, server)
				response.addServerHeaders(request.httpVersion)
			case "HEAD":
				response = handleGetRequest(request, server)
				response.addServerHeaders(request.httpVersion)
				response.body = ""
			}
			fmt.Printf("Worker %v: %v [%v] : '%v %v' %v %v\n", id, cl.RemoteAddr(), time.Now().String(), request.method, request.path, response.status, response.headers["Content-Length"])
		}
		cl.Write(response.toByteSlice())
		cl.Close()
	}
}
