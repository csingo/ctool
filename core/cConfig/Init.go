package cConfig

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"framework/core/cHelper"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	"time"
)

// Inject 注入配置
func Inject(instance interface{}) {
	if !IsConf(instance) {
		return
	}
	name := instance.(ConfigInterface).ConfigName()
	configContainer.instances[name] = instance
	configContainer.versions[name] = "0"
}

// Load 加载配置文件
func Load() {
	// 获取文件配置路径
	initLocalFile()
	// 获取配置中心配置
	initConfigCenter()

	// 加载配置
	log.Println("start load cConfig ...")
	for _, instance := range configContainer.instances {
		name := instance.(ConfigInterface).ConfigName()
		jsonBytes, err := readConf(name)
		if err != nil {
			log.Printf("read cConfig err: [%s] %+v", name, err)
			continue
		}
		err = json.Unmarshal(jsonBytes, instance)
		if err != nil {
			log.Printf("json unmarshal cConfig err: [%s] %+v", name, jsonBytes)
			continue
		}
		hex := md5.Sum(jsonBytes)
		configContainer.versions[name] = fmt.Sprintf("%x", hex)
		// 添加配置监听器
		//ListenConf([]string{name}, nil)
	}
}

// GetConf 获取配置
func GetConf(name string) (interface{}, error) {
	var conf interface{} = nil

	indexArr := strings.Split(name, ".")
	indexLen := len(indexArr)

	if _, ok := configContainer.instances[indexArr[0]]; ok {
		conf = configContainer.instances[indexArr[0]]
	} else {
		return nil, errors.New("cConfig not fount")
	}

	if indexLen > 1 {
		for i := 1; i < indexLen; i++ {
			if cHelper.IsInt(indexArr[i]) {
				index := cHelper.ToInt(indexArr[i])
				temp := reflect.ValueOf(conf).Index(index)
				if temp.IsValid() {
					conf = temp.Interface()
				} else {
					return nil, errors.New("cConfig field is invalid")
				}
			} else {
				temp := reflect.ValueOf(conf).Elem().FieldByName(indexArr[i])
				if temp.IsValid() {
					conf = temp.Interface()
				} else {
					return nil, errors.New("cConfig field is invalid")
				}
			}
		}
	}

	return conf, nil
}

func IsConf(instance interface{}) bool {
	return reflect.TypeOf(instance).Implements(reflect.TypeOf((*ConfigInterface)(nil)).Elem())
}

// ListenConf 添加配置监听器
func ListenConf(confs []string, handler HandleFunc) {
	interval := configContainer.updateInterval
	for _, conf := range confs {
		if _, ok := configContainer.instances[conf]; ok {
			// 添加监听器处理函数
			if handler != nil {
				configContainer.updateFunc[conf] = append(configContainer.updateFunc[conf], handler)
			}

			// 监听器不存在则新建
			if _, tok := configContainer.updateTickers[conf]; !tok {
				ticker := time.NewTicker(time.Duration(interval) * time.Second)
				configContainer.updateTickers[conf] = ticker

				// 启动监听器处理协程
				log.Printf("start listen cConfig: %s", conf)
				go func(name string, ticker *time.Ticker) {
					for {
						<-ticker.C
						// 读取配置
						jsonBytes, err := readConf(name)
						if err != nil {
							//log.Printf("read cConfig err: [%s] %+v", name, err)
							continue
						}
						hex := md5.Sum(jsonBytes)
						encode := fmt.Sprintf("%x", hex)

						// 比较版本，不一致时才执行更新
						if encode != configContainer.versions[name] {
							// 更新配置
							configContainer.versions[name] = encode
							instance := configContainer.instances[name]
							err := json.Unmarshal(jsonBytes, instance)
							if err != nil {
								log.Printf("json unmarshal cConfig err: [%s] %+v", name, err)
							}
							// 执行监听器函数
							for _, f := range configContainer.updateFunc[name] {
								if f != nil {
									go func(function func()) {
										function()
										defer func() {
											if r := recover(); r != nil {
												log.Printf("listen cConfig error: %+v", function)
											}
										}()
									}(f)
								}
							}
						}
					}
				}(conf, ticker)
			}
		}
	}
}

// readConf 读取配置文件内容，启用配置中心时读取配置中心
func readConf(name string) ([]byte, error) {
	var data []byte
	var err error
	if _, ok := configContainer.configCenter.listeners[name]; configContainer.configCenter.enable && ok {
		group := configContainer.configCenter.listeners[name]
		data, err = readConfigCenter(name, group)
		if len(data) <= 0 && err == nil {
			err = errors.New("cConfig center is empty: " + name)
		}
	} else {
		filePath := configContainer.path + "/" + name + ".json"
		data, err = ioutil.ReadFile(filePath)
		if len(data) <= 0 && err == nil {
			err = errors.New("file cConfig is empty: " + name)
		}
	}

	return data, err
}
