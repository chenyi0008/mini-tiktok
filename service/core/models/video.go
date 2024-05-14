package models

import (
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"mini-tiktok/common/consts"
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

// AggregateVideoFavorCount 聚合查询视频视频点赞数量
func (m *DefaultVideoModel) AggregateVideoFavorCount() (*map[int]int, error) {
	// 查询每个视频有多少个不同的用户ID关联
	var videoFavorCounts []struct {
		VideoID   int
		UserCount int
	}
	err := m.Db.Table(consts.FAVORITE).
		Select("video_id, COUNT(user_id) as user_count").
		Group("video_id").
		Scan(&videoFavorCounts).Error
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	// 将结果存储到map中
	videoFavorCountMap := make(map[int]int)
	for _, item := range videoFavorCounts {
		videoFavorCountMap[item.VideoID] = item.UserCount
	}
	return &videoFavorCountMap, nil
}

// AggregateVideoCommentCount 聚合查询视频评论数量
func (m *DefaultVideoModel) AggregateVideoCommentCount() (*map[int]int, error) {

	var videoCommentCounts []struct {
		VideoID   int
		UserCount int
	}
	err := m.Db.Table(consts.COMMENT).
		Select("video_id, COUNT(user_id) as user_count").
		Group("video_id").
		Scan(&videoCommentCounts).Error
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	// 将结果存储到map中
	videoFavorCountMap := make(map[int]int)
	for _, item := range videoCommentCounts {
		videoFavorCountMap[item.VideoID] = item.UserCount
	}
	return &videoFavorCountMap, nil
}

// UpdateVideoFavorCount 更新点赞数量
func (m *DefaultVideoModel) UpdateVideoFavorCount(id, count int) (*VideoModel, error) {
	video := &VideoModel{}
	err := m.Db.Table(consts.VIDEO).Where("id = ?", id).UpdateColumn("favorite_count", count).Error
	return video, err
}

// UpdateVideoCommentCount 更新评论数量
func (m *DefaultVideoModel) UpdateVideoCommentCount(id, count int) (*VideoModel, error) {
	video := &VideoModel{}
	err := m.Db.Preload(consts.VIDEO).Where("id = ?", id).UpdateColumn("comment_count", count).Error
	return video, err
}
