package cServer

import (
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"gitee.com/csingo/ctool/core/cCommand"
	"gitee.com/csingo/ctool/core/cConfig"
	"gitee.com/csingo/ctool/core/cLog"
	"gitee.com/csingo/ctool/core/cMiddleware"
	"gitee.com/csingo/ctool/core/cRpc"
)

// Inject 注入实体
func Inject(instance interface{}) {
	name := reflect.TypeOf(instance).Elem().Name()
	path := reflect.TypeOf(instance).Elem().PkgPath()
	index := path + "/" + name
	serverContainer.instances[index] = instance

	cConfig.Inject(instance)
	cCommand.Inject(instance)
	cMiddleware.Inject(instance)
	cRpc.Inject(instance)
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

	if _, ok := serverContainer.instances[name]; ok {
		return serverContainer.instances[name]
	}

	return nil
}

func Load() {
	listenExit()
	cConfig.ListenConf([]string{"ServerConf"}, ReloadHTTP)
}

func Init() {
	cLog.Load()
	cConfig.Load()
	cCommand.Load()
	Load()

	cCommand.SetSysExitChannel(app.exit)

	// 启动服务
	go StartHTTP()
}

func Start() {
	// 系统信号监听
	signal.Notify(app.signal, os.Interrupt, os.Kill, syscall.SIGTERM)
	// 进程常驻
	for {
		select {
		case <-app.signal:
			app.exitCounter = app.exitCounter + 2
			Exit()
			cCommand.Exit()
		case exit := <-app.exit:
			app.exitCounter--
			if exit && app.exitCounter == 0 {
				os.Exit(0)
			}
		}
	}
}
