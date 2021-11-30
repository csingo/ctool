package app

import (
	"context"
)

type RpcService struct {
	##SERVICE##Server
}

func (s *RpcService) RpcServiceName() (string, string) {
	return "##APP##", "##SERVICE##"
}

type ##SERVICE##HttpClient struct {
	host string
}

func (s *##SERVICE##HttpClient) RegisterServerHost(host string) {
	s.host = host
}

var ##SERVICE## = ##SERVICE##HttpClient{}
