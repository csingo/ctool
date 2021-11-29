package service

import (
	"context"
	"gitee.com/csingo/ctool/base/app"
	"log"
)

type HelloService struct {
	app.RpcService
}

func (s *HelloService) Say(ctx context.Context, req *app.SayRequest) (*app.SayReply, error) {
	log.Println("say hello")
	log.Printf("%v", req)
	return nil, nil
}
