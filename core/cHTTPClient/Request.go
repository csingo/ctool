package cHTTPClient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HTTPMethod string

const (
	MethodHEAD    HTTPMethod = "HEAD"
	MethodGET     HTTPMethod = "GET"
	MethodPOST    HTTPMethod = "POST"
	MethodOPTIONS HTTPMethod = "OPTIONS"
	MethodDELETE  HTTPMethod = "DELETE"
	MethodPUT     HTTPMethod = "PUT"
	MethodPATCH   HTTPMethod = "PATCH"
	MethodTRACE   HTTPMethod = "TRACE"
	MethodCONNECT HTTPMethod = "CONNECT"
)

type Option struct {
	Host    string
	Uri     string
	Method  HTTPMethod
	Headers map[string]string
	Query   map[string]string
	Data    string
}

func Request(option Option) (header http.Header, body []byte, err error) {
	values := url.Values{}
	for k, v := range option.Query {
		values.Set(k, v)
	}

	host := strings.TrimRight(option.Host, "/")
	uri := strings.TrimPrefix(option.Uri, "/")
	query := values.Encode()
	var url = host
	if uri != "" {
		url = fmt.Sprintf("%s/%s", url, uri)
	}
	if query != "" {
		url = fmt.Sprintf("%s?%s", url, query)
	}

	var req *http.Request
	client := &http.Client{}

	req, err = http.NewRequest(string(option.Method), url, strings.NewReader(option.Data))
	if err != nil {
		return
	}

	switch option.Method {
	case MethodPOST:
		if option.Headers == nil {
			break
		}
		for k, v := range option.Headers {
			req.Header.Set(k, v)
		}
	}

	rsp, err := client.Do(req)
	header = rsp.Header
	defer rsp.Body.Close()
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}

	return
}
