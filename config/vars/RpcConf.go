package vars

import (
	"gitee.com/csingo/ctool/base/app"
	"gitee.com/csingo/ctool/config/typs"
)

var RpcConf = &typs.RpcConf{
	ServiceHosts: []*typs.RpcConf_ServieHost{
		{app.HelloService, "http://127.0.0.1:8080"},
	},
}
