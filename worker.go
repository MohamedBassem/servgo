package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func handleGetRequest(req *Request) Response {
	res := NewResponse()
	res.setBody("Hello World")
	res.setStatusSuccess()
	return res
}

func handleErrorResponse(err error) Response {
	res := NewResponse()
	res.setBody(err.Error())
	switch err.(type) {
	case *UnparsableRequestError:
		res.setStatusBadRequest()
	case *NotAllowedMethodError:
		res.setStatusNotAllowedMethod()
	case *NotFoundError:
		res.setStatusNotFound()
	}
	return res
}

func worker(id int, requestQueue <-chan net.Conn) {
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
				response = handleGetRequest(request)
				response.addServerHeaders(request.httpVersion)
			case "HEAD":
				response = handleGetRequest(request)
				response.addServerHeaders(request.httpVersion)
				response.body = ""
			}
			fmt.Printf("Worker %v: %v [%v] : '%v %v' %v %v\n", id, cl.RemoteAddr(), time.Now().String(), request.method, request.path, response.status, response.headers["Content-Length"])
		}
		cl.Write(response.toByteSlice())
		cl.Close()
	}
}
