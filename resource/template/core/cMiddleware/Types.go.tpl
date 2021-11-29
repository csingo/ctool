package cMiddleware

import "github.com/gin-gonic/gin"

type MiddlewareInterface interface {
	Middleware() string
	Handler(c *gin.Context)
}

type container struct {
	instances    map[string]interface{}
	iandlerFuncs map[string]gin.HandlerFunc
}

var middlewareContainer = &container{
	instances:    map[string]interface{}{},
	iandlerFuncs: map[string]gin.HandlerFunc{},
}
