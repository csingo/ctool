package controller

import (
	"##PROJECT##/app/service"
	"##PROJECT##/response"
	"github.com/gin-gonic/gin"
)

type ##CONTROLLER## struct{}

func (i *##CONTROLLER##) Hello(c *gin.Context) (int, interface{}) {
	return response.Success("data")
}
