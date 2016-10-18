package models

import (
	"time"
)

type UserTeamRequest struct {
	Id        string    `json:"id"`
	UserId    string    `db:"user_id" json:"userId"`
	TeamId    string    `db:"team_id" json:"teamId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`

	Decision  *bool      `db:"decision" json:"decision"`
	DecidedAt *time.Time `db:"decided_at" json:"decidedAt"`

	UserDecision  *bool      `db:"user_decision" json:"userDecision"`
	UserDecidedAt *time.Time `db:"user_decided_at" json:"userDecidedAt"`

	LeaderDecision  *bool      `db:"leader_decision" json:"leaderDecision"`
	LeaderDecidedAt *time.Time `db:"leader_decided_at" json:"leaderDecidedAt"`
	LeaderDecidedBy *string    `db:"leader_decided_by" json:"leaderDecidedBy"`

	AdminDecision  *bool      `db:"admin_decision" json:"adminDecision"`
	AdminDecidedAt *time.Time `db:"admin_decided_at" json:"adminDecidedAt"`
	AdminDecidedBy *string    `db:"admin_decided_by" json:"adminDecidedBy"`
}

func (userTeamRequest *UserTeamRequest) IsAdminNeeded() bool {
	return !(userTeamRequest.AdminDecision != nil &&
		userTeamRequest.AdminDecidedAt == nil)
}

// TODO: think through what needs to be declassified

func (e *Env) CreateUserTeamRequest(userTeamRequest *UserTeamRequest) error {
	return e.Db.Get(
		userTeamRequest, `
    INSERT INTO user_team_request (
      user_id,
      team_id,
      user_decision,
      leader_decision,
      leader_decided_by,
      admin_decision
    )
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *`,
		userTeamRequest.UserId,
		userTeamRequest.TeamId,
		userTeamRequest.UserDecision,
		userTeamRequest.LeaderDecision,
		userTeamRequest.LeaderDecidedBy,
		userTeamRequest.AdminDecision,
	)
}

func (e *Env) GetUserTeamRequests(modifier *QueryModifier) ([]UserTeamRequest, error) {
	userTeamRequests := make([]UserTeamRequest, 0)

	sql, args, err := modifier.ToSql("user_team_request", "*")
	if err != nil {
		return userTeamRequests, err
	}

	err = e.Db.Select(&userTeamRequests, sql, args...)

	return userTeamRequests, err
}

func (e *Env) GetUserTeamRequestsByUser(
	userId string,
) ([]UserTeamRequest, error) {
	userTeamRequests := make([]UserTeamRequest, 0)
	err := e.Db.Select(
		&userTeamRequests, `
    SELECT *
    FROM user_team_request
    WHERE
      user_id=$1 AND
      decision IS NULL`,
		userId,
	)
	return userTeamRequests, err
}

func (e *Env) GetUserTeamRequestsByTeam(
	teamId string,
) ([]UserTeamRequest, error) {
	userTeamRequests := make([]UserTeamRequest, 0)
	err := e.Db.Select(
		&userTeamRequests, `
    SELECT *
    FROM user_team_request
    WHERE
      team_id=$1 AND
      decision IS NULL`,
		teamId,
	)
	return userTeamRequests, err
}

func (e *Env) GetUserTeamRequestById(id string) (*UserTeamRequest, error) {
	var userTeamRequest UserTeamRequest
	err := e.Db.Get(
		&userTeamRequest, `
    SELECT *
    FROM user_team_request
    WHERE id=$1`,
		id,
	)
	return &userTeamRequest, BetterGetterErrors(err)
}

func (e *Env) GetUserTeamRequestByUserTeam(
	userId, teamId string,
) (*UserTeamRequest, error) {
	var userTeamRequest UserTeamRequest
	err := e.Db.Get(
		&userTeamRequest, `
    SELECT *
    FROM user_team_request
    WHERE
      user_id=$1 AND
      team_id=$2 AND
      decision IS NULL
    ORDER BY created_at DESC
    LIMIT 1`,
		userId,
		teamId,
	)
	return &userTeamRequest, BetterGetterErrors(err)
}

func (e *Env) UpdateUserTeamRequest(userTeamRequest *UserTeamRequest) error {
	return e.Db.Get(
		userTeamRequest, `
    UPDATE user_team_request
    SET
      decision=$1,
      decided_at=$2,
      user_decision=$3,
      user_decided_at=$4,
      leader_decision=$5,
      leader_decided_at=$6,
      leader_decided_by=$7,
      admin_decision=$8,
      admin_decided_at=$9,
      admin_decided_by=$10
    WHERE id=$11
    RETURNING *`,
		userTeamRequest.Decision,
		userTeamRequest.DecidedAt,
		userTeamRequest.UserDecision,
		userTeamRequest.UserDecidedAt,
		userTeamRequest.LeaderDecision,
		userTeamRequest.LeaderDecidedAt,
		userTeamRequest.LeaderDecidedBy,
		userTeamRequest.AdminDecision,
		userTeamRequest.AdminDecidedAt,
		userTeamRequest.AdminDecidedBy,
		userTeamRequest.Id,
	)
}
