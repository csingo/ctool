package autoload

import (
	AppController "github.com/csingo/ctool/app/app/controller"
	"github.com/csingo/ctool/core/cServer"
)

func initController() {
	cServer.Inject(&AppController.HomeController{})
}
