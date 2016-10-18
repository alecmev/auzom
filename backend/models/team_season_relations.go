package models

import (
	"github.com/jmoiron/sqlx"
)

type TeamSeasonRelation struct {
	TeamId    string `db:"team_id"`
	SeasonId  string `db:"season_id"`
	IsRequest bool   `db:"is_request"`
}

func (e *Env) GetTeamSeasonRelationsBySeason(season *Season) (
	[]TeamSeasonRelation, error,
) {
	relations := make([]TeamSeasonRelation, 0)
	err := e.Db.Select(
		&relations, `
    SELECT team_id, season_id, FALSE AS is_request
    FROM team_season
    WHERE season_id=$1 AND left_at IS NULL
    UNION ALL
    SELECT team_id, season_id, TRUE AS is_request
    FROM team_season_request
    WHERE season_id=$1 AND decision IS NULL`,
		season.Id,
	)
	return relations, err
}

func (e *Env) GetTeamSeasonRelationsByTeam(team *Team) (
	[]TeamSeasonRelation, error,
) {
	_, ok := e.Db.(*sqlx.Tx)
	if ok {
		return getTeamSeasonRelationsByTeam(e.Db, team.Id)
	}

	var relations []TeamSeasonRelation
	err := e.Atomic(func(etx *Env) error {
		var inerr error
		relations, inerr = getTeamSeasonRelationsByTeam(etx.Db, team.Id)
		return inerr
	})
	return relations, err
}

func getTeamSeasonRelationsByTeam(db database, teamId string) (
	[]TeamSeasonRelation, error,
) {
	relations := make([]TeamSeasonRelation, 0)
	_, err := db.Exec(`
    CREATE TEMP TABLE tmp ON COMMIT DROP AS
    SELECT season_id
    FROM team_season
    WHERE team_id=$1 AND left_at IS NULL
    UNION ALL
    SELECT season_id
    FROM team_season_request
    WHERE team_id=$1 AND decision IS NULL`,
		teamId,
	)
	if err != nil {
		return nil, err
	}

	err = db.Select(
		&relations, `
    SELECT ts.team_id, ts.season_id, FALSE AS is_request
    FROM team_season ts
    INNER JOIN tmp
    ON ts.season_id=tmp.season_id
    WHERE ts.left_at IS NULL AND ts.team_id!=$1
    UNION ALL
    SELECT tsr.team_id, tsr.season_id, TRUE AS is_request
    FROM team_season_request tsr
    INNER JOIN tmp
    ON tsr.season_id=tmp.season_id
    WHERE tsr.decision IS NULL AND tsr.team_id!=$1`,
		teamId,
	)
	return relations, err
}
