package main

import (
	"github.com/csingo/ctool/autoload"
	"github.com/csingo/ctool/core/cCommand"
	"os"
)

func main() {
	autoload.Init()

	if len(os.Args) <= 1 || os.Args[1] == "help" {
		cCommand.Run("help::doc")
	} else {
		cCommand.Run(os.Args[1], os.Args[2:]...)
	}

}
