package controller

import (
	"gitee.com/csingo/ctool/app/service"
	"gitee.com/csingo/ctool/response"
	"github.com/gin-gonic/gin"
)

type HomeController struct{}

func (i *HomeController) Ping(c *gin.Context) (int, interface{}) {
	data := "pong"
	return response.Success(data)
}

func (i *HomeController) Hello(c *gin.Context) (int, interface{}) {
	s := service.HelloService{}
	data, _ := s.Say(c, nil)

	return response.Success(data)
}
