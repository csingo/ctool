package cServer

import (
	"framework/core/cCommand"
	"framework/core/cConfig"
	"framework/core/cMiddleware"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

// Inject 注入实体
func Inject(instance interface{}) {
	name := reflect.TypeOf(instance).Elem().Name()
	path := reflect.TypeOf(instance).Elem().PkgPath()
	index := path + "/" + name
	ServerContainer.Instances[index] = instance

	cConfig.Inject(instance)
	cCommand.Inject(instance)
	cMiddleware.Inject(instance)
}

// Instance 获取容器内实体
func Instance(name string) interface{} {
	mod, err := cConfig.GetConf("ServerConf.GoMod")
	if err != nil {
		mod = ""
	}
	if mod.(string) != "" {
		name = mod.(string) + "/" + name
	}

	if _, ok := ServerContainer.Instances[name]; ok {
		return ServerContainer.Instances[name]
	}

	return nil
}

func Load() {
	listenExit()
	cConfig.ListenConf([]string{"ServerConf"}, ReloadHTTP)
}

func Init() {
	cConfig.Load()
	cCommand.Load()
	Load()

	cCommand.SetSysExitChannel(App.Exit)

	// 启动服务
	go StartHTTP()
}

func Start() {
	// 系统信号监听
	signal.Notify(App.Signal, os.Interrupt, os.Kill, syscall.SIGTERM)
	// 进程常驻
	for {
		select {
		case <-App.Signal:
			App.ExitCounter = App.ExitCounter + 2
			Exit()
			cCommand.Exit()
		case exit := <-App.Exit:
			App.ExitCounter--
			if exit && App.ExitCounter == 0 {
				os.Exit(0)
			}
		}
	}
}
