package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	Name    string `json:"name"`
	IsGroup bool   `json:"is_group"`
}
