package models

import (
	"time"
)

type TeamSeasonRequest struct {
	Id          string     `json:"id"`
	TeamId      string     `db:"team_id" json:"teamId"`
	SeasonId    string     `db:"season_id" json:"seasonId"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy   string     `db:"created_by" json:"createdBy"`
	Decision    *bool      `db:"decision" json:"decision"`
	DecidedAt   *time.Time `db:"decided_at" json:"decidedAt"`
	DecidedBy   *string    `db:"decided_by" json:"decidedBy"`
	CancelledBy *string    `db:"cancelled_by" json:"cancelledBy"`
}

func (e *Env) CreateTeamSeasonRequest(
	teamSeasonRequest *TeamSeasonRequest,
) error {
	return e.Db.Get(
		teamSeasonRequest, `
    INSERT INTO team_season_request (
      team_id,
      season_id,
      created_by
    )
    VALUES ($1, $2, $3)
    RETURNING *`,
		teamSeasonRequest.TeamId,
		teamSeasonRequest.SeasonId,
		teamSeasonRequest.CreatedBy,
	)
}

func (e *Env) GetTeamSeasonRequestById(id string) (*TeamSeasonRequest, error) {
	var teamSeasonRequest TeamSeasonRequest
	err := e.Db.Get(
		&teamSeasonRequest, `
    SELECT *
    FROM team_season_request
    WHERE id=$1`,
		id,
	)
	return &teamSeasonRequest, BetterGetterErrors(err)
}

func (e *Env) GetTeamSeasonRequests(
	modifier *QueryModifier,
) ([]TeamSeasonRequest, error) {
	teamSeasonRequests := make([]TeamSeasonRequest, 0)

	sql, args, err := modifier.ToSql("team_season_request", "*")
	if err != nil {
		return teamSeasonRequests, err
	}

	err = e.Db.Select(&teamSeasonRequests, sql, args...)

	return teamSeasonRequests, err
}

func (e *Env) GetTeamSeasonRequestsByTeam(
	teamId string,
) ([]TeamSeasonRequest, error) {
	teamSeasonRequests := make([]TeamSeasonRequest, 0)
	err := e.Db.Select(
		&teamSeasonRequests, `
    SELECT *
    FROM team_season_request
    WHERE
      team_id=$1 AND
      decision IS NULL`,
		teamId,
	)
	return teamSeasonRequests, err
}

func (e *Env) GetTeamSeasonRequestsBySeason(
	seasonId string,
) ([]TeamSeasonRequest, error) {
	teamSeasonRequests := make([]TeamSeasonRequest, 0)
	err := e.Db.Select(
		&teamSeasonRequests, `
    SELECT *
    FROM team_season_request
    WHERE
      season_id=$1 AND
      decision IS NULL`,
		seasonId,
	)
	return teamSeasonRequests, err
}

func (e *Env) GetTeamSeasonRequestByTeamSeason(
	teamId, seasonId string,
) (*TeamSeasonRequest, error) {
	var teamSeasonRequest TeamSeasonRequest
	err := e.Db.Get(
		&teamSeasonRequest, `
    SELECT *
    FROM team_season_request
    WHERE
      team_id=$1 AND
      season_id=$2 AND
      decision IS NULL
    ORDER BY created_at DESC, id
    LIMIT 1`,
		teamId,
		seasonId,
	)
	return &teamSeasonRequest, BetterGetterErrors(err)
}

func (e *Env) UpdateTeamSeasonRequest(
	teamSeasonRequest *TeamSeasonRequest,
) error {
	return e.Db.Get(
		teamSeasonRequest, `
    UPDATE team_season_request
    SET
      decision=$2,
      decided_at=$3,
      decided_by=$4,
      cancelled_by=$5
    WHERE id=$1
    RETURNING *`,
		teamSeasonRequest.Id,
		teamSeasonRequest.Decision,
		teamSeasonRequest.DecidedAt,
		teamSeasonRequest.DecidedBy,
		teamSeasonRequest.CancelledBy,
	)
}
