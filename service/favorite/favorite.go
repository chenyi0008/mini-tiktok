package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"mini-tiktok/service/favorite/internal/config"
	"mini-tiktok/service/favorite/internal/server"
	"mini-tiktok/service/favorite/internal/svc"
	"mini-tiktok/service/favorite/pb/favorite"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "service/favorite/etc/favorite.yaml", "the config file")
var nacosConfigFile = flag.String("q", "service/favorite/etc/nacos.yaml", "the config file")
var mode = flag.String("mode", "dev", "run environment")

func main() {
	flag.Parse()

	var c config.Config
	var nacosConfig config.NacosConf
	conf.MustLoad(*nacosConfigFile, &nacosConfig)
	if *mode == "dev" {
		c.Log.Encoding = "explan"
	}
	nacosConfig.LoadConfig(&c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		favorite.RegisterFavoriteServer(grpcServer, server.NewFavoriteServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
