package cMiddleware

import "github.com/gin-gonic/gin"

type MiddlewareInterface interface {
	Middleware() string
	Handler(c *gin.Context)
}

type container struct {
	Instances    map[string]interface{}
	HandlerFuncs map[string]gin.HandlerFunc
}

var MiddlewareContainer = &container{
	Instances:    map[string]interface{}{},
	HandlerFuncs: map[string]gin.HandlerFunc{},
}
