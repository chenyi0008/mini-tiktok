package models

import (
	"gorm.io/gorm"
	"mini-tiktok/common/customgorm"
	"mini-tiktok/service/core/helper"
)

type VideoModel struct {
	Author   user   `json:"author" gorm:"foreignKey:UserId" `
	UserId   int64  `json:"user_id"`
	CoverURL string `json:"cover_url"` // 视频封面地址
	PlayURL  string `json:"play_url"`  // 视频播放地址
	Title    string `json:"title"`     // 视频标题
	customgorm.CustomModel
}

type user struct {
	ID   int64  `json:"id"`   // 用户id
	Name string `json:"name"` // 用户名称
}

func (user) TableName() string {
	return "tb_user"
}

func (VideoModel) TableName() string {
	return "tb_video"
}

type Video struct {
	VideoModel
	FavoriteCount int64 `json:"favorite_count"`       // 视频的点赞总数
	IsFavorite    bool  `json:"is_favorite" gorm:"-"` // true-已点赞，false-未点赞
	CommentCount  int64 `json:"comment_count"`        // 视频的评论总数
}

type DefaultVideoModel struct {
	Db *gorm.DB
}

func NewVideoModel(db *gorm.DB) *DefaultVideoModel {
	return &DefaultVideoModel{
		Db: db,
	}
}

func (m *DefaultVideoModel) ListByCreatedAt(latestTime int64, N uint) ([]VideoModel, error) {
	videoList := make([]VideoModel, N)
	err := m.Db.Preload("Author").Where("created_at < ?", helper.TransformUnixToDate(latestTime)).Order("id desc").Limit(int(N)).Find(&videoList).Error
	return videoList, err
}

func (m *DefaultVideoModel) ListAll() ([]VideoModel, error) {
	videoList := make([]VideoModel, 0)
	err := m.Db.Preload("Author").Order("id desc").Find(&videoList).Error
	return videoList, err
}

func (m *DefaultVideoModel) List() ([]VideoModel, error) {
	videoList := make([]VideoModel, 0)
	err := m.Db.Model(&VideoModel{}).Preload("Author").Find(&videoList).Error
	return videoList, err
}

func (m *DefaultVideoModel) ListByUserId(userId uint) ([]VideoModel, error) {
	videoList := make([]VideoModel, 0)
	err := m.Db.Where("user_id = ?", userId).Find(&videoList).Error
	return videoList, err
}

func (m *DefaultVideoModel) Create(userId int64, playURL string, coverURL string, title string) error {
	video := &VideoModel{
		UserId:   userId,
		PlayURL:  playURL,
		CoverURL: coverURL,
		Title:    title,
	}
	err := m.Db.Create(&video).Error
	return err
}

func (m *DefaultVideoModel) ListInIds(arr []uint64) ([]VideoModel, error) {
	videoList := make([]VideoModel, 0)
	err := m.Db.Preload("Author").Where("id in ?", arr).Find(&videoList).Error
	return videoList, err
}

func (m *DefaultVideoModel) FindOneById(id int) (*VideoModel, error) {
	video := &VideoModel{}
	err := m.Db.Preload("Author").Where("id = ?", id).Find(&video).Error
	return video, err
}
