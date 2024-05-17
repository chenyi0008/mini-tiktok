package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Favorite struct {
	UserId  uint `json:"userId"`
	VideoId uint `json:"videoId"`
	gorm.Model
}

func (Favorite) TableName() string {
	return "tb_favorite"
}

type DefaultFavoriteModel struct {
	Db *gorm.DB
}

func NewFavoriteModel(db *gorm.DB) *DefaultFavoriteModel {
	return &DefaultFavoriteModel{
		Db: db,
	}
}

func (m *DefaultFavoriteModel) CreateIgnore(userId, videoId int) error {
	// INSERT IGNORE INTO tb_favorite (t1, t2) VALUES (2, 3);
	favorite := &Favorite{
		UserId:  uint(userId),
		VideoId: uint(videoId),
	}
	err := m.Db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(favorite).Error
	return err
}

func (m *DefaultFavoriteModel) Del(userId, videoId int) error {
	err := m.Db.Unscoped().Delete(&Favorite{}, "user_id = ? and video_id = ?", userId, videoId).Error
	return err
}
