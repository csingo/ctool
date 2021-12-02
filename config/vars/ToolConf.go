package vars

import "gitee.com/csingo/ctool/config/typs"

var Tool = &typs.ToolConf{
	Version:               "v0.0.27",
	ProtoGenGoPackage:     "google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1",
	ProtoGenGoGrpcPackage: "google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0",
	ProtocDownload: &typs.ToolConf_ProtocAddr{
		Win:   "https://github.com/protocolbuffers/protobuf/releases/download/v3.19.1/protoc-3.19.1-win64.zip",
		Mac:   "https://github.com/protocolbuffers/protobuf/releases/download/v3.19.1/protoc-3.19.1-osx-x86_64.zip",
		Linux: "https://github.com/protocolbuffers/protobuf/releases/download/v3.19.1/protoc-3.19.1-linux-x86_64.zip",
	},
	WriteFiles: map[string][]string{
		"project::create": {
			"/.gitignore.tpl",
			"/Makefile.tpl",
			"/app/.gitkeep.tpl",
			"/autoload/.gitkeep.tpl",
			"/autoload/command.go.tpl",
			"/autoload/config.go.tpl",
			"/autoload/controller.go.tpl",
			"/autoload/loader.go.tpl",
			"/autoload/middleware.go.tpl",
			"/autoload/rpc.go.tpl",
			"/autoload/service.go.tpl",
			"/base/.gitkeep.tpl",
			"/bin/.gitkeep.tpl",
			"/bin/init.sh.tpl",
			"/config/.gitkeep.tpl",
			"/config/typs/.gitkeep.tpl",
			"/config/typs/CommandConf.go.tpl",
			"/config/typs/ConfigCenterConf.go.tpl",
			"/config/typs/DatabaseConf.go.tpl",
			"/config/typs/LogConf.go.tpl",
			"/config/typs/RedisConf.go.tpl",
			"/config/typs/RouteConf.go.tpl",
			"/config/typs/RpcConf.go.tpl",
			"/config/typs/ServerConf.go.tpl",
			"/config/vars/.gitkeep.tpl",
			"/config/vars/CommandConf.go.tpl",
			"/config/vars/ConfigCenterConf.go.tpl",
			"/config/vars/DatabaseConf.go.tpl",
			"/config/vars/LogConf.go.tpl",
			"/config/vars/RedisConf.go.tpl",
			"/config/vars/RouteConf.go.tpl",
			"/config/vars/RpcConf.go.tpl",
			"/config/vars/ServerConf.go.tpl",
			"/conn/.gitkeep.tpl",
			"/conn/database.go.tpl",
			"/conn/redis.go.tpl",
			"/core/.gitkeep.tpl",
			"/core/cCommand/.gitkeep.tpl",
			"/core/cCommand/Init.go.tpl",
			"/core/cCommand/State.go.tpl",
			"/core/cCommand/Types.go.tpl",
			"/core/cConfig/.gitkeep.tpl",
			"/core/cConfig/Center.go.tpl",
			"/core/cConfig/File.go.tpl",
			"/core/cConfig/Init.go.tpl",
			"/core/cConfig/Nacos.go.tpl",
			"/core/cConfig/Types.go.tpl",
			"/core/cHTTPClient/Request.go.tpl",
			"/core/cHelper/.gitkeep.tpl",
			"/core/cHelper/Array.go.tpl",
			"/core/cHelper/Env.go.tpl",
			"/core/cHelper/File.go.tpl",
			"/core/cHelper/Math.go.tpl",
			"/core/cHelper/String.go.tpl",
			"/core/cHelper/constants/randomType/Index.go.tpl",
			"/core/cLog/.gitkeep.tpl",
			"/core/cLog/Init.go.tpl",
			"/core/cMiddleware/.gitkeep.tpl",
			"/core/cMiddleware/Init.go.tpl",
			"/core/cMiddleware/Types.go.tpl",
			"/core/cRpc/.gitkeep.tpl",
			"/core/cRpc/Init.go.tpl",
			"/core/cRpc/Types.go.tpl",
			"/core/cServer/.gitkeep.tpl",
			"/core/cServer/Http.go.tpl",
			"/core/cServer/Init.go.tpl",
			"/core/cServer/State.go.tpl",
			"/core/cServer/Types.go.tpl",
			"/docker/.gitkeep.tpl",
			"/docker/Dockerfile.tpl",
			"/docker/deployment.yml.tpl",
			"/response/.gitkeep.tpl",
			"/response/Function.go.tpl",
			"/response/Type.go.tpl",
			"/server/.gitkeep.tpl",
			"/server/main.go.tpl",
		},
		"app::create": {
			"/app/app/.gitkeep.tpl",
			"/app/app/command/.gitkeep.tpl",
			"/app/app/controller/.gitkeep.tpl",
			"/app/app/middleware/.gitkeep.tpl",
			"/app/app/service/.gitkeep.tpl",
			"/app/app/validator/.gitkeep.tpl",
		},
	},
}
