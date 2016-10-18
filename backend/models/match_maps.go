package models

import (
	"time"
)

type MatchMapPublic struct {
	Id          string     `json:"id"`
	MatchId     string     `db:"match_id" json:"matchId"`
	GameMapId   string     `db:"game_map_id" json:"gameMapId"`
	TeamId      *string    `db:"team_id" json:"teamId"`
	IsBan       bool       `db:"is_ban" json:"isBan"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	DiscardedAt *time.Time `db:"discarded_at" json:"discardedAt"`
}

type MatchMap struct {
	MatchMapPublic
	CreatedBy   string  `db:"created_by" json:"createdBy"`
	DiscardedBy *string `db:"discarded_by" json:"discardedBy"`
}

func (e *Env) CreateMatchMap(matchMap *MatchMap) error {
	return e.Db.Get(
		matchMap, `
    INSERT INTO match_map (match_id, game_map_id, team_id, is_ban, created_by)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *`,
		matchMap.MatchId,
		matchMap.GameMapId,
		matchMap.TeamId,
		matchMap.IsBan,
		matchMap.CreatedBy,
	)
}

func (e *Env) GetMatchMapById(id string) (*MatchMap, error) {
	var matchMap MatchMap
	err := e.Db.Get(
		&matchMap, `
    SELECT *
    FROM match_map
    WHERE id=$1`,
		id,
	)
	return &matchMap, BetterGetterErrors(err)
}

func (e *Env) GetMatchMaps(modifier *QueryModifier) ([]MatchMap, error) {
	matchMaps := make([]MatchMap, 0)
	sql, args, err := modifier.ToSql("match_map", "*")
	if err != nil {
		return matchMaps, err
	}

	err = e.Db.Select(&matchMaps, sql, args...)
	return matchMaps, err
}

func (e *Env) UpdateMatchMap(matchMap *MatchMap) error {
	err := e.Db.Get(
		matchMap, `
    UPDATE match_map
    SET discarded_at=$2, discarded_by=$3
    WHERE id=$1
    RETURNING *`,
		matchMap.Id,
		matchMap.DiscardedAt,
		matchMap.DiscardedBy,
	)
	return err
}
