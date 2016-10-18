package models

import (
	"github.com/jmoiron/sqlx"
)

type UserTeamRelation struct {
	UserId    string `db:"user_id"`
	TeamId    string `db:"team_id"`
	IsRequest bool   `db:"is_request"`
}

func (e *Env) GetUserTeamRelationsByUser(user *User) (
	[]UserTeamRelation, error,
) {
	relations := make([]UserTeamRelation, 0)
	err := e.Db.Select(
		&relations, `
    SELECT user_id, team_id, FALSE AS is_request
    FROM user_team
    WHERE user_id=$1 AND left_at IS NULL
    UNION ALL
    SELECT user_id, team_id, TRUE AS is_request
    FROM user_team_request
    WHERE user_id=$1 AND decision IS NULL`,
		user.Id,
	)
	return relations, err
}

func (e *Env) GetUserTeamRelationsByTeam(team *Team) (
	[]UserTeamRelation, error,
) {
	_, ok := e.Db.(*sqlx.Tx)
	if ok {
		return getUserTeamRelationsByTeam(e.Db, team.Id)
	}

	var relations []UserTeamRelation
	err := e.Atomic(func(etx *Env) error {
		var inerr error
		relations, inerr = getUserTeamRelationsByTeam(etx.Db, team.Id)
		return inerr
	})
	return relations, err
}

func getUserTeamRelationsByTeam(db database, teamId string) (
	[]UserTeamRelation, error,
) {
	relations := make([]UserTeamRelation, 0)
	_, err := db.Exec(`
    CREATE TEMP TABLE tmp ON COMMIT DROP AS
    SELECT user_id
    FROM user_team
    WHERE team_id=$1 AND left_at IS NULL
    UNION ALL
    SELECT user_id
    FROM user_team_request
    WHERE team_id=$1 AND decision IS NULL`,
		teamId,
	)
	if err != nil {
		return nil, err
	}

	err = db.Select(
		&relations, `
    SELECT ut.user_id, ut.team_id, FALSE AS is_request
    FROM user_team ut
    INNER JOIN tmp
    ON ut.user_id=tmp.user_id
    WHERE ut.left_at IS NULL AND ut.team_id!=$1
    UNION ALL
    SELECT utr.user_id, utr.team_id, TRUE AS is_request
    FROM user_team_request utr
    INNER JOIN tmp
    ON utr.user_id=tmp.user_id
    WHERE utr.decision IS NULL AND utr.team_id!=$1`,
		teamId,
	)
	return relations, err
}
