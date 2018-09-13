package models

import "github.com/jinzhu/gorm"

type postGorm struct {
	db *gorm.DB
}

var _ PostDB = &postGorm{}

type PostDB interface {
	GetOneById(id uint) (*Post, error)
	GetAllByUserId(userId uint) ([]Post, error)
	Create(post *Post) error
	Update(post *Post) error
	Delete(id uint) error
}

func (pg *postGorm) GetOneById(id uint) (*Post, error) {
	var post Post
	db := pg.db.Where("id = ?", id)
	err := First(db, &post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (pg *postGorm) GetAllByUserId(userId uint) ([]Post, error) {
	var posts []Post
	db := pg.db.Where("user_id = ?", userId)
	err := db.Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (pg *postGorm) Create(post *Post) error {
	return pg.db.Create(post).Error
}

func (pg *postGorm) Update(post *Post) error {
	return pg.db.Save(post).Error
}

func (pg *postGorm) Delete(id uint) error {
	post := &Post{Model: gorm.Model{ID: id}}
	return pg.db.Delete(post).Error
}
