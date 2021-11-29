package main

import (
	"##PROJECT##/autoload"
	"##PROJECT##/core/cCommand"
)

const (
	VERSION = "v0.0.2"
)

func main() {
	autoload.Init()

	cCommand.Run("test::hello", "--name=cxy")
}
