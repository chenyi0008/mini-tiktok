package main

import (
	"flag"
	"fmt"
	"mini-tiktok/service/core/corn"
	"mini-tiktok/service/core/internal/config"
	"mini-tiktok/service/core/internal/handler"
	"mini-tiktok/service/core/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

// var configFile = flag.String("f", "service/core/etc/core-api.yaml", "the config file")
var nacosConfigFile = flag.String("q", "service/core/etc/nacos.yaml", "the config file")

func main() {
	flag.Parse()

	var nacosConf config.NacosConf
	var c config.Config
	//conf.MustLoad(*configFile, &c)
	conf.MustLoad(*nacosConfigFile, &nacosConf)
	nacosConf.LoadConfig(&c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	go corn.CornInit(ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()

}
