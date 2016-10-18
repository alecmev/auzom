package models

import (
	"time"
)

type StagePublic struct {
	Id        string     `json:"id"`
	SeasonId  string     `db:"season_id" json:"seasonId"`
	Slug      string     `json:"slug"`
	Name      string     `json:"name"`
	Abbr      string     `json:"abbr"`
	StartedAt *time.Time `db:"started_at" json:"startedAt"`
}

type Stage struct {
	StagePublic
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy string     `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateStage(stage *Stage) error {
	return e.Db.Get(
		stage, `
    INSERT INTO stage (
      season_id,
      slug,
      name,
      abbr,
      started_at,
      created_by
    )
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *`,
		stage.SeasonId,
		stage.Slug,
		stage.Name,
		stage.Abbr,
		stage.StartedAt,
		stage.CreatedBy,
	)
}

func (e *Env) GetStageById(id string) (*Stage, error) {
	var stage Stage
	err := e.Db.Get(
		&stage, `
    SELECT *
    FROM stage
    WHERE id=$1`,
		id,
	)
	return &stage, BetterGetterErrors(err)
}

func (e *Env) GetStages(modifier *QueryModifier) ([]Stage, error) {
	stages := make([]Stage, 0)
	sql, args, err := modifier.ToSql("stage", "*")
	if err != nil {
		return stages, err
	}

	err = e.Db.Select(&stages, sql, args...)
	return stages, err
}

func (e *Env) UpdateStage(stage *Stage, updatedBy string) error {
	err := e.Db.Get(
		stage, `
    UPDATE stage
    SET
      slug=$2,
      name=$3,
      abbr=$4,
      started_at=$5,
      updated_by=$6
    WHERE id=$1
    RETURNING *`,
		stage.Id,
		stage.Slug,
		stage.Name,
		stage.Abbr,
		stage.StartedAt,
		updatedBy,
	)
	return err
}
