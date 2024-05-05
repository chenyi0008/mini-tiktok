package models

import (
	"fmt"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserId  uint
	VideoId uint
	Content string
	User    User `json:"author" customgorm:"foreignKey:UserId"`
}

func (Comment) TableName() string {
	return "tb_comment"
}

type DefaultCommentModel struct {
	Db *gorm.DB
}

func NewCommentModel(db *gorm.DB) *DefaultCommentModel {
	return &DefaultCommentModel{
		Db: db,
	}
}

func (m *DefaultCommentModel) GetByVideoId(videoId uint) ([]Comment, error) {
	commentList := make([]Comment, 0)
	err := m.Db.Preload("User").Where("video_id = ?", videoId).Find(&commentList).Error
	return commentList, err
}

func (m *DefaultCommentModel) Create(userId uint, videoId uint, content string) (*Comment, error) {
	comment := &Comment{
		UserId:  userId,
		VideoId: videoId,
		Content: content,
	}
	err := m.Db.Create(comment).Error
	fmt.Printf("comment:%+v", comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (m *DefaultCommentModel) CountByVideoId(videoId uint) (int64, error) {
	var count int64
	err := m.Db.Model(&Comment{}).Where("video_id = ?", videoId).Count(&count).Error
	return count, err
}
