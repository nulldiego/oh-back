package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(15);uniqueIndex;not null" json:"username"`
	Password string `json:"-"`
	Chats    []Chat
}
