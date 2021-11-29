package cConfig

import (
	"gitee.com/csingo/ctool/core/cHelper"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
	"reflect"
)

func initNacos() {
	scheme, err1 := GetConf("ConfigCenterConf.Channels.Nacos.Client.Scheme")
	host, err2 := GetConf("ConfigCenterConf.Channels.Nacos.Client.Host")
	port, err3 := GetConf("ConfigCenterConf.Channels.Nacos.Client.Port")
	username, err4 := GetConf("ConfigCenterConf.Channels.Nacos.Client.Username")
	password, err5 := GetConf("ConfigCenterConf.Channels.Nacos.Client.Password")
	namespace, err6 := GetConf("ConfigCenterConf.Channels.Nacos.Client.Namespace")
	path, err7 := GetConf("ConfigCenterConf.Channels.Nacos.Client.Path")
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil {
		log.Printf("init nacos err : %+v", []error{err1, err2, err3, err4, err5, err6, err7})
		configContainer.configCenter.enable = false
		return
	}
	var err error
	configContainer.configCenter.nacosClient, err = clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig: &constant.ClientConfig{
			NamespaceId:          namespace.(string),
			UpdateThreadNum:      0,
			NotLoadCacheAtStart:  true,
			UpdateCacheWhenEmpty: false,
			Username:             username.(string),
			Password:             password.(string),
		},
		ServerConfigs: []constant.ServerConfig{
			{
				IpAddr:      host.(string),
				ContextPath: path.(string),
				Port:        port.(uint64),
				Scheme:      scheme.(string),
			},
		},
	})
	if err != nil {
		log.Printf("create nacos client error : %+v", err)
		configContainer.configCenter.enable = false
		return
	}

	listeners, err := GetConf("ConfigCenterConf.Channels.Nacos.Listeners")
	if err != nil {
		return
	}
	len := reflect.ValueOf(listeners).Len()
	if len > 0 {
		for i := 0; i < len; i++ {
			groupIndex := "ConfigCenterConf.Channels.Nacos.Listeners." + cHelper.ToString(i) + ".Group"
			dataIndex := "ConfigCenterConf.Channels.Nacos.Listeners." + cHelper.ToString(i) + ".Data"

			group, errG := GetConf(groupIndex)
			data, errD := GetConf(dataIndex)
			if errG != nil || errD != nil {
				continue
			}
			configContainer.configCenter.listeners[data.(string)] = group.(string)
		}
	}
}

func readNacos(data, group string) ([]byte, error) {
	content, err := configContainer.configCenter.nacosClient.GetConfig(vo.ConfigParam{
		DataId: data,
		Group:  group,
	})
	return []byte(content), err
}
