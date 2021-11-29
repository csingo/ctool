package cCommand

import (
	"github.com/robfig/cron"
	"reflect"
	"sync"
)

// CommandInterface command父接口
type CommandInterface interface {
	Help() *CommandHelpDoc
}

// CommandHelpDoc command帮助文档
type CommandHelpDoc struct {
	AutoExit    bool
	CommandDesc CommandDesc
	MethodDesc  []MethodDesc
	OptionDesc  []OptionDesc
}

// CommandDesc command描述
type CommandDesc struct {
	Name string
	Desc string
}

// MethodDesc command方法描述
type MethodDesc struct {
	Name    string
	Desc    string
	Options []string
}

// OptionDesc command方法参数描述
type OptionDesc struct {
	Name string
	Desc string
}

// Option command方法参数
type Option struct {
	Name  string
	Value string
}

// CommandArgvs command命令行输入内容
type CommandArgvs struct {
	Command CommandInterface
	Method  reflect.Value
	Options []*Option
	Name    string
}

// container 容器
type container struct {
	Instances map[string]interface{}
	Labels    []string
}

var CommandContainer = &container{
	Instances: make(map[string]interface{}),
	Labels:    []string{},
}

// CommandState command状态
type CommandState struct {
	Enable         bool
	Cron           *cron.Cron
	ExitChannel    chan bool
	AppExitChannel chan bool
	Running        map[string]string // 进行中任务
	Lock           *sync.Mutex       // map 并发控制
}

var state = &CommandState{
	Enable:         true,
	Cron:           cron.New(),
	ExitChannel:    make(chan bool),
	AppExitChannel: nil,
	Running:        map[string]string{},
	Lock:           &sync.Mutex{},
}
