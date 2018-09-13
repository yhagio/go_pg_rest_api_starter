package models

import (
	"strings"

	"github.com/jinzhu/gorm"
)

const (
	ErrNotFound               modelError   = "models: Resource not found"
	ErrInvalidID              modelError   = "models: ID provided was invalid"
	ErrInvalidEmailOrPassword modelError   = "models: Email or password is incorrect"
	ErrEmailRequired          modelError   = "models: Email is required"
	ErrEmailInvalid           modelError   = "models: Email is not valid"
	ErrEmailTaken             modelError   = "models: Email is already taken"
	ErrPasswordRequired       modelError   = "models: Password is required"
	ErrPasswordTooShort       modelError   = "models: Password must be at least 8 characters long"
	ErrTokenRequired          privateError = "models: Token is required"
	ErrTokenTooShort          privateError = "models: Token must be at least 32 bytes"
	ErrUserIDRequired         privateError = "models: UserID is required"
	ErrTokenInvalid           modelError   = "models: token provided is not valid"
	ErrServiceRequired        privateError = "models: Service is required"
	ErrTitleRequired          modelError   = "models: Title is required"
	ErrDescRequired           modelError   = "models: Description is required"
)

func First(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

type modelError string

func (e modelError) Error() string {
	return string(e)
}

// Creates public facing error message
func (e modelError) Public() string {
	cleanedStr := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(cleanedStr, " ")
	split[0] = strings.Title(split[0]) // Capitalize first letter
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
