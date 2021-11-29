package main

import (
	"gitee.com/csingo/ctool/autoload"
	"gitee.com/csingo/ctool/core/cServer"
)

func main() {
	autoload.Init()
	cServer.Init()
	
	cServer.Start()
}
