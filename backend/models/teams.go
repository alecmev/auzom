package models

import (
	"time"
)

type TeamPublic struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Abbr        string     `json:"abbr"`
	Logo        string     `json:"logo"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy   string     `db:"created_by" json:"createdBy"`
	DisbandedAt *time.Time `db:"disbanded_at" json:"disbandedAt"`
	DisbandedBy *string    `db:"disbanded_by" json:"disbandedBy"`
}

type Team struct {
	TeamPublic
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateTeam(team *Team) error {
	return e.Db.Get(
		team, `
    INSERT INTO team (name, abbr, logo, created_by)
    VALUES ($1, $2, $3, $4)
    RETURNING *`,
		team.Name,
		team.Abbr,
		team.Logo,
		team.CreatedBy,
	)
}

func (e *Env) GetTeamById(id string) (*Team, error) {
	var team Team
	err := e.Db.Get(
		&team, `
    SELECT *
    FROM team
    WHERE id=$1`,
		id,
	)
	return &team, BetterGetterErrors(err)
}

func (e *Env) GetTeams(modifier *QueryModifier) ([]Team, error) {
	teams := make([]Team, 0)
	sql, args, err := modifier.ToSql("team", "*")
	if err != nil {
		return teams, err
	}

	err = e.Db.Select(&teams, sql, args...)
	return teams, err
}

func (e *Env) UpdateTeam(team *Team, updatedBy string) error {
	err := e.Db.Get(
		team, `
    UPDATE team
    SET
      name=$2,
      abbr=$3,
      logo=$4,
      disbanded_at=$5,
      disbanded_by=$6,
      updated_by=$7
    WHERE id=$1
    RETURNING *`,
		team.Id,
		team.Name,
		team.Abbr,
		team.Logo,
		team.DisbandedAt,
		team.DisbandedBy,
		updatedBy,
	)
	return err
}
