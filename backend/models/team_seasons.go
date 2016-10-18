package models

import (
	"time"
)

type TeamSeasonPublic struct {
	Id        string     `json:"id"`
	TeamId    string     `db:"team_id" json:"teamId"`
	SeasonId  string     `db:"season_id" json:"seasonId"`
	RequestId *string    `db:"request_id" json:"requestId"`
	LeftAt    *time.Time `db:"left_at" json:"leftAt"`
	IsDone    bool       `db:"is_done" json:"isDone"`
	KickedBy  *string    `db:"kicked_by" json:"kickedBy"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
}

type TeamSeason struct {
	TeamSeasonPublic
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateTeamSeason(teamSeason *TeamSeason) error {
	return e.Db.Get(
		teamSeason, `
    INSERT INTO team_season (team_id, season_id, request_id)
    VALUES ($1, $2, $3)
    RETURNING *`,
		teamSeason.TeamId,
		teamSeason.SeasonId,
		teamSeason.RequestId,
	)
}

func (e *Env) GetTeamSeasons(modifier *QueryModifier) ([]TeamSeason, error) {
	teamSeasons := make([]TeamSeason, 0)
	sql, args, err := modifier.ToSql("team_season", "*")
	if err != nil {
		return teamSeasons, err
	}

	err = e.Db.Select(&teamSeasons, sql, args...)
	return teamSeasons, BetterGetterErrors(err)
}

func (e *Env) GetTeamSeasonsBySeason(season *Season) ([]TeamSeason, error) {
	teamSeasons := make([]TeamSeason, 0)
	err := e.Db.Select(
		&teamSeasons, `
    SELECT *
    FROM team_season
    WHERE season_id=$1`,
		season.Id,
	)
	return teamSeasons, BetterGetterErrors(err)
}

func (e *Env) GetTeamSeasonById(id string) (*TeamSeason, error) {
	var teamSeason TeamSeason
	err := e.Db.Get(
		&teamSeason, `
    SELECT *
    FROM team_season
    WHERE id=$1`,
		id,
	)
	return &teamSeason, BetterGetterErrors(err)
}

func (e *Env) GetTeamSeasonByTeamSeason(
	teamId, seasonId string,
) (*TeamSeason, error) {
	var teamSeason TeamSeason
	err := e.Db.Get(
		&teamSeason, `
    SELECT *
    FROM team_season
    WHERE
      team_id=$1 AND
      season_id=$2 AND
      left_at IS NULL
    ORDER BY created_at DESC, id
    LIMIT 1`,
		teamId,
		seasonId,
	)
	return &teamSeason, BetterGetterErrors(err)
}

func (e *Env) UpdateTeamSeason(teamSeason *TeamSeason, updatedBy string) error {
	return e.Db.Get(
		teamSeason, `
    UPDATE team_season
    SET
      left_at=$2,
      is_done=$3,
      kicked_by=$4,
      updated_by=$5
    WHERE id=$1
    RETURNING *`,
		teamSeason.Id,
		teamSeason.LeftAt,
		teamSeason.IsDone,
		teamSeason.KickedBy,
		updatedBy,
	)
}
