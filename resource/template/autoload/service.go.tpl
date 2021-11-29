package autoload

import (
	AppService "##PROJECT##/app/service"
	"##PROJECT##/core/cServer"
)

func initService() {
	cServer.Inject(&AppService.HelloService{})
}
