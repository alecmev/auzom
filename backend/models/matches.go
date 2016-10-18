package models

import (
	"errors"
	"net/http"
	"time"

	"app/utils"
)

type MatchPublic struct {
	Id                string     `json:"id"`
	BracketId         string     `db:"bracket_id" json:"bracketId"`
	BracketRound      int        `db:"bracket_round" json:"bracketRound"`
	SortNumber        int        `db:"sort_number" json:"sortNumber"`
	AreMapsReady      bool       `db:"are_maps_ready" json:"areMapsReady"`
	StartedAt         time.Time  `db:"started_at" json:"startedAt"`
	ReportingClosedAt *time.Time `db:"reporting_closed_at" json:"reportingClosedAt"`

	SeedX          *int    `db:"seed_x" json:"seedX"`
	SeedY          *int    `db:"seed_y" json:"seedY"`
	TeamX          *string `db:"team_x" json:"teamX"`
	TeamY          *string `db:"team_y" json:"teamY"`
	ParentX        *string `db:"parent_x" json:"parentMatchX"`
	ParentXIsLoser bool    `db:"parent_x_is_loser" json:"parentMatchXIsLoser"`
	ParentY        *string `db:"parent_y" json:"parentMatchY"`
	ParentYIsLoser bool    `db:"parent_y_is_loser" json:"parentMatchYIsLoser"`

	MatchReportId *string  `db:"match_report_id" json:"matchReportId"`
	ScoreX        *float64 `db:"score_x" json:"scoreX"`
	ScoreY        *float64 `db:"score_y" json:"scoreY"`
	RawScoreX     *float64 `db:"raw_score_x" json:"rawScoreX"`
	RawScoreY     *float64 `db:"raw_score_y" json:"rawScoreY"`
	IsOverridden  bool     `db:"is_overridden" json:"isOverridden"`
	IsPenalized   bool     `db:"is_penalized" json:"isPenalized"`
}

type Match struct {
	MatchPublic
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy string     `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateMatch(match *Match) error {
	return e.Db.Get(
		match, `
    INSERT INTO match (
      bracket_id, bracket_round, sort_number, are_maps_ready, started_at,
      reporting_closed_at,

      seed_x, seed_y,
      team_x, team_y,
      parent_x, parent_x_is_loser,
      parent_y, parent_y_is_loser,

      created_by
    )
    VALUES (
      $1, $2, $3, $4, $5, $6,
      $7, $8, $9, $10, $11, $12, $13, $14,
      $15
    )
    RETURNING *`,
		match.BracketId,
		match.BracketRound,
		match.SortNumber,
		match.AreMapsReady,
		match.StartedAt,
		match.ReportingClosedAt,

		match.SeedX,
		match.SeedY,
		match.TeamX,
		match.TeamY,
		match.ParentX,
		match.ParentXIsLoser,
		match.ParentY,
		match.ParentYIsLoser,

		match.CreatedBy,
	)
}

func (e *Env) GetMatchById(id string) (*Match, error) {
	var match Match
	err := e.Db.Get(
		&match, `
    SELECT *
    FROM match
    WHERE id=$1`,
		id,
	)
	return &match, BetterGetterErrors(err)
}

func (e *Env) GetMatches(modifier *QueryModifier) ([]Match, error) {
	matches := make([]Match, 0)
	sql, args, err := modifier.ToSql("match", "*")
	if err != nil {
		return matches, err
	}

	err = e.Db.Select(&matches, sql, args...)
	return matches, err
}

func (e *Env) GetChildMatches(parent *Match) ([]Match, error) {
	matches := make([]Match, 0)
	err := e.Db.Select(
		&matches, `
    SELECT *
    FROM match
    WHERE parent_x=$1 OR parent_y=$1`,
		parent.Id,
	)
	return matches, err
}

func (e *Env) UpdateMatch(match *Match, updatedBy string) error {
	return e.Db.Get(
		match, `
    UPDATE match
    SET
      are_maps_ready=$2,
      started_at=$3,

      team_x=$4,
      team_y=$5,

      match_report_id=$6,
      score_x=$7,
      score_y=$8,
      raw_score_x=$9,
      raw_score_y=$10,
      is_overridden=$11,
      is_penalized=$12,

      reporting_closed_at=$13,

      seed_x=$14,
      seed_y=$15,

      updated_by=$16
    WHERE id=$1
    RETURNING *`,
		match.Id,

		match.AreMapsReady,
		match.StartedAt,

		match.TeamX,
		match.TeamY,

		match.MatchReportId,
		match.ScoreX,
		match.ScoreY,
		match.RawScoreX,
		match.RawScoreY,
		match.IsOverridden,
		match.IsPenalized,

		match.ReportingClosedAt,

		match.SeedX,
		match.SeedY,

		updatedBy,
	)
}

var ErrNotSeededYet = errors.New("match not seeded yet")

func (match *Match) UserIsLeaderOf(
	e *Env, userId string,
) (map[string]struct{}, error, int) {
	res := map[string]struct{}{}
	checkLeadership := func(teamId *string) (error, int) {
		if teamId == nil {
			return ErrNotSeededYet, http.StatusBadRequest
		}

		userTeam, inerr := e.GetUserTeamByUserTeam(userId, *teamId)
		if inerr != nil && inerr != utils.ErrNotFound {
			return inerr, http.StatusInternalServerError
		} else if inerr == utils.ErrNotFound || !userTeam.IsLeader {
			return utils.ErrUnauthorized, http.StatusUnauthorized
		}

		res[*teamId] = struct{}{}
		return nil, http.StatusOK
	}

	err, isInternal := checkLeadership(match.TeamX)
	if err != nil && err != utils.ErrUnauthorized {
		return nil, err, isInternal
	}

	err, isInternal = checkLeadership(match.TeamY)
	if err != nil && err != utils.ErrUnauthorized {
		return nil, err, isInternal
	}

	return res, nil, http.StatusOK
}
