package autoload

import (
	"##PROJECT##/config/vars"
	"##PROJECT##/core/cServer"
)

func initConfig() {
	cServer.Inject(vars.Command)
	cServer.Inject(vars.ConfigCenter)
	cServer.Inject(vars.Database)
	cServer.Inject(vars.Redis)
	cServer.Inject(vars.Route)
	cServer.Inject(vars.Server)
	cServer.Inject(vars.RpcConf)
}
