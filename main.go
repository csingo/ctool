package main

import (
	"gitee.com/csingo/ctool/autoload"
	"gitee.com/csingo/ctool/core/cCommand"
)

func main() {
	autoload.Init()

	//log.Println(os.Getwd())
	cCommand.Run("project::create", "--name=cxy", "--path=.")
}
