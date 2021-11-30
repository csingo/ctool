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
	rpcServiceContainer.services[name] = instance.(RpcServiceInterface)
}

func IsRpcService(instance interface{}) bool {
	return reflect.TypeOf(instance).Implements(reflect.TypeOf((*RpcServiceInterface)(nil)).Elem())
}

func GetService(app, service string) (RpcServiceInterface, error) {
	name := fmt.Sprintf("%s/%s", app, service)
	if instance, ok := rpcServiceContainer.services[name]; ok {
		return instance, nil
	}

	return nil, errors.New("service not found")
}
