package command

import (
	"##PROJECT##/core/cCommand"
	"log"
)

type ##COMMAND## struct{}

func (i *##COMMAND##) Help() *cCommand.CommandHelpDoc {
	return &cCommand.CommandHelpDoc{
		CommandDesc: cCommand.CommandDesc{Name: "##COMMANDNAME##", Desc: ""},
		MethodDesc: []cCommand.MethodDesc{
			{Name: "hello", Desc: "你好", Options: []string{"name"}},
		},
		OptionDesc: []cCommand.OptionDesc{
			{Name: "name", Desc: "名称"},
		},
	}
}

func (i *##COMMAND##) Hello(name cCommand.Option) {
	log.Printf("hello %s", name.Value)
}
