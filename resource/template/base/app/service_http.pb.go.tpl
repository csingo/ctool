package ##APP##

import (
	"context"
)

type Rpc##SERVICE## struct {
	##SERVICE##Server
}

func (s *Rpc##SERVICE##) RpcServiceName() (string, string) {
	return "##APP##", "##SERVICE##"
}

type ##SERVICE##HttpClient struct {
	host string
}

func (s *##SERVICE##HttpClient) RegisterServerHost(host string) {
	s.host = host
}

var ##SERVICE## = &##SERVICE##HttpClient{}
