package models

import (
	"fmt"
	"gorm.io/gorm"
	"mini-tiktok/service/core/models"
	"mini-tiktok/utils"
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

func NewUserModel(db *gorm.DB) *DefaultFavoriteModel {
	return &DefaultFavoriteModel{
		Db: db,
	}
}

func (m *DefaultFavoriteModel) GiveLike(userId, videoId uint64) error {
	favorite := &Favorite{UserId: uint(userId), VideoId: uint(videoId)}
	err := m.Db.Create(favorite).Error
	return err
}

func (m *DefaultFavoriteModel) CancelLike(userId, videoId uint64) error {
	err := m.Db.Unscoped().Where("user_id = ? and video_id = ?", userId, videoId).Delete(&Favorite{}).Error
	return err
}

func (m *DefaultFavoriteModel) GetByUserId(userId uint64) (*[]Favorite, error) {
	favorite := make([]Favorite, 0)
	err := m.Db.Where("user_id = ?", userId).Find(&favorite).Error
	return &favorite, err
}

func (m *DefaultFavoriteModel) CountByVideoId(videoId uint) (int, error) {
	var count int64
	err := m.Db.Model(&Favorite{}).Where("video_id = ?", videoId).Count(&count).Error
	return int(count), err
}

func (m *DefaultFavoriteModel) GetVideoFavoriteCountBatch(videoId []int) ([]int, error) {
	length := len(videoId)
	count := make([]int, length)
	videoList := make([]models.Video, length)
	orderBy := utils.BuildSQLStringForIDs(videoId)
	err := m.Db.Model(&models.VideoModel{}).Select("id, favorite_count").Clauses(orderBy).Where("id in ?", videoId).Find(&videoList).Error
	fmt.Println()
	for _, video := range videoList {
		fmt.Println(video)
	}
	fmt.Println()
	for i := 0; i < length; i++ {
		count[i] = int(videoList[i].FavoriteCount)
	}
	return count, err
}

func (m *DefaultFavoriteModel) IsFavor(videoId, userId uint) (bool, error) {
	var count int64
	err := m.Db.Model(&Favorite{}).Where("user_id = ? and video_id = ?", userId, videoId).Count(&count).Error
	if count > 0 {
		return true, err
	} else {
		return false, err
	}
	return false, nil
}
