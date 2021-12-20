package autoload

import (
    //TODO:ImportConfig
	"##PROJECT##/config/vars"
	"##PROJECT##/core/cServer"
)

func initConfig() {
    //TODO:InitConfig
	cServer.Inject(vars.Command)
	cServer.Inject(vars.ConfigCenter)
	cServer.Inject(vars.Database)
	cServer.Inject(vars.Redis)
	cServer.Inject(vars.Route)
	cServer.Inject(vars.Server)
	cServer.Inject(vars.Rpc)
	cServer.Inject(vars.Log)
}
