package models

import "github.com/jinzhu/gorm"

type passwordReset struct {
	gorm.Model
	UserID    uint   `gorm:"not null"`
	Token     string `gorm:"-"`
	TokenHash string `gorm:"not null;unique_index"`
}
