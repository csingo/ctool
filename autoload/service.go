package autoload

import (
	AppService "framework/app/service"
	"framework/core/cServer"
)

func initService() {
	cServer.Inject(&AppService.HelloService{})
}
