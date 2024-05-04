package main

import (
	"flag"
	"fmt"

	"mini-tiktok/service/follow/internal/config"
	"mini-tiktok/service/follow/internal/server"
	"mini-tiktok/service/follow/internal/svc"
	"mini-tiktok/service/follow/pb/follow"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var nacosConfigFile = flag.String("f", "service/follow/etc/nacos.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	var nacosConf config.NacosConf

	conf.MustLoad(*nacosConfigFile, &nacosConf)

	nacosConf.LoadConfig(&c)

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		follow.RegisterFollowServer(grpcServer, server.NewFollowServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
