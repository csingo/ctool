package autoload

import "framework/config/vars"

func initRpcClient() {
	for _, service := range vars.RpcConf.ServiceHosts {
		service.Client.RegisterServerHost(service.Host)
	}
}
