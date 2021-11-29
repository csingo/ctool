package cRpc

type RpcClientInterface interface {
	RegisterServerHost(host string)
}

type RpcServiceInterface interface {
	RpcServiceName() (string, string)
}

type RpcService struct{}

func (s *RpcService) RpcServiceName() (string, string) {
	return "cRpc", "RpcServer"
}

type container struct {
	Services map[string]RpcServiceInterface
}

var RpcServiceContainer = &container{
	Services: map[string]RpcServiceInterface{},
}
