package cServer

import (
	"context"
	"##PROJECT##/core/cConfig"
	"##PROJECT##/core/cHelper"
	"##PROJECT##/core/cMiddleware"
	"##PROJECT##/core/cRpc"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// StopHTTP 停止http服务
func StopHTTP() {
	if state.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), state.timeout)
		defer cancel()
		state.httpServer.Shutdown(ctx)
	}
}

// ReloadHTTP 重启http服务
func ReloadHTTP() {
	// 退出 http
	StopHTTP()
	// 启动 http
	StartHTTP()
}

// StartHTTP 启动http服务
func StartHTTP() {
	enable, err := cConfig.GetConf("ServerConf.HttpServer.Enable")
	if err != nil {
		enable = false
	}
	if !enable.(bool) {
		return
	}
	port, err := cConfig.GetConf("ServerConf.HttpServer.Port")
	if err != nil {
		port = 80
	}
	readTimeout, err := cConfig.GetConf("ServerConf.HttpServer.ReadTimeout")
	if err != nil {
		readTimeout = 10
	}
	writeTimeout, err := cConfig.GetConf("ServerConf.HttpServer.WriteTimeout")
	if err != nil {
		writeTimeout = 10
	}
	exitTimeout, err := cConfig.GetConf("ServerConf.HttpServer.ExitTimeout")
	if err != nil {
		exitTimeout = 10
	}
	maxHeaderBytes, err := cConfig.GetConf("ServerConf.HttpServer.MaxHeaderBytes")
	if err != nil {
		maxHeaderBytes = 20
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// 注册路由
	dispatch(router)

	// 注册 rpc server
	router.POST("/rpc/call", handleRpc)

	// 定义 HTTP 配置
	s := &http.Server{
		Addr:           ":" + cHelper.ToString(port.(int)),
		Handler:        router,
		ReadTimeout:    time.Duration(readTimeout.(int)) * time.Second,
		WriteTimeout:   time.Duration(writeTimeout.(int)) * time.Second,
		MaxHeaderBytes: 1 << maxHeaderBytes.(int),
	}

	state.httpServer = s
	state.timeout = time.Duration(exitTimeout.(int)) * time.Second

	// 启动 HTTP 服务
	s.ListenAndServe()
}

// dispatch 分发路由
func dispatch(router *gin.Engine) {
	// 读取路由配置
	routeConf, err := cConfig.GetConf("RouteConf")
	if err != nil {
		log.Printf("%+v", err)
		return
	}

	baseGroup := router.Group("")
	parseRoute(routeConf, baseGroup)
}

// parseRoute 解析路由
func parseRoute(conf interface{}, router *gin.RouterGroup) {
	routes := reflect.ValueOf(conf).Elem().FieldByName("Routes")
	path := reflect.ValueOf(conf).Elem().FieldByName("Path").Interface().(string)

	var frontMiddlewares, postMiddlewares []string
	middlewares := reflect.ValueOf(conf).Elem().FieldByName("Middlewares")
	if !middlewares.IsNil() {
		frontMiddlewares = middlewares.Elem().FieldByName("Fronts").Interface().([]string)
		postMiddlewares = middlewares.Elem().FieldByName("Posts").Interface().([]string)
	}

	routeLen := routes.Len()
	// routeLen > 0 为路由组
	if routeLen > 0 {
		// allMiddlewares 组装所有中间件
		var allMiddlewares []gin.HandlerFunc
		for _, m := range frontMiddlewares {
			h := cMiddleware.GetMiddleware(m)
			if h != nil {
				allMiddlewares = append(allMiddlewares, h)
			}
		}
		for _, m := range postMiddlewares {
			h := cMiddleware.GetMiddleware(m)
			if h != nil {
				allMiddlewares = append(allMiddlewares, h)
			}
		}
		group := router.Group(path, allMiddlewares...)
		for i := 0; i < routeLen; i++ {
			route := routes.Index(i).Interface()
			parseRoute(route, group)
		}
	} else {
		// 路由处理方法参数集合
		var routeHandlerArgs = []reflect.Value{reflect.ValueOf(path)}

		// 获取控制器方法
		handler := reflect.ValueOf(conf).Elem().FieldByName("Handler").Interface().(string)
		handlerArr := strings.Split(handler, "::")
		if len(handlerArr) < 2 {
			return
		}
		controllerName := handlerArr[0]
		controllerMethodName := handlerArr[1]
		controller := Instance(controllerName)
		if controller == nil {
			log.Printf("controller not exists: %s", controllerName)
			return
		}
		controllerMethod := reflect.ValueOf(controller).MethodByName(controllerMethodName)
		if !controllerMethod.IsValid() {
			log.Printf("controller method is not valid: %s", controllerMethodName)
			return
		}
		// 定义控制器路由处理
		var controllerHandler gin.HandlerFunc
		controllerHandler = func(c *gin.Context) {
			state.httpConnCounter++
			results := controllerMethod.Call([]reflect.Value{
				reflect.ValueOf(c),
			})
			if len(results) > 0 {
				httpCode := results[0].Interface().(int)
				result := results[1].Interface()
				c.JSON(httpCode, result)
			}
			state.httpConnCounter--
		}

		// 路由处理方法-前置中间件
		for _, m := range frontMiddlewares {
			h := cMiddleware.GetMiddleware(m)
			if h != nil {
				routeHandlerArgs = append(routeHandlerArgs, reflect.ValueOf(h))
			}
		}
		// 路由处理方法-controller
		routeHandlerArgs = append(routeHandlerArgs, reflect.ValueOf(controllerHandler))
		// 路由处理方法-后置中间件
		for _, m := range postMiddlewares {
			h := cMiddleware.GetMiddleware(m)
			if h != nil {
				routeHandlerArgs = append(routeHandlerArgs, reflect.ValueOf(h))
			}
		}

		// 注册路由处理方法
		method := reflect.ValueOf(conf).Elem().FieldByName("Method").Interface().(string)
		reflect.ValueOf(router).MethodByName(method).Call(routeHandlerArgs)
	}
}

func handleRpc(c *gin.Context) {
	rpcApp := c.Request.Header.Get("Rpc-App")
	rpcService := c.Request.Header.Get("Rpc-Service")
	rpcMethod := c.Request.Header.Get("Rpc-Method")

	serviceInstance, err := cRpc.GetService(rpcApp, rpcService)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		c.Abort()
		return
	}

	caller := reflect.ValueOf(serviceInstance).MethodByName(rpcMethod)
	typ := caller.Type().In(1)
	param := reflect.New(typ)

	jsonCall := reflect.ValueOf(c).MethodByName("BindJSON")
	jsonRes := jsonCall.Call([]reflect.Value{param})
	if !jsonRes[0].IsNil() {
		c.String(http.StatusInternalServerError, jsonRes[0].Interface().(error).Error())
		c.Abort()
		return
	}

	responseValues := caller.Call([]reflect.Value{reflect.ValueOf(c), param.Elem()})
	if !responseValues[1].IsNil() {
		err = responseValues[1].Interface().(error)
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, responseValues[0].Interface())
}
