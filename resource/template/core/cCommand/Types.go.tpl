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
	instances map[string]interface{}
	labels    []string
}

var commandContainer = &container{
	instances: make(map[string]interface{}),
	labels:    []string{},
}

// CommandState command状态
type commandState struct {
	enable         bool
	cron           *cron.Cron
	exitChannel    chan bool
	appExitChannel chan bool
	running        map[string]string // 进行中任务
	lock           *sync.Mutex       // map 并发控制
}

var state = &commandState{
	enable:         true,
	cron:           cron.New(),
	exitChannel:    make(chan bool),
	appExitChannel: nil,
	running:        map[string]string{},
	lock:           &sync.Mutex{},
}
