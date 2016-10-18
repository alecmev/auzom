package models

import (
	"time"
)

type BracketMapPublic struct {
	Id        string `json:"id"`
	BracketId string `db:"bracket_id" json:"bracketId"`
	GameMapId string `db:"game_map_id" json:"gameMapId"`
	SubPool   int    `db:"sub_pool" json:"subPool"`
	IsEnabled bool   `db:"is_enabled" json:"isEnabled"`
}

type BracketMap struct {
	BracketMapPublic
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy string     `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateBracketMap(bracketMap *BracketMap) error {
	return e.Db.Get(
		bracketMap, `
    INSERT INTO bracket_map (bracket_id, game_map_id, sub_pool, created_by)
    VALUES ($1, $2, $3, $4)
    RETURNING *`,
		bracketMap.BracketId,
		bracketMap.GameMapId,
		bracketMap.SubPool,
		bracketMap.CreatedBy,
	)
}

func (e *Env) GetBracketMapById(id string) (*BracketMap, error) {
	var bracketMap BracketMap
	err := e.Db.Get(
		&bracketMap, `
    SELECT *
    FROM bracket_map
    WHERE id=$1`,
		id,
	)
	return &bracketMap, BetterGetterErrors(err)
}

func (e *Env) GetBracketMaps(modifier *QueryModifier) ([]BracketMap, error) {
	bracketMaps := make([]BracketMap, 0)
	sql, args, err := modifier.ToSql("bracket_map", "*")
	if err != nil {
		return bracketMaps, err
	}

	err = e.Db.Select(&bracketMaps, sql, args...)
	return bracketMaps, err
}
