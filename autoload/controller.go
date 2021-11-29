package autoload

import (
	AppController "framework/app/controller"
	"framework/core/cServer"
)

func initController() {
	cServer.Inject(&AppController.HomeController{})
}
