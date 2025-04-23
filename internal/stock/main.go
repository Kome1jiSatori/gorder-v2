package main

import (
	"github.com/Kome1jiSatori/gorder-v2/common/genproto/stockpb"
	"github.com/Kome1jiSatori/gorder-v2/common/server"
	"github.com/Kome1jiSatori/gorder-v2/stock/ports"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")
	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) { // 启动gprc服务
			svc := ports.NewGRPCServer()
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
		// 暂时不用
	default:
		panic("unexpected server type")
	}

}
