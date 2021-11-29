package autoload

import (
	AppMiddleware "##PROJECT##/app/middleware"
	"##PROJECT##/core/cServer"
)

func initMiddleware() {
	cServer.Inject(&AppMiddleware.TestMiddleware{})
}
