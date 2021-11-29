package vars

import (
	"gitee.com/csingo/ctool/config/typs"
	"gitee.com/csingo/ctool/core/cHelper"
)

var ConfigCenter = &typs.ConfigCenterConf{
	Enable:   cHelper.EnvToBool("CONFIG_CENTER_ENABLE", true),
	Driver:   cHelper.EnvToString("CONFIG_CENTER_ENABLE", "nacos"),
	Interval: cHelper.EnvToInt("CONFIG_CENTER_ENABLE", 1),
	Channels: &typs.ConfigCenterConf_Channels{
		Nacos: &typs.ConfigCenterConf_NacosChannel{
			Client: &typs.ConfigCenterConf_NacosChannel_Client{
				Scheme:    "http",
				Host:      cHelper.EnvToString("NACOS_HOST", "nacos.qdtech.ai"),
				Port:      uint64(cHelper.EnvToInt("NACOS_PORT", 8848)),
				Username:  cHelper.EnvToString("NACOS_USERNAME", "nacos"),
				Password:  cHelper.EnvToString("NACOS_PASSWORD", "nacos"),
				Namespace: cHelper.EnvToString("APP_ENV", "dev"),
				Path:      "/nacos",
			},
			Listeners: []*typs.ConfigCenterConf_NacosChannel_Listener{
				{Group: "base-go", Data: "DatabaseConf"},
				{Group: "base-go", Data: "RedisConf"},
				//{Group: "eztalk-monitor", Data: "ServerConf"},
			},
		},
		Apollo: nil,
	},
}
