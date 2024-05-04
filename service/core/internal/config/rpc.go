package config

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"mini-tiktok/service/favorite/pb/favorite"
	"mini-tiktok/service/follow/follow"
	"mini-tiktok/service/message/message"
)

var FavoriteClient favorite.FavoriteClient
var FollowClient follow.FollowClient
var MessageClient message.MessageClient

func InitFollowClient(etcdHost string) *follow.FollowClient {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{ // 通过 etcd 服务发现时，只需要给 Etcd 配置即可
			Hosts: []string{etcdHost},
			Key:   "follow.rpc",
		},
	})
	client := follow.NewFollowClient(conn.Conn())
	return &client
}

func InitFavoriteClient(etcdHost string) *favorite.FavoriteClient {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{ // 通过 etcd 服务发现时，只需要给 Etcd 配置即可
			Hosts: []string{etcdHost},
			Key:   "favorite.rpc",
		},
	})
	client := favorite.NewFavoriteClient(conn.Conn())
	return &client
}

func InitMessageClient(etcdHost string) *message.MessageClient {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{ // 通过 etcd 服务发现时，只需要给 Etcd 配置即可
			Hosts: []string{etcdHost},
			Key:   "message.rpc",
		},
	})
	client := message.NewMessageClient(conn.Conn())
	return &client
}
