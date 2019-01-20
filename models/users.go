package models

import (
	"time"

	"go_rest_pg_starter/utils"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username     string `gorm:"not null; unique_index"`
	Email        string `gorm:"not null; unique_index"`
	Role         string `gorm:"not null; default:'standard'"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null`
	Token        string `gorm:"-"`
	TokenHash    string `gorm:"not null; unique_index"`
}

// UserService is a set of methods used to manipulate and
// work with the user model
type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
	InitiateReset(email string) (string, error)
	CompleteReset(token, newPassword string) (*User, error)
}

func NewUserService(db *gorm.DB, pepper, hmacKey string) UserService {
	ug := &userGorm{db}
	hmac := utils.NewHMAC(hmacKey)
	uv := newUserValidator(ug, hmac, pepper)

	return &userService{
		UserDB:          uv,
		pepper:          pepper,
		passwordResetDB: newPasswordResetValidator(&passwordResetGorm{db}, hmac),
	}
}

var _ UserService = &userService{}

type userService struct {
	UserDB
	pepper          string
	passwordResetDB passwordResetDB
}

// Authenticate user. Checks email and password.
func (us *userService) Authenticate(email, password string) (*User, error) {
	user, err := us.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password+us.pepper),
	)

	if err != nil {
		return nil, ErrInvalidEmailOrPassword
	}

	return user, nil
}

// InitiateReset will complete all the model-related tasks
// to start the password reset process for the user with
// the provided email address. Once completed, it will
// return the token, or an error if there was one.
func (us *userService) InitiateReset(email string) (string, error) {
	user, err := us.GetByEmail(email)
	if err != nil {
		return "", err
	}

	pwr := passwordReset{
		UserID: user.ID,
	}

	err = us.passwordResetDB.Create(&pwr)
	if err != nil {
		return "", err
	}
	return pwr.Token, nil
}

// CompleteReset will complete all the model-related tasks
// to complete the password reset process for the user that
// the token matches, including updating that user's pw.
// If the token has expired, or if it is invalid for any
// other reason the ErrTokenInvalid error will be returned.
func (us *userService) CompleteReset(token, newPw string) (*User, error) {
	pwr, err := us.passwordResetDB.GetOneByToken(token)
	if err != nil {
		if err == ErrNotFound {
			return nil, ErrTokenInvalid
		}
		return nil, err
	}

	// If the password rest is over 12 hours old, it is invalid
	if time.Now().Sub(pwr.CreatedAt) > (12 * time.Hour) {
		return nil, ErrTokenInvalid
	}

	user, err := us.GetById(pwr.UserID)
	if err != nil {
		return nil, err
	}

	user.Password = newPw
	err = us.Update(user)
	if err != nil {
		return nil, err
	}

	us.passwordResetDB.Delete(pwr.ID)
	return user, nil
}
