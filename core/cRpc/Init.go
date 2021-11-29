package cRpc

import (
	"errors"
	"fmt"
	"reflect"
)

func Inject(instance interface{}) {
	if !IsRpcService(instance) {
		return
	}

	app, service := instance.(RpcServiceInterface).RpcServiceName()
	name := fmt.Sprintf("%s/%s", app, service)
	RpcServiceContainer.Services[name] = instance.(RpcServiceInterface)
}

func IsRpcService(instance interface{}) bool {
	return reflect.TypeOf(instance).Implements(reflect.TypeOf((*RpcServiceInterface)(nil)).Elem())
}

func GetService(app, service string) (RpcServiceInterface, error) {
	name := fmt.Sprintf("%s/%s", app, service)
	if instance, ok := RpcServiceContainer.Services[name]; ok {
		return instance, nil
	}

	return nil, errors.New("service not found")
}

//
//func Load() {
//	serviceHosts, err := cConfig.GetConf("RpcConf.ServiceHosts")
//	if err != nil {
//		return
//	}
//
//	length := reflect.ValueOf(serviceHosts).Len()
//	if length <= 0 {
//		return
//	}
//
//	for i := 0; i < length; i++ {
//		instance := reflect.ValueOf(serviceHosts).Index(i).FieldByName("Service").Interface()
//		host := reflect.ValueOf(serviceHosts).Index(i).FieldByName("Host").Interface().(string)
//		if !IsRpc(instance) {
//			continue
//		}
//
//		RpcContainer.RegisterService(instance.(RpcServiceInterface), host)
//		instance.(RpcServiceInterface).RegisterServerHost(host)
//	}
//}
