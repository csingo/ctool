package autoload

import (
	AppService "gitee.com/csingo/ctool/app/service"
	"gitee.com/csingo/ctool/core/cServer"
)

func initService() {
	cServer.Inject(&AppService.HelloService{})
}
