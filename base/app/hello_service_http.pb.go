package app

import (
	"context"
)

type RpcService struct {
	HelloServiceServer
}

func (s *RpcService) RpcServiceName() (string, string) {
	return "app", "HelloService"
}

type HelloServiceHttpClient struct {
	host string
}

func (s *HelloServiceHttpClient) RegisterServerHost(host string) {
	s.host = host
}

var HelloService = &HelloServiceHttpClient{}

func (s *HelloServiceHttpClient) Say(ctx context.Context, req *SayRequest) (rsp *SayReply, err error) {
	err = call(ctx, s.host, "app", "HelloService", "Say", req, rsp)
	if err != nil {
		return
	}

	return
}
