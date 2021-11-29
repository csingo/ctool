package main

import (
	"gitee.com/csingo/ctool/autoload"
	"gitee.com/csingo/ctool/core/cCommand"
)

func main() {
	autoload.Init()

	cCommand.Run("test::hello", "--name=cxy")
}
