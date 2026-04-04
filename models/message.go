package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Content string `json:"content"`
	UserID  uint   `json:"user_id"`
	ChatID  uint   `json:"chat_id"`
}
