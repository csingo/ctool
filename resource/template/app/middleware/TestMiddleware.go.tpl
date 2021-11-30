package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
)

type ##MIDDLEWARE## struct{}

func (i *##MIDDLEWARE##) Middleware() string {
	return "Base"
}

func (i *##MIDDLEWARE##) Handler(c *gin.Context) {
	log.Println("handler")
}
