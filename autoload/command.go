package autoload

import (
	"github.com/csingo/ctool/app/app/command"
	"github.com/csingo/ctool/core/cServer"
)

func initCommand() {
	cServer.Inject(&command.ProjectCommand{})
	cServer.Inject(&command.AppCommand{})
	cServer.Inject(&command.ToolCommand{})
	cServer.Inject(&command.ConfigCommand{})
	cServer.Inject(&command.HelpCommand{})
	cServer.Inject(&command.SdkCommand{})
}
