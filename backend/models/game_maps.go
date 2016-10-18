package models

import (
	"time"
)

type GameMapPublic struct {
	Id        string `json:"id"`
	GameId    string `db:"game_id" json:"gameId"`
	Name      string `json:"name"`
	Abbr      string `json:"abbr"`
	SideX     string `db:"side_x" json:"sideX"`
	SideXAbbr string `db:"side_x_abbr" json:"sideXAbbr"`
	SideY     string `db:"side_y" json:"sideY"`
	SideYAbbr string `db:"side_y_abbr" json:"sideYAbbr"`
}

type GameMap struct {
	GameMapPublic
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy string     `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateGameMap(gameMap *GameMap) error {
	return e.Db.Get(
		gameMap, `
    INSERT INTO game_map (
      game_id, name, abbr, side_x, side_x_abbr, side_y, side_y_abbr, created_by
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING *`,
		gameMap.GameId,
		gameMap.Name,
		gameMap.Abbr,
		gameMap.SideX,
		gameMap.SideXAbbr,
		gameMap.SideY,
		gameMap.SideYAbbr,
		gameMap.CreatedBy,
	)
}

func (e *Env) GetGameMapById(id string) (*GameMap, error) {
	var gameMap GameMap
	err := e.Db.Get(
		&gameMap, `
    SELECT *
    FROM game_map
    WHERE id=$1`,
		id,
	)
	return &gameMap, BetterGetterErrors(err)
}

func (e *Env) GetGameMaps(modifier *QueryModifier) ([]GameMap, error) {
	gameMaps := make([]GameMap, 0)
	sql, args, err := modifier.ToSql("game_map", "*")
	if err != nil {
		return gameMaps, err
	}

	err = e.Db.Select(&gameMaps, sql, args...)
	return gameMaps, err
}

func (e *Env) UpdateGameMap(gameMap *GameMap, updatedBy string) error {
	return e.Db.Get(
		gameMap, `
    UPDATE game_map
    SET
      name=$2,
      abbr=$3,
      side_x=$4,
      side_x_abbr=$5,
      side_y=$6,
      side_y_abbr=$7,
      updated_by=$8
    WHERE id=$1
    RETURNING *`,
		gameMap.Id,
		gameMap.Name,
		gameMap.Abbr,
		gameMap.SideX,
		gameMap.SideXAbbr,
		gameMap.SideY,
		gameMap.SideYAbbr,
		updatedBy,
	)
}
