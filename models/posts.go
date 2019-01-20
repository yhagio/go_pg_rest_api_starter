package models

import (
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	UserID      uint   `gorm:"not null; index"`
}

// PostService is a set of methods used to manipulate and
// work with the Post model
type PostService interface {
	PostDB
}

func NewPostService(db *gorm.DB) PostService {
	return &postService{
		PostDB: &postValidator{
			PostDB: &postGorm{
				db: db,
			},
		},
	}
}

var _ PostService = &postService{}

type postService struct {
	PostDB
}
