package models

import (
	"gorm.io/gorm"
)

type User struct {
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
	FollowCount     int64  `json:"follow_count"`     // 关注总数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
	//ID              int64   `json:"id"`               // 用户id
	Name           string       `json:"name"`            // 用户名称
	Signature      string       `json:"signature"`       // 个人简介
	TotalFavorited int64        `json:"total_favorited"` // 获赞数量
	WorkCount      int64        `json:"work_count"`      // 作品数
	Password       string       `json:"password" //密码`
	VideoList      []VideoModel `json:"video_list" `
	gorm.Model
}

func (User) TableName() string {
	return "tb_user"
}

type DefaultUserModel struct {
	Db *gorm.DB
}

func NewUserModel(db *gorm.DB) *DefaultUserModel {
	return &DefaultUserModel{
		Db: db,
	}
}

func (m *DefaultUserModel) GetByName(name string) (*User, error) {
	user := new(User)
	err := m.Db.Where("name = ?", name).First(user).Error
	return user, err
}

func (m *DefaultUserModel) GetById(id uint) (*User, error) {
	user := new(User)
	err := m.Db.Where("id = ?", id).First(user).Error
	return user, err
}

func (m *DefaultUserModel) Create(name string, password string) error {
	user := new(User)
	user.Name = name
	user.Password = password
	err := m.Db.Create(user).Error
	return err
}
