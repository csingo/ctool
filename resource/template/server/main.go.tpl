package main

import (
	"##PROJECT##/autoload"
	"##PROJECT##/core/cServer"
)

func main() {
	autoload.Init()
	cServer.Init()
	
	cServer.Start()
}
