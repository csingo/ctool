package command

import (
	"framework/core/cCommand"
	"log"
)

type TestCommand struct{}

func (i *TestCommand) Help() *cCommand.CommandHelpDoc {
	return &cCommand.CommandHelpDoc{
		CommandDesc: cCommand.CommandDesc{Name: "test", Desc: "测试"},
		MethodDesc: []cCommand.MethodDesc{
			{Name: "hello", Desc: "你好", Options: []string{"name"}},
			{Name: "test", Desc: "生成手机号", Options: []string{"name"}},
		},
		OptionDesc: []cCommand.OptionDesc{
			{Name: "name", Desc: "名称"},
		},
	}
}

func (i *TestCommand) Hello(name cCommand.Option) {
	//app.HelloService.Say(context.Background(), nil)
	log.Printf("hello %s", name.Value)
}

func (i *TestCommand) Test() {
}
