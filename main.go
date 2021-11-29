package main

import (
	"gitee.com/csingo/ctool/autoload"
	"log"
	"os"
)

func main() {
	autoload.Init()

	log.Println(os.Getwd())
	//cCommand.Run("project::create", "--name=cxy", "--path=.")
}
