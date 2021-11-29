package vars

import (
	"framework/base/app"
	"framework/config/typs"
)

var RpcConf = &typs.RpcConf{
	ServiceHosts: []*typs.RpcConf_ServieHost{
		{app.HelloService, "http://127.0.0.1:8080"},
	},
}
