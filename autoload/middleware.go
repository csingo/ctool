package autoload

import (
	AppMiddleware "github.com/csingo/ctool/app/app/middleware"
	"github.com/csingo/ctool/core/cServer"
)

func initMiddleware() {
	cServer.Inject(&AppMiddleware.TestMiddleware{})
}
