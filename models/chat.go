package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	Name     string    `json:"name"`
	IsGroup  bool      `json:"is_group"`
	Users    []User    `gorm:"many2many:user_chats;" json:"users,omitempty"`
	Messages []Message `gorm:"foreignKey:ChatID" json:"messages,omitempty"`
}
