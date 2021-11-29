package vars

import (
	"framework/config/typs"
	"framework/core/cHelper"
)

var Server = &typs.ServerConf{
	GoMod: "framework",
	HttpServer: &typs.ServerConf_HttpServer{
		Enable:         cHelper.EnvToBool("HTTP_SERVER_ENABLE", true),
		Port:           cHelper.EnvToInt("HTTP_SERVER_PORT", 8080),
		ReadTimeout:    cHelper.EnvToInt("HTTP_SERVER_READ_TIMEOUT", 5),
		WriteTimeout:   cHelper.EnvToInt("HTTP_SERVER_WRITE_TIMEOUT", 10),
		ExitTimeout:    cHelper.EnvToInt("HTTP_SERVER_EXIT_TIMEOUT", 30),
		MaxHeaderBytes: cHelper.EnvToInt("HTTP_SERVER_MAX_HEADER_BYTES", 20),
	},
}
