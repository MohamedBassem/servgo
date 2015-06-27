package main

import (
	"bufio"
	"fmt"
	"net"
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
		fmt.Printf("A new request is handled by worker %v\n", id)
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
		}
		cl.Write(response.toByteSlice())
		cl.Close()
	}
}
