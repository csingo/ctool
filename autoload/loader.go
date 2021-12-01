package autoload

func Init() {
	initConfig()
	initCommand()
	initController()
	initMiddleware()
	initService()
	initRpcClient()
}
