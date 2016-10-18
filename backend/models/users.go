package models

import (
	"strings"
	"time"

	"app/utils"
)

type UserPublic struct {
	Id        string    `json:"id"`
	Nickname  string    `json:"nickname"`
	Fullname  string    `json:"fullname"`
	IsAdmin   bool      `db:"is_admin" json:"isAdmin"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Gravatar  string    `json:"gravatar"`
}

type User struct {
	UserPublic
	Email           string     `json:"email"`
	Password        []byte     `json:"-"`
	IsEmailVerified bool       `db:"is_email_verified" json:"isEmailVerified"`
	GravatarEmail   string     `db:"gravatar_email" json:"gravatarEmail"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy       *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateUser(user *User) error {
	return e.Db.Get(
		user, `
    INSERT INTO "user" (
      email, password, is_email_verified, nickname, is_admin, gravatar
    )
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *`,
		user.Email,
		user.Password,
		user.IsEmailVerified,
		strings.Split(user.Email, "@")[0],
		user.IsAdmin,
		utils.EmailToGravatar(user.Email),
	)
}

func (e *Env) GetUsers(modifier *QueryModifier) ([]User, error) {
	users := make([]User, 0)
	sql, args, err := modifier.ToSql("\"user\"",
		"id", "nickname", "fullname", "is_admin", "created_at", "gravatar")
	if err != nil {
		return users, err
	}

	err = e.Db.Select(&users, sql, args...)
	return users, err
}

func (e *Env) GetUserById(id string) (*User, error) {
	var user User
	err := e.Db.Get(
		&user, `
    SELECT *
    FROM "user"
    WHERE id=$1`,
		id,
	)
	return &user, BetterGetterErrors(err)
}

// http://stackoverflow.com/q/9807909/242684
func (e *Env) GetUserByEmail(email string) (*User, error) {
	var user User
	err := e.Db.Get(
		&user, `
    SELECT *
    FROM "user"
    WHERE LOWER(email)=LOWER($1)`,
		email,
	)
	return &user, BetterGetterErrors(err)
}

func (e *Env) UpdateUser(user *User, updatedBy string) error {
	return e.Db.Get(
		user, `
    UPDATE "user"
    SET
      nickname=$2,
      fullname=$3,
      is_admin=$4,
      gravatar=$5,
      email=$6,
      password=$7,
      is_email_verified=$8,
      gravatar_email=$9,
      updated_by=$10
    WHERE id=$1
    RETURNING *`,
		user.Id,
		user.Nickname,
		user.Fullname,
		user.IsAdmin,
		utils.EmailToGravatar(user.GravatarEmail),
		user.Email,
		user.Password,
		user.IsEmailVerified,
		user.GravatarEmail,
		updatedBy,
	)
}
