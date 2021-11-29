package main

import (
	"gitee.com/csingo/ctool/autoload"
	"gitee.com/csingo/ctool/core/cCommand"
)

const (
	VERSION = "v0.0.2"
)

func main() {
	autoload.Init()

	cCommand.Run("test::hello", "--name=cxy")
}
