package svc

import (
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
	"mini-tiktok/service/core/internal/config"
	"mini-tiktok/service/core/internal/middleware"
	"mini-tiktok/service/core/models"
	"mini-tiktok/service/favorite/pb/favorite"
	"mini-tiktok/service/follow/follow"
	"mini-tiktok/service/message/message"
)

type ServiceContext struct {
	Config      config.Config
	Engine      *gorm.DB
	RDB         *redis.Client
	RedisCli    *models.RedisCliModel
	Auth        rest.Middleware
	UserModel   *models.DefaultUserModel
	VideoModel  *models.DefaultVideoModel
	FollowRpc   follow.FollowClient
	FavoriteRpc favorite.FavoriteClient
	MessageRpc  message.MessageClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	engine := models.InitMysql(c.Mysql.DataSource)
	redisEngine := models.InitRedis(c)
	engine.Logger.LogMode(4)
	return &ServiceContext{
		Config:      c,
		Engine:      engine,
		RDB:         redisEngine,
		RedisCli:    models.NewRedisCli(redisEngine),
		UserModel:   models.NewUserModel(engine),
		VideoModel:  models.NewVideoModel(engine),
		Auth:        middleware.NewAuthMiddleware().Handle,
		FollowRpc:   *config.InitFollowClient(c.Etcd.Host),
		FavoriteRpc: *config.InitFavoriteClient(c.Etcd.Host),
		MessageRpc:  *config.InitMessageClient(c.Etcd.Host),
	}
}
