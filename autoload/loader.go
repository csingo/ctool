package autoload

import "gitee.com/csingo/ctool/autoload/qdLog"

func Init() {
	qdLog.InitLog()
	initConfig()
	initCommand()
	initController()
	initMiddleware()
	initService()
	initRpcClient()

}
