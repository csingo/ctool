package vars

import (
	"gitee.com/csingo/ctool/config/typs"
	"gitee.com/csingo/ctool/core/cHelper"
)

var Redis = &typs.RedisConf{
	Host: cHelper.EnvToString("REDIS_HOST", "127.0.0.1"),
	Port: cHelper.EnvToInt("REDIS_PORT", 6379),
	Auth: cHelper.EnvToString("REDIS_AUTH", ""),
	Db:   cHelper.EnvToInt("REDIS_DB", 0),
}
