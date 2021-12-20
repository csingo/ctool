package typs

import "github.com/csingo/ctool/core/cRpc"

type RpcConf struct {
	ServiceHosts []*RpcConf_ServieHost
}

type RpcConf_ServieHost struct {
	Client cRpc.RpcClientInterface
	Host   string
}

func (i *RpcConf) ConfigName() string {
	return "RpcConf"
}
