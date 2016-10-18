package models

import (
	"time"
)

type MatchReportPublic struct {
	Id      string `json:"id"`
	MatchId string `db:"match_id" json:"matchId"`

	ScoreX    float64 `db:"score_x" json:"scoreX"`
	ScoreY    float64 `db:"score_y" json:"scoreY"`
	RawScoreX float64 `db:"raw_score_x" json:"rawScoreX"`
	RawScoreY float64 `db:"raw_score_y" json:"rawScoreY"`

	OverrideReason    string   `db:"override_reason" json:"overrideReason"`
	IsPenalOverride   bool     `db:"is_penal_override" json:"isPenalOverride"`
	ScoreXOverride    *float64 `db:"score_x_override" json:"scoreXOverride"`
	ScoreYOverride    *float64 `db:"score_y_override" json:"scoreYOverride"`
	RawScoreXOverride *float64 `db:"raw_score_x_override" json:"rawScoreXOverride"`
	RawScoreYOverride *float64 `db:"raw_score_y_override" json:"rawScoreYOverride"`

	MapsPlayed   int `db:"maps_played" json:"mapsPlayed"`
	MapsX        int `db:"maps_x" json:"mapsX"`
	MapsY        int `db:"maps_y" json:"mapsY"`
	RoundsPlayed int `db:"rounds_played" json:"roundsPlayed"`
	RoundsX      int `db:"rounds_x" json:"roundsX"`
	RoundsY      int `db:"rounds_y" json:"roundsY"`
}

type MatchReport struct {
	MatchReportPublic
	TeamBy       *string    `db:"team_by" json:"teamBy"`
	AgreedUponAt *time.Time `db:"agreed_upon_at" json:"agreedUponAt"`
	AgreedUponBy *string    `db:"agreed_upon_by" json:"agreedUponBy"`
	CreatedAt    time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy    string     `db:"created_by" json:"createdBy"`
}

func (e *Env) ClaimMatchReportId() (id string, err error) {
	err = e.Db.Get(&id, `SELECT nextval('match_report_id_seq')`)
	return
}

func (e *Env) CreateMatchReport(matchReport *MatchReport) error {
	if matchReport.Id == "" {
		id, err := e.ClaimMatchReportId()
		if err != nil {
			return err
		}

		matchReport.Id = id
	}

	return e.Db.Get(
		matchReport, `
    INSERT INTO match_report (
      id, match_id,

      score_x, score_y, raw_score_x, raw_score_y,

      override_reason, is_penal_override,
      score_x_override, score_y_override,
      raw_score_x_override, raw_score_y_override,

      maps_played, maps_x, maps_y,
      rounds_played, rounds_x, rounds_y,

      team_by, created_by
    )
    VALUES (
      $1, $2,
      $3, $4, $5, $6,
      $7, $8, $9, $10, $11, $12,
      $13, $14, $15, $16, $17, $18,
      $19, $20
    )
    RETURNING *`,
		matchReport.Id,
		matchReport.MatchId,

		matchReport.ScoreX,
		matchReport.ScoreY,
		matchReport.RawScoreX,
		matchReport.RawScoreY,

		matchReport.OverrideReason,
		matchReport.IsPenalOverride,
		matchReport.ScoreXOverride,
		matchReport.ScoreYOverride,
		matchReport.RawScoreXOverride,
		matchReport.RawScoreYOverride,

		matchReport.MapsPlayed,
		matchReport.MapsX,
		matchReport.MapsY,
		matchReport.RoundsPlayed,
		matchReport.RoundsX,
		matchReport.RoundsY,

		matchReport.TeamBy,
		matchReport.CreatedBy,
	)
}

func (e *Env) GetMatchReportById(id string) (*MatchReport, error) {
	var matchReport MatchReport
	err := e.Db.Get(
		&matchReport, `
    SELECT *
    FROM match_report
    WHERE id=$1`,
		id,
	)
	return &matchReport, BetterGetterErrors(err)
}

func (e *Env) GetMatchReports(modifier *QueryModifier) ([]MatchReport, error) {
	matchReports := make([]MatchReport, 0)
	sql, args, err := modifier.ToSql("match_report", "*")
	if err != nil {
		return matchReports, err
	}

	err = e.Db.Select(&matchReports, sql, args...)
	return matchReports, err
}

func (e *Env) UpdateMatchReport(matchReport *MatchReport) error {
	err := e.Db.Get(
		matchReport, `
    UPDATE match_report
    SET agreed_upon_at=$2, agreed_upon_by=$3
    WHERE id=$1
    RETURNING *`,
		matchReport.Id,
		matchReport.AgreedUponAt,
		matchReport.AgreedUponBy,
	)
	return err
}
