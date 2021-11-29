package main

import (
	"framework/autoload"
	"framework/core/cServer"
)

func main() {
	autoload.Init()
	cServer.Init()
	
	cServer.Start()
}
