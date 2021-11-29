package autoload

import (
	AppMiddleware "framework/app/middleware"
	"framework/core/cServer"
)

func initMiddleware() {
	cServer.Inject(&AppMiddleware.TestMiddleware{})
}
