package main

import (
	"framework/autoload"
	"framework/core/cCommand"
)

func main() {
	autoload.Init()

	cCommand.Run("test::hello", "--name=cxy")
}
