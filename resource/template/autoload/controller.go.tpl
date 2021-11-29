package autoload

import (
	AppController "##PROJECT##/app/controller"
	"##PROJECT##/core/cServer"
)

func initController() {
	cServer.Inject(&AppController.HomeController{})
}
