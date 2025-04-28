package server

import (
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"net"
)

func init() {
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger))
}

func RunGRPCServer(ServiceName string, registerServer func(server *grpc.Server)) {
	addr := viper.Sub(ServiceName).GetString("grpc-addr")
	if addr == "" {
		// TODO: warning log
		addr = viper.GetString("fallback-grpc-addr")
	}
	RunGRPCServerOnAddr(addr, registerServer)
}

func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			grpc_tags.UnaryServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
		),
		grpc.ChainStreamInterceptor(
			grpc_tags.StreamServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry), // 通过extractor提取结构体然后打日志
		),
	)
	registerServer(grpcServer)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Infof("Starting gRPC server, Listening %s", addr)
	if err := grpcServer.Serve(listen); err != nil {
		logrus.Panic(err)
	}
}
