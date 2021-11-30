package main

import (
	"gitee.com/csingo/ctool/autoload"
	"gitee.com/csingo/ctool/core/cCommand"
)

func main() {
	autoload.Init()

	//args := os.Args
	//
	//cCommand.Run(args[1], args[2:]...)

	cCommand.Run("tool::init")
}
