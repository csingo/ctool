package app

import (
	"context"
	"encoding/json"
	"gitee.com/csingo/ctool/core/cHTTPClient"
)

func call(ctx context.Context, host string, req interface{}, rsp interface{}) (err error) {
	data, err := json.Marshal(req)
	if err != nil {
		return
	}

	var body []byte
	_, body, err = cHTTPClient.Request(cHTTPClient.Option{
		Host:   host,
		Uri:    "/rpc/call",
		Method: cHTTPClient.MethodPOST,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Rpc-App":      "app",
			"Rpc-Service":  "HelloService",
			"Rpc-Method":   "Say",
		},
		Query: nil,
		Data:  string(data),
	})
	if err != nil {
		return
	}

	err = json.Unmarshal(body, rsp)
	if err != nil {
		return
	}

	return
}
