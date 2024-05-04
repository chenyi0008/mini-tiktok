package svc

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"mini-tiktok/service/favorite/internal/config"
	"mini-tiktok/service/favorite/models"
)

type ServiceContext struct {
	Config        config.Config
	Mysql         *gorm.DB
	Redis         *redis.Client
	FavoriteModel *models.DefaultFavoriteModel
	CommentModel  *models.DefaultCommentModel
	RedisCli      *models.DefaultRedisCliModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := models.InitMysql(c.Mysql.DataSource)
	redis := models.InitRedis(c.RedisInfo.Addr, c.RedisInfo.Password, c.RedisInfo.DB)

	return &ServiceContext{
		Config:        c,
		Mysql:         db,
		Redis:         redis,
		RedisCli:      models.NewRedisCli(redis),
		FavoriteModel: models.NewUserModel(db),
		CommentModel:  models.NewCommentModel(db),
	}
}
