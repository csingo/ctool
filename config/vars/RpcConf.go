package vars

import (
	"github.com/csingo/ctool/base/app"
	"github.com/csingo/ctool/config/typs"
)

var RpcConf = &typs.RpcConf{
	ServiceHosts: []*typs.RpcConf_ServieHost{
		{app.HelloService, "http://127.0.0.1:8080"},
	},
}
