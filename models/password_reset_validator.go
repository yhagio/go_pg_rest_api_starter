package models

import "go_rest_pg_starter/utils"

func newPasswordResetValidator(db passwordResetDB, hmac utils.HMAC) *passwordResetValidator {
	return &passwordResetValidator{
		passwordResetDB: db,
		hmac:            hmac,
	}
}

type passwordResetValidator struct {
	passwordResetDB
	hmac utils.HMAC
}

func (pwrv *passwordResetValidator) requireUserID(pwr *passwordReset) error {
	if pwr.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (pwrv *passwordResetValidator) setTokenIfUnset(pwr *passwordReset) error {
	if pwr.Token != "" {
		return nil
	}
	token, err := utils.GenerateToken()
	if err != nil {
		return err
	}
	pwr.Token = token
	return nil
}

func (pwrv *passwordResetValidator) hmacToken(pwr *passwordReset) error {
	if pwr.Token == "" {
		return nil
	}
	pwr.TokenHash = pwrv.hmac.Hash(pwr.Token)
	return nil
}

type pwResetValFunc func(*passwordReset) error

func runPwResetValFuncs(pwr *passwordReset, fns ...pwResetValFunc) error {
	for _, fn := range fns {
		err := fn(pwr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pwrv *passwordResetValidator) GetOneByToken(token string) (*passwordReset, error) {
	pwr := passwordReset{Token: token}
	err := runPwResetValFuncs(&pwr, pwrv.hmacToken)
	if err != nil {
		return nil, err
	}
	return pwrv.passwordResetDB.GetOneByToken(pwr.TokenHash)
}

func (pwrv *passwordResetValidator) Create(pwr *passwordReset) error {
	err := runPwResetValFuncs(pwr,
		pwrv.requireUserID,
		pwrv.setTokenIfUnset,
		pwrv.hmacToken,
	)
	if err != nil {
		return err
	}
	return pwrv.passwordResetDB.Create(pwr)
}

func (pwrv *passwordResetValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrInvalidID
	}
	return pwrv.passwordResetDB.Delete(id)
}
