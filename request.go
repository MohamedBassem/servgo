package main

import "strings"

var allowedMethods = map[string]bool{
	"GET": true,
}

type Request struct {
	method      string
	path        string
	httpVersion string
	headers     map[string]string
}

func (rq *Request) addHeader(key, val string) {
	rq.headers[key] = val
}

func ParseRequest(requestLines []string) (*Request, error) {
	var request Request
	request.headers = make(map[string]string)
	if len(requestLines) < 1 {
		return nil, &UnparsableRequestError{"Error Parsing the Header Line"}
	}
	headerLineFields := strings.Fields(requestLines[0])
	if len(headerLineFields) < 3 {
		return nil, &UnparsableRequestError{"Error Parsing the Header Line"}
	}
	request.method = headerLineFields[0]
	request.path = headerLineFields[1]
	request.httpVersion = headerLineFields[2]

	if !allowedMethods[request.method] {
		return nil, &NotAllowedMethodError{"Method is not allowed"}
	}

	for i := 1; i < len(requestLines); i++ {
		tmp := strings.Split(requestLines[i], ":")
		if len(tmp) < 2 {
			return nil, &UnparsableRequestError{"Error Parsing the Headers"}
		}
		request.addHeader(strings.TrimSpace(tmp[0]), strings.TrimSpace(tmp[1]))
	}
	return &request, nil
}
