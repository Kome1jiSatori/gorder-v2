package main

import (
	"context"
	"github.com/Kome1jiSatori/gorder-v2/common/config"
	"github.com/Kome1jiSatori/gorder-v2/common/discovery"
	"github.com/Kome1jiSatori/gorder-v2/common/genproto/stockpb"
	"github.com/Kome1jiSatori/gorder-v2/common/server"
	"github.com/Kome1jiSatori/gorder-v2/stock/ports"
	"github.com/Kome1jiSatori/gorder-v2/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")

	// app注入
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := service.NewApplication(ctx)

	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) { // 启动gprc服务
			svc := ports.NewGRPCServer(application)
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
		// 暂时不用
	default:
		panic("unexpected server type")
	}

}
