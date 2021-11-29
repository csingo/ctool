package autoload

import (
	AppMiddleware "gitee.com/csingo/ctool/app/middleware"
	"gitee.com/csingo/ctool/core/cServer"
)

func initMiddleware() {
	cServer.Inject(&AppMiddleware.TestMiddleware{})
}
