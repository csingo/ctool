package main

import (
	"github.com/csingo/ctool/autoload"
	"github.com/csingo/ctool/core/cServer"
)

func main() {
	autoload.Init()
	cServer.Init()
	
	cServer.Start()
}
