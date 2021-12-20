package autoload

import "gitee.com/csingo/ctool/config/vars"

func initRpcClient() {
	for _, service := range vars.RpcConf.ServiceHosts {
		service.Client.RegisterServerHost(service.Host)
	}
}
