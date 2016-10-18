package models

import (
	"time"
)

type UserTeamPublic struct {
	Id        string     `json:"id"`
	UserId    string     `db:"user_id" json:"userId"`
	TeamId    string     `db:"team_id" json:"teamId"`
	RequestId *string    `db:"request_id" json:"requestId"`
	IsLeader  bool       `db:"is_leader" json:"isLeader"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	LeftAt    *time.Time `db:"left_at" json:"leftAt"`
	KickedBy  *string    `db:"kicked_by" json:"kickedBy"`
}

type UserTeam struct {
	UserTeamPublic
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateUserTeam(userTeam *UserTeam) error {
	return e.Db.Get(
		userTeam, `
    INSERT INTO user_team (user_id, team_id, request_id, is_leader)
    VALUES ($1, $2, $3, $4)
    RETURNING *`,
		userTeam.UserId,
		userTeam.TeamId,
		userTeam.RequestId,
		userTeam.IsLeader,
	)
}

func (e *Env) GetUserTeams(modifier *QueryModifier) ([]UserTeam, error) {
	userTeams := make([]UserTeam, 0)
	sql, args, err := modifier.ToSql("user_team", "*")
	if err != nil {
		return userTeams, err
	}

	err = e.Db.Select(&userTeams, sql, args...)
	return userTeams, BetterGetterErrors(err)
}

func (e *Env) GetUserTeamById(id string) (*UserTeam, error) {
	var userTeam UserTeam
	err := e.Db.Get(
		&userTeam, `
    SELECT *
    FROM user_team
    WHERE id=$1`,
		id,
	)
	return &userTeam, BetterGetterErrors(err)
}

func (e *Env) GetUserTeamByUserTeam(userId, teamId string) (*UserTeam, error) {
	var userTeam UserTeam
	err := e.Db.Get(
		&userTeam, `
    SELECT *
    FROM user_team
    WHERE
      user_id=$1 AND
      team_id=$2 AND
      left_at IS NULL
    ORDER BY created_at DESC, id
    LIMIT 1`,
		userId,
		teamId,
	)
	return &userTeam, BetterGetterErrors(err)
}

func (e *Env) UpdateUserTeam(userTeam *UserTeam, updatedBy string) error {
	return e.Db.Get(
		userTeam, `
    UPDATE user_team
    SET
      is_leader=$2,
      left_at=$3,
      kicked_by=$4,
      updated_by=$5
    WHERE id=$1
    RETURNING *`,
		userTeam.Id,
		userTeam.IsLeader,
		userTeam.LeftAt,
		userTeam.KickedBy,
		updatedBy,
	)
}
