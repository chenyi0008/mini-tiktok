package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	RedisInfo struct {
		Addr     string
		Password string
		DB       int
	}
	Queue struct {
		Addr  string
		Topic string
	}
}
