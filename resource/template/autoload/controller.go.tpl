package autoload

import (
	AppController "gitee.com/csingo/ctool/app/controller"
	"gitee.com/csingo/ctool/core/cServer"
)

func initController() {
	cServer.Inject(&AppController.HomeController{})
}
