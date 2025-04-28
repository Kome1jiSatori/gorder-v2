package server

import (
	"github.com/Kome1jiSatori/gorder-v2/common/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
	setMiddlewares(apiRouter)
	wrapper(apiRouter)      // 通过传过来的函数对gin修改
	apiRouter.Group("/api") // 路由组

	if err := apiRouter.Run(addr); err != nil {
		panic(err)
	}
}

func setMiddlewares(r *gin.Engine) {
	r.Use(middleware.StructuredLog(logrus.NewEntry(logrus.StandardLogger())))
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("default_server"))
}
