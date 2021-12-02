package command

import (
	"fmt"
	"gitee.com/csingo/ctool/core/cCommand"
	"github.com/fatih/color"
)

type HelpCommand struct{}

type commandHelpDoc struct {
	name        string
	method      string
	desc        string
	options     []string
	optionsDesc map[string]string
}

func (i *HelpCommand) Help() *cCommand.CommandHelpDoc {
	return &cCommand.CommandHelpDoc{
		CommandDesc: cCommand.CommandDesc{Name: "help", Desc: "帮助"},
		MethodDesc: []cCommand.MethodDesc{
			{Name: "doc", Desc: "文档", Options: []string{}},
		},
		OptionDesc: []cCommand.OptionDesc{},
	}
}

func (i *HelpCommand) Doc() {
	var docs []commandHelpDoc
	instances := cCommand.GetAllCommand()
	for _, commandInterface := range instances {
		var optionMap = make(map[string]string)
		commandDoc := commandInterface.Help()
		commandName := commandDoc.CommandDesc.Name
		commandDesc := commandDoc.CommandDesc.Desc
		for _, desc := range commandDoc.OptionDesc {
			optionMap[desc.Name] = desc.Desc
		}
		for _, methodDesc := range commandDoc.MethodDesc {
			var methodOptions = make(map[string]string)
			for _, option := range methodDesc.Options {
				if _, ok := optionMap[option]; ok {
					methodOptions[option] = optionMap[option]
				} else {
					methodOptions[option] = ""
				}
			}

			docs = append(docs, commandHelpDoc{
				name:        commandName,
				method:      methodDesc.Name,
				desc:        fmt.Sprintf("%s%s", commandDesc, methodDesc.Desc),
				options:     methodDesc.Options,
				optionsDesc: methodOptions,
			})
		}
	}

	for _, doc := range docs {
		var optionStr string
		for _, option := range doc.options {
			optionStr = fmt.Sprintf("%s --%s=%s", optionStr, option, option)
		}
		color.Set(color.FgGreen)
		fmt.Printf("\n    %s::%s [%s ] %s\n", doc.name, doc.method, optionStr, doc.desc)
		for _, option := range doc.options {
			color.Set(color.FgBlue)
			fmt.Printf("        --%-10s%s\n", option, doc.optionsDesc[option])
		}
	}
}
