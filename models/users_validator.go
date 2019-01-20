package models

import (
	"errors"
	"regexp"
	"strings"

	"go_rest_pg_starter/utils"

	"golang.org/x/crypto/bcrypt"
)

// userValidator is our validation layer that validates
// and normalizes data before passing it on to the next
// UserDB in our interface chain.
type userValidator struct {
	UserDB
	hmac       utils.HMAC
	pepper     string
	emailRegex *regexp.Regexp
}

func newUserValidator(udb UserDB, hmac utils.HMAC, pepper string) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		pepper:     pepper,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

func (uv *userValidator) GetById(id uint) (*User, error) {
	if id <= 0 {
		return nil, errors.New("Invalid ID")
	}
	return uv.UserDB.GetById(id)
}

func (uv *userValidator) GetByToken(token string) (*User, error) {
	user := User{Token: token}
	err := userValidationFuncs(&user, uv.hmacHashToken)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.GetByToken(user.TokenHash)
}

func (uv *userValidator) GetByEmail(email string) (*User, error) {
	user := User{Email: email}
	err := userValidationFuncs(&user, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.GetByEmail(user.Email)
}

func (uv *userValidator) Create(user *User) error {
	err := userValidationFuncs(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.generatePasswordHash,
		uv.passwordHashRequired,
		uv.setTokenIfNotSet,
		uv.tokenMinBytes,
		uv.hmacHashToken,
		uv.tokenHashRequired,
		uv.requireEmail,
		uv.normalizeEmail,
		uv.checkEmailFormat,
		uv.checkEmailAvailability)
	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	err := userValidationFuncs(user,
		uv.passwordMinLength,
		uv.generatePasswordHash,
		uv.passwordHashRequired,
		uv.tokenMinBytes,
		uv.hmacHashToken,
		uv.tokenHashRequired,
		uv.requireEmail,
		uv.normalizeEmail,
		uv.checkEmailFormat,
		uv.checkEmailAvailability)
	if err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := userValidationFuncs(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}

	return uv.UserDB.Delete(id)
}

func (uv *userValidator) generatePasswordHash(user *User) error {
	// If password is not changed, do nothing
	if user.Password == "" {
		return nil
	}

	hasedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password+uv.pepper),
		bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.PasswordHash = string(hasedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) setTokenIfNotSet(user *User) error {
	if user.Token != "" {
		return nil
	}

	token, err := utils.GenerateToken()
	if err != nil {
		return err
	}
	user.Token = token

	return nil
}

func (uv *userValidator) hmacHashToken(user *User) error {
	if user.Token == "" {
		return nil
	}
	user.TokenHash = uv.hmac.Hash(user.Token)
	return nil
}

// Closure way, dynamically validates with argument
func (uv *userValidator) idGreaterThan(num uint) userValidationFunc {
	return userValidationFunc(func(user *User) error {
		if user.ID <= num {
			return ErrInvalidID
		}
		return nil
	})
}

func (uv *userValidator) tokenMinBytes(user *User) error {
	if user.Token == "" {
		return nil
	}
	num, err := utils.NumberOfBytes(user.Token)
	if err != nil {
		return err
	}
	if num < 32 {
		return ErrTokenTooShort
	}
	return nil
}

func (uv *userValidator) tokenHashRequired(user *User) error {
	if user.TokenHash == "" {
		return ErrTokenRequired
	}
	return nil
}

///////////////////////////////////////////////////////////
// Eamil validation
///////////////////////////////////////////////////////////

// Normalize Email
func (uv *userValidator) normalizeEmail(user *User) error {
	// trim space and make it lowercase
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) checkEmailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) checkEmailAvailability(user *User) error {
	existing, err := uv.GetByEmail(user.Email)
	if err == ErrNotFound {
		// Email is available if the email is not found in our db
		return nil
	}
	if err != nil {
		return err
	}

	// Email is found but user id is different
	if user.ID != existing.ID {
		return ErrEmailTaken
	}
	return nil
}

///////////////////////////////////////////////////////////
// Password validation
///////////////////////////////////////////////////////////

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}

///////////////////////////////////////////////////////////
// Reusable validation functions helper
///////////////////////////////////////////////////////////

type userValidationFunc func(*User) error

func userValidationFuncs(user *User, funcs ...userValidationFunc) error {
	for _, fn := range funcs {
		err := fn(user)
		if err != nil {
			return err
		}
	}
	return nil
}
