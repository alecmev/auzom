package models

import (
	"time"
)

type MatchPenaltyPublic struct {
	Id            string  `json:"id"`
	MatchReportId string  `db:"match_report_id" json:"matchReportId"`
	MatchRoundId  *string `db:"match_round_id" json:"matchRoundId"`

	Reason    string  `json:"reason"`
	ScoreX    float64 `db:"score_x" json:"scoreX"`
	ScoreY    float64 `db:"score_y" json:"scoreY"`
	RawScoreX float64 `db:"raw_score_x" json:"rawScoreX"`
	RawScoreY float64 `db:"raw_score_y" json:"rawScoreY"`
}

type MatchPenalty struct {
	MatchPenaltyPublic
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	CreatedBy string    `db:"created_by" json:"createdBy"`
}

func (e *Env) CreateMatchPenalty(matchPenalty *MatchPenalty) error {
	return e.Db.Get(
		matchPenalty, `
    INSERT INTO match_penalty (
      match_report_id, match_round_id,
      reason, score_x, score_y, raw_score_x, raw_score_y,
      created_by
    )
    VALUES (
      $1, $2,
      $3, $4, $5, $6, $7,
      $8
    )
    RETURNING *`,
		matchPenalty.MatchReportId,
		matchPenalty.MatchRoundId,

		matchPenalty.Reason,
		matchPenalty.ScoreX,
		matchPenalty.ScoreY,
		matchPenalty.RawScoreX,
		matchPenalty.RawScoreY,

		matchPenalty.CreatedBy,
	)
}

func (e *Env) GetMatchPenaltyById(id string) (*MatchPenalty, error) {
	var matchPenalty MatchPenalty
	err := e.Db.Get(
		&matchPenalty, `
    SELECT *
    FROM match_penalty
    WHERE id=$1`,
		id,
	)
	return &matchPenalty, BetterGetterErrors(err)
}

func (e *Env) GetMatchPenalties(
	modifier *QueryModifier,
) ([]MatchPenalty, error) {
	matchPenalties := make([]MatchPenalty, 0)
	sql, args, err := modifier.ToSql("match_penalty", "*")
	if err != nil {
		return matchPenalties, err
	}

	err = e.Db.Select(&matchPenalties, sql, args...)
	return matchPenalties, err
}
