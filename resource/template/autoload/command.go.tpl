package autoload

import (
	AppCommand "##PROJECT##/app/command"
	"##PROJECT##/core/cServer"
)

func initCommand() {
	cServer.Inject(&AppCommand.TestCommand{})
}
