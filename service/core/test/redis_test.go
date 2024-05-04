package models

import (
	"context"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"mini-tiktok/service/core/internal/config"
	"mini-tiktok/service/core/internal/svc"
	"os"
	"testing"
)

var ctx *svc.ServiceContext

func TestMain(m *testing.M) {
	// 初始化操作
	var nacosConfigFile = flag.String("q", "service/core/etc/nacos.yaml", "the config file")
	var c config.Config
	var nacosConf config.NacosConf
	conf.MustLoad(*nacosConfigFile, &nacosConf)
	nacosConf.LoadConfig(&c)
	ctx = svc.NewServiceContext(c)

	// 运行测试
	code := m.Run()

	// 退出
	os.Exit(code)
}

func TestZRangeByScore(t *testing.T) {
	score, err := ctx.RedisCli.ZRangeByScore(context.Background(), 12)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range score {
		t.Log(item.Score)
	}
}
