package autoload

import (
	AppCommand "gitee.com/csingo/ctool/app/command"
	"gitee.com/csingo/ctool/core/cServer"
)

func initCommand() {
	cServer.Inject(&AppCommand.ProjectCommand{})
	cServer.Inject(&AppCommand.AppCommand{})
	cServer.Inject(&AppCommand.ToolCommand{})
	cServer.Inject(&AppCommand.ConfigCommand{})
	cServer.Inject(&AppCommand.HelpCommand{})
}
