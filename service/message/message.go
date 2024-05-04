package main

import (
	"flag"
	"fmt"

	"mini-tiktok/service/message/internal/config"
	"mini-tiktok/service/message/internal/server"
	"mini-tiktok/service/message/internal/svc"
	"mini-tiktok/service/message/pb/message"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// var configFile = flag.String("f", "service/message/etc/message.yaml", "the config file")
var nacosConfigFile = flag.String("f", "service/message/etc/nacos.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	var nacosConf config.NacosConf

	conf.MustLoad(*nacosConfigFile, &nacosConf)
	nacosConf.LoadConfig(&c)

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		message.RegisterMessageServer(grpcServer, server.NewMessageServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
