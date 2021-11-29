package cMiddleware

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

func Inject(instance interface{}) {
	if !IsMiddleware(instance) {
		return
	}

	name := instance.(MiddlewareInterface).Middleware()
	var function gin.HandlerFunc
	function = func(c *gin.Context) {
		instance.(MiddlewareInterface).Handler(c)
	}

	MiddlewareContainer.Instances[name] = instance
	MiddlewareContainer.HandlerFuncs[name] = function
}

func IsMiddleware(o interface{}) bool {
	return reflect.TypeOf(o).Implements(reflect.TypeOf((*MiddlewareInterface)(nil)).Elem())
}

func GetMiddleware(name string) gin.HandlerFunc {
	if _, ok := MiddlewareContainer.HandlerFuncs[name]; !ok {
		return nil
	}

	return MiddlewareContainer.HandlerFuncs[name]
}
