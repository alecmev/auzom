package models

import (
	"time"
)

type GamePublic struct {
	Id          string     `json:"id"`
	Slug        string     `json:"slug"`
	FranchiseId *string    `db:"franchise_id" json:"franchiseId"`
	Name        string     `json:"name"`
	Abbr        string     `json:"abbr"`
	ReleasedAt  *time.Time `db:"released_at" json:"releasedAt"`
	Cover       string     `json:"cover"`
	Summary     string     `json:"summary"`
}

type Game struct {
	GamePublic
	VerificationHandle *string    `db:"verification_handle" json:"verificationHandle"`
	CreatedAt          time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy          string     `db:"created_by" json:"createdBy"`
	UpdatedAt          *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy          *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateGame(game *Game) error {
	return e.Db.Get(
		game, `
    INSERT INTO game (
      slug,
      name,
      abbr,
      released_at,
      created_by,
      cover,
      summary,
      verification_handle
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING *`,
		game.Slug,
		game.Name,
		game.Abbr,
		game.ReleasedAt,
		game.CreatedBy,
		game.Cover,
		game.Summary,
		game.VerificationHandle,
	)
}

func (e *Env) GetGameById(id string) (*Game, error) {
	var game Game
	err := e.Db.Get(
		&game, `
    SELECT *
    FROM game
    WHERE id=$1`,
		id,
	)
	return &game, BetterGetterErrors(err)
}

func (e *Env) GetGameBySlug(slug string) (*Game, error) {
	var game Game
	err := e.Db.Get(
		&game, `
    SELECT *
    FROM game
    WHERE slug=$1`,
		slug,
	)
	return &game, BetterGetterErrors(err)
}

func (e *Env) GetGames(modifier *QueryModifier) ([]Game, error) {
	games := make([]Game, 0)
	sql, args, err := modifier.ToSql("game", "*")
	if err != nil {
		return games, err
	}

	err = e.Db.Select(&games, sql, args...)
	return games, err
}

func (e *Env) UpdateGame(game *Game, updatedBy string) error {
	return e.Db.Get(
		game, `
    UPDATE game
    SET
      slug=$2,
      name=$3,
      abbr=$4,
      released_at=$5,
      cover=$6,
      summary=$7,
      verification_handle=$8,
      updated_by=$9
    WHERE id=$1
    RETURNING *`,
		game.Id,
		game.Slug,
		game.Name,
		game.Abbr,
		game.ReleasedAt,
		game.Cover,
		game.Summary,
		game.VerificationHandle,
		updatedBy,
	)
}
