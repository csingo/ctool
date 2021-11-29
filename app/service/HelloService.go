package service

import (
	"context"
	"framework/base/app"
	"log"
)

type HelloService struct {
	app.RpcService
}

func (s *HelloService) Say(ctx context.Context, req *app.SayRequest) (*app.SayReply, error) {
	log.Println("say hello")
	return nil, nil
}
