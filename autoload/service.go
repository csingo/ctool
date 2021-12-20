package autoload

import (
	AppService "github.com/csingo/ctool/app/app/service"
	"github.com/csingo/ctool/core/cServer"
)

func initService() {
	cServer.Inject(&AppService.HelloService{})
}
