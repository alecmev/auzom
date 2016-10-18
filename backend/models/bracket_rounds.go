package models

import (
	"time"
)

type BracketRoundPublic struct {
	Id               string  `json:"id"`
	BracketId        string  `db:"bracket_id" json:"bracketId"`
	Number           int     `json:"number"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ByeTeamId        *string `db:"bye_team_id" json:"byeTeamId"`
	MapVetoProcedure string  `db:"map_veto_procedure" json:"mapVetoProcedure"`
}

type BracketRound struct {
	BracketRoundPublic
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy string     `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateBracketRound(bracketRound *BracketRound) error {
	return e.Db.Get(
		bracketRound, `
    INSERT INTO bracket_round (
      bracket_id, number, name, description, bye_team_id, map_veto_procedure,
      created_by
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING *`,
		bracketRound.BracketId,
		bracketRound.Number,
		bracketRound.Name,
		bracketRound.Description,
		bracketRound.ByeTeamId,
		bracketRound.MapVetoProcedure,
		bracketRound.CreatedBy,
	)
}

func (e *Env) GetBracketRoundById(id string) (*BracketRound, error) {
	var bracketRound BracketRound
	err := e.Db.Get(
		&bracketRound, `
    SELECT *
    FROM bracket_round
    WHERE id=$1`,
		id,
	)
	return &bracketRound, BetterGetterErrors(err)
}

func (e *Env) GetBracketRoundByMatch(match *Match) (*BracketRound, error) {
	var bracketRound BracketRound
	err := e.Db.Get(
		&bracketRound, `
    SELECT *
    FROM bracket_round
    WHERE bracket_id=$1 AND number=$2`,
		match.BracketId,
		match.BracketRound,
	)
	return &bracketRound, BetterGetterErrors(err)
}

func (e *Env) GetBracketRounds(
	modifier *QueryModifier,
) ([]BracketRound, error) {
	brackets := make([]BracketRound, 0)
	sql, args, err := modifier.ToSql("bracket_round", "*")
	if err != nil {
		return brackets, err
	}

	err = e.Db.Select(&brackets, sql, args...)
	return brackets, err
}

func (e *Env) UpdateBracketRound(
	bracketRound *BracketRound, updatedBy string,
) error {
	return e.Db.Get(
		bracketRound, `
    UPDATE bracket_round
    SET
      name=$2,
      description=$3,
      map_veto_procedure=$4,
      updated_by=$5
    WHERE id=$1
    RETURNING *`,
		bracketRound.Id,
		bracketRound.Name,
		bracketRound.Description,
		bracketRound.MapVetoProcedure,
		updatedBy,
	)
}
