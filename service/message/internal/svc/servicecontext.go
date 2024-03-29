package svc

import (
	"mini-tiktok/service/message/internal/config"
	"mini-tiktok/service/message/models"
)

type ServiceContext struct {
	Config       config.Config
	MessageModel *models.DefaultMessageModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := models.InitMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:       c,
		MessageModel: models.NewMessageModel(db),
	}
}
