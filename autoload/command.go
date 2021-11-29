package autoload

import (
	AppCommand "framework/app/command"
	"framework/core/cServer"
)

func initCommand() {
	cServer.Inject(&AppCommand.TestCommand{})
}
