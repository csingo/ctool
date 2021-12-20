package vars

import (
	"github.com/csingo/ctool/config/typs"
	"github.com/csingo/ctool/core/cHelper"
)

var Database = &typs.DatabaseConf{
	Driver: cHelper.EnvToString("DB_DRIVER", "default"),
	Channels: &typs.DatabaseConf_Channels{
		Default: &typs.DatabaseConf_MysqlChannel{
			Dsn:             cHelper.EnvToString("DB_DSN", ""),
			Host:            cHelper.EnvToString("DB_HOST", "127.0.0.1"),
			Port:            cHelper.EnvToInt("DB_PORT", 3306),
			Database:        cHelper.EnvToString("DB_NAME", "test"),
			Username:        cHelper.EnvToString("DB_USERNAME", "root"),
			Password:        cHelper.EnvToString("DB_PASSWORD", "123456"),
			Charset:         "utf8mb4",
			Parsetime:       "True",
			Loc:             "Local",
			MaxIgleConns:    5,
			MaxOpenConns:    10,
			ConnMaxLifetime: 60,
		},
	},
}
