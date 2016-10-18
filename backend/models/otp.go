package models

import (
	"time"

	"app/utils"
)

type OTP struct {
	TokenHash []byte    `db:"token" json:"-"`
	Token     string    `db:"-" json:"token"` // filled only upon creation
	UserId    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (e *Env) CreateOTP(userId string) (*OTP, error) {
	token, err := utils.GenerateToken(32, false)
	if err != nil {
		return nil, err
	}

	var otp OTP
	err = e.Db.Get(
		&otp, `
    INSERT INTO otp (token, user_id)
    VALUES ($1, $2)
    RETURNING *`,
		utils.Blake2b256(token),
		userId,
	)
	otp.Token = token
	return &otp, err
}

func (e *Env) DeleteOTPByToken(token string) (*OTP, error) {
	var otp OTP
	err := e.Db.Get(
		&otp, `
    DELETE FROM otp
    WHERE token=$1
    RETURNING *`,
		utils.Blake2b256(token),
	)
	return &otp, BetterGetterErrors(err)
}
