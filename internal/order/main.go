package main

import (
	"context"
	"github.com/Kome1jiSatori/gorder-v2/common/tracing"
	"log"

	"github.com/Kome1jiSatori/gorder-v2/common/broker"
	"github.com/Kome1jiSatori/gorder-v2/common/config"
	"github.com/Kome1jiSatori/gorder-v2/common/discovery"
	"github.com/Kome1jiSatori/gorder-v2/common/genproto/orderpb"
	"github.com/Kome1jiSatori/gorder-v2/common/logging"
	"github.com/Kome1jiSatori/gorder-v2/common/server"
	"github.com/Kome1jiSatori/gorder-v2/order/infrastructure/consumer"
	"github.com/Kome1jiSatori/gorder-v2/order/ports"
	"github.com/Kome1jiSatori/gorder-v2/order/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

// 开启服务
func main() {
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	defer shutdown(ctx)

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	// 注册到consul
	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	// mq初始化
	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeCh()
	}()

	go consumer.NewConsumer(application).Listen(ch)

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		svc := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, svc)
	})

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		router.StaticFile("/success", "../../public/success.html")
		ports.RegisterHandlersWithOptions(router, HTTPServer{
			app: application,
		}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})
}
