package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
	addr := viper.Sub(serviceName).GetString("http-addr") // 适配不同配置
	if addr == "" {
		panic("Empty http address")
	}
	RunHTTPServerOnAddr(addr, wrapper)
}

func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
	apiRouter := gin.New()
	wrapper(apiRouter)      // 通过传过来的函数对gin修改
	apiRouter.Group("/api") // 路由组

	if err := apiRouter.Run(addr); err != nil {
		panic(err)
	}
}
