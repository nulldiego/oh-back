package model

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	UserID   uint `gorm:"not null"`
	Messages []Message
}

type Message struct {
	gorm.Model        // Already includes CreatedAt
	ChatID     uint   `gorm:"not null" binding:"required"`
	Author     string `gorm:"not null"` // HUMAN or GPT
	Text       string `gorm:"type:text" binding:"required"`
}
