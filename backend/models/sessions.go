package models

import (
	"time"

	"app/utils"
)

type Session struct {
	TokenHash  []byte    `db:"token" json:"-"`
	Token      string    `db:"-" json:"token"` // filled only upon creation
	UserId     string    `db:"user_id" json:"userId"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	Remember   bool      `json:"remember"`
	LastUsedAt time.Time `db:"last_used_at" json:"lastUsedAt"`
	LastUsedIp string    `db:"last_used_ip" json:"lastUsedIp"`
}

func (e *Env) CreateSession(
	userId string, remember bool, ip string,
) (*Session, error) {
	token, err := utils.GenerateToken(32, false)
	if err != nil {
		return nil, err
	}

	var session Session
	err = e.Db.Get(
		&session, `
    INSERT INTO session (token, user_id, remember, last_used_ip)
    VALUES ($1, $2, $3, $4)
    RETURNING *`,
		utils.Blake2b256(token),
		userId,
		remember,
		ip,
	)
	session.Token = token
	return &session, err
}

func (e *Env) GetSessionByToken(token string) (*Session, error) {
	var session Session
	err := e.Db.Get(
		&session, `
    SELECT *
    FROM session
    WHERE token=$1`,
		utils.Blake2b256(token),
	)
	return &session, BetterGetterErrors(err)
}

func (e *Env) UpdateSessionLastUsed(session *Session) error {
	_, err := e.Db.Exec(`
    UPDATE session
    SET last_used_at=$1, last_used_ip=$2
    WHERE token=$3`,
		session.LastUsedAt,
		session.LastUsedIp,
		session.TokenHash,
	)
	return err
}

func (e *Env) DeleteSession(session *Session) error {
	_, err := e.Db.Exec(`
    DELETE FROM session
    WHERE token=$1`,
		session.TokenHash,
	)
	return err
}
