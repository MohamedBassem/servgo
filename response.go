package main

import (
	"strconv"
	"time"
)

var statusReason = map[int]string{
	200: "OK",
	400: "Bad Request",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
}

type Response struct {
	httpVersion string
	status      int
	headers     map[string]string
	body        string
}

func NewResponse() Response {
	var res Response
	res.headers = make(map[string]string)
	res.status = 200
	return res
}

func (res *Response) addServerHeaders(httpVersion string) {
	res.headers["Date"] = time.Now().String()
	res.headers["Content-Length"] = strconv.Itoa(len([]byte(res.body)))
	res.headers["Server"] = "GoServe"
	res.httpVersion = httpVersion
}

func (res *Response) setStatusCode(code int) {
	res.status = code
}

func (res *Response) addHeader(key, val string) {
	res.headers[key] = val
}

func (res *Response) setBody(body string) {
	res.body = body
}

func (res *Response) toByteSlice() []byte {
	var ret string
	ret += res.httpVersion + " " + strconv.Itoa(res.status) + " " + statusReason[res.status] + "\r\n"
	for k, v := range res.headers {
		ret += k + ": " + v + "\r\n"
	}
	ret += "\r\n"
	ret += res.body
	return []byte(ret)
}
