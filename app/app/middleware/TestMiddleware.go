package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
)

type TestMiddleware struct{}

func (i *TestMiddleware) Middleware() string {
	return "Base"
}

func (i *TestMiddleware) Handler(c *gin.Context) {
	log.Printf("[%s] %s", c.Request.Method, c.Request.RequestURI)
}
