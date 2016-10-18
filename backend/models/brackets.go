package models

import (
	"time"
)

type BracketPublic struct {
	Id               string `json:"id"`
	StageId          string `db:"stage_id" json:"stageId"`
	Slug             string `json:"slug"`
	Name             string `json:"name"`
	Abbr             string `json:"abbr"`
	Order            int    `json:"order"`
	Type             string `json:"type"`
	Size             int    `json:"size"`
	MapVetoProcedure string `db:"map_veto_procedure" json:"mapVetoProcedure"`
}

type Bracket struct {
	BracketPublic
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy string     `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateBracket(bracket *Bracket) error {
	return e.Db.Get(
		bracket, `
    INSERT INTO bracket (
      stage_id, slug, name, abbr, "order", type, size, map_veto_procedure,
      created_by
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING *`,
		bracket.StageId,
		bracket.Slug,
		bracket.Name,
		bracket.Abbr,
		bracket.Order,
		bracket.Type,
		bracket.Size,
		bracket.MapVetoProcedure,
		bracket.CreatedBy,
	)
}

func (e *Env) GetBracketById(id string) (*Bracket, error) {
	var bracket Bracket
	err := e.Db.Get(
		&bracket, `
    SELECT *
    FROM bracket
    WHERE id=$1`,
		id,
	)
	return &bracket, BetterGetterErrors(err)
}

func (e *Env) GetMatchesForBracket(bracketId string) ([]Match, error) {
	matches := make([]Match, 0)

	err := e.Db.Select(
		&matches, `
    SELECT match.*
    FROM bracket, match
	  WHERE bracket.id=match.bracket_id AND bracket.id=$1`,
		bracketId,
	)

	return matches, err
}

// Note: currently assumes single elimination
// Maybe the root node(s) can be marked in the db instead?
func (e *Env) GetFinalMatchForBracket(bracketId string) (*Match, error) {
	var match Match

	err := e.Db.Get(
		&match, `
    SELECT m.*
    FROM bracket, match m
    WHERE bracket.id=m.bracket_id AND bracket.id=$1 AND (
      SELECT count(*)
      FROM match m2
      WHERE m2.parent_x=m.id OR m2.parent_y=m.Id
    )=0`,
		bracketId,
	)

	return &match, err
}

func (e *Env) GetBrackets(modifier *QueryModifier) ([]Bracket, error) {
	brackets := make([]Bracket, 0)

	sql, args, err := modifier.ToSql("bracket", "*")
	if err != nil {
		return brackets, err
	}

	err = e.Db.Select(&brackets, sql, args...)

	return brackets, err
}

func (e *Env) UpdateBracket(bracket *Bracket, updatedBy string) error {
	err := e.Db.Get(
		bracket, `
    UPDATE bracket
    SET
      slug=$2,
      name=$3,
      abbr=$4,
      "order"=$5,
      map_veto_procedure=$6,
      updated_by=$7
    WHERE id=$1
    RETURNING *`,
		bracket.Id,
		bracket.Slug,
		bracket.Name,
		bracket.Abbr,
		bracket.Order,
		bracket.MapVetoProcedure,
		updatedBy,
	)
	return err
}
