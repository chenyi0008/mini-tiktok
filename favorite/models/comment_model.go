package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserId  uint
	VideoId uint
	Content string
	User    User `json:"author" gorm:"foreignKey:UserId"`
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

func (m *DefaultFavoriteModel) Create(userId uint, videoId uint, content string) error {
	comment := &Comment{
		UserId:  userId,
		VideoId: videoId,
		Content: content,
	}
	err := m.Db.Create(comment).Error
	if err != nil {
		return err
	}
	return nil
}
