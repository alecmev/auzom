package models

import (
	"time"
)

type MatchRoundPublic struct {
	Id            string `json:"id"`
	MatchReportId string `db:"match_report_id" json:"matchReportId"`
	GameMapId     string `db:"game_map_id" json:"gameMapId"`

	IsTeamXOnSideY bool    `db:"is_team_x_on_side_y" json:"isTeamXOnSideY"`
	IsNotPlayed    bool    `db:"is_not_played" json:"isNotPlayed"`
	RawScoreX      float64 `db:"raw_score_x" json:"rawScoreX"`
	RawScoreY      float64 `db:"raw_score_y" json:"rawScoreY"`

	OverrideReason    string   `db:"override_reason" json:"overrideReason"`
	IsPenalOverride   bool     `db:"is_penal_override" json:"isPenalOverride"`
	RawScoreXOverride *float64 `db:"raw_score_x_override" json:"rawScoreXOverride"`
	RawScoreYOverride *float64 `db:"raw_score_y_override" json:"rawScoreYOverride"`
}

type MatchRound struct {
	MatchRoundPublic
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	CreatedBy string    `db:"created_by" json:"createdBy"`
}

func (e *Env) ClaimMatchRoundId() (id string, err error) {
	err = e.Db.Get(&id, `SELECT nextval('match_round_id_seq')`)
	return
}

func (e *Env) CreateMatchRound(matchRound *MatchRound) error {
	if matchRound.Id == "" {
		id, err := e.ClaimMatchRoundId()
		if err != nil {
			return err
		}

		matchRound.Id = id
	}

	return e.Db.Get(
		matchRound, `
    INSERT INTO match_round (
      id, match_report_id, game_map_id,

      is_team_x_on_side_y, is_not_played, raw_score_x, raw_score_y,

      override_reason, is_penal_override,
      raw_score_x_override, raw_score_y_override,

      created_by
    )
    VALUES (
      $1, $2, $3,
      $4, $5, $6, $7,
      $8, $9, $10, $11,
      $12
    )
    RETURNING *`,
		matchRound.Id,
		matchRound.MatchReportId,
		matchRound.GameMapId,

		matchRound.IsTeamXOnSideY,
		matchRound.IsNotPlayed,
		matchRound.RawScoreX,
		matchRound.RawScoreY,

		matchRound.OverrideReason,
		matchRound.IsPenalOverride,
		matchRound.RawScoreXOverride,
		matchRound.RawScoreYOverride,

		matchRound.CreatedBy,
	)
}

func (e *Env) GetMatchRoundById(id string) (*MatchRound, error) {
	var matchRound MatchRound
	err := e.Db.Get(
		&matchRound, `
    SELECT *
    FROM match_round
    WHERE id=$1`,
		id,
	)
	return &matchRound, BetterGetterErrors(err)
}

func (e *Env) GetMatchRoundsByReportId(reportId string) ([]MatchRound, error) {
	matchRounds := make([]MatchRound, 0)
	err := e.Db.Select(
		&matchRounds, `
    SELECT *
    FROM match_round
    WHERE match_report_id=$1`,
		reportId,
	)
	return matchRounds, err
}

func (e *Env) GetMatchRounds(modifier *QueryModifier) ([]MatchRound, error) {
	matchRounds := make([]MatchRound, 0)
	sql, args, err := modifier.ToSql("match_round", "*")
	if err != nil {
		return matchRounds, err
	}

	err = e.Db.Select(&matchRounds, sql, args...)
	return matchRounds, err
}
