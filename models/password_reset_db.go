package models

import "github.com/jinzhu/gorm"

type passwordResetDB interface {
	GetOneByToken(token string) (*passwordReset, error)
	Create(pwr *passwordReset) error
	Delete(id uint) error
}

type passwordResetGorm struct {
	db *gorm.DB
}

func (pwrg *passwordResetGorm) GetOneByToken(tokenHash string) (*passwordReset, error) {
	var pwr passwordReset
	err := First(pwrg.db.Where("token_hash = ?", tokenHash), &pwr)
	if err != nil {
		return nil, err
	}
	return &pwr, nil
}

func (pwrg *passwordResetGorm) Create(pwr *passwordReset) error {
	return pwrg.db.Create(pwr).Error
}

func (pwrg *passwordResetGorm) Delete(id uint) error {
	pwr := passwordReset{
		Model: gorm.Model{ID: id},
	}
	return pwrg.db.Delete(&pwr).Error
}
