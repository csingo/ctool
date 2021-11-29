package vars

import "gitee.com/csingo/ctool/config/typs"

var Tool = &typs.ToolConf{
	CopyFiles: []string{
		// autoload, 自动加载目录
		"autoload/.gitkeep",
		"autoload/command.go",
		"autoload/controller.go",
		"autoload/loader.go",
		"autoload/middleware.go",
		"autoload/rpc.go",
		"autoload/service.go",
		// base, sdk客户端目录
		"base/.gitkeep",
		// bin, 二进制执行文件目录
		"bin/.gitkeep",
		// config, 配置目录
		"config/.gitkeep",
		"config/typs/.gitkeep",
		"config/vars/.gitkeep",
		// config/typs，配置定义目录
		"config/typs/CommandConf.go",
		"config/typs/ConfigCenterConf.go",
		"config/typs/DatabaseConf.go",
		"config/typs/RedisConf.go",
		"config/typs/RouteConf.go",
		"config/typs/RpcConf.go",
		"config/typs/ServerConf.go",
		// config/vars，配置实例目录
		"config/vars/CommandConf.go",
		"config/vars/ConfigCenterConf.go",
		"config/vars/DatabaseConf.go",
		"config/vars/RedisConf.go",
		"config/vars/RouteConf.go",
		"config/vars/RpcConf.go",
		"config/vars/ServerConf.go",
		// core, 核心框架目录
		"core/.gitkeep",
		// core/cCommand
		"core/cCommand/.gitkeep",
		"core/cCommand/Init.go",
		"core/cCommand/State.go",
		"core/cCommand/Types.go",
		// core/cConfig
		"core/cConfig/.gitkeep",
		"core/cConfig/Center.go",
		"core/cConfig/File.go",
		"core/cConfig/Init.go",
		"core/cConfig/Nacos.go",
		"core/cConfig/Types.go",
		// core/cHelper
		"core/cHelper/.gitkeep",
		"core/cHelper/Array.go",
		"core/cHelper/Env.go",
		"core/cHelper/File.go",
		"core/cHelper/Math.go",
		"core/cHelper/String.go",
		"core/cHelper/constants/randomType/Index.go",
		"core/cHTTPClient/Request.go",
		// core/cHTTPClient
		"core/cHTTPClient/.gitkeep",
		"core/cHTTPClient/Request.go",
		// core/cMiddleware
		"core/cMiddleware/.gitkeep",
		"core/cMiddleware/Init.go",
		"core/cMiddleware/Types.go",
		// core/cRpc
		"core/cRpc/.gitkeep",
		"core/cRpc/Init.go",
		"core/cRpc/Types.go",
		// core/cServer
		"core/cServer/.gitkeep",
		"core/cServer/Http.go",
		"core/cServer/Init.go",
		"core/cServer/State.go",
		"core/cServer/Types.go",
		// docker
		"docker/.gitkeep",
		"docker/deployment.yml",
		"docker/Dockerfile",
		// resource
		"resource/.gitkeep",
		// response
		"response/.gitkeep",
		"response/Function.go",
		"response/Type.go",
		// server
		"server/.gitkeep",
		"server/main.go",
		// ./
		"Makefile",
		".gitignore",
	},
	RewiteParams: []string{},
}
