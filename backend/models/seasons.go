package models

import (
	"time"
)

type SeasonPublic struct {
	Id              string     `json:"id"`
	Slug            string     `json:"slug"`
	TournamentId    string     `db:"tournament_id" json:"tournamentId"`
	Name            string     `json:"name"`
	Abbr            string     `json:"abbr"`
	PublishedAt     *time.Time `db:"published_at" json:"publishedAt"`
	Description     string     `json:"description"`
	Rules           string     `json:"rules"`
	TeamSize        int        `db:"team_size" json:"teamSize"`
	TeamSizeMax     int        `db:"team_size_max" json:"teamSizeMax"`
	Capacity        int        `json:"capacity"`
	Duration        int        `json:"duration"`
	YoutubePlaylist string     `db:"youtube_playlist" json:"youtubePlaylist"`
	Sponsors        string     `json:"sponsors"`
	SignupsOpenedAt *time.Time `db:"signups_opened_at" json:"signupsOpenedAt"`
	SignupsClosedAt *time.Time `db:"signups_closed_at" json:"signupsClosedAt"`
	EndedAt         *time.Time `db:"ended_at" json:"endedAt"`
}

type Season struct {
	SeasonPublic
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy string     `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateSeason(season *Season) error {
	return e.Db.Get(
		season, `
    INSERT INTO season (
      slug,
      tournament_id,
      name,
      abbr,
      published_at,
      description,
      rules,
      team_size,
      team_size_max,
      capacity,
      duration,
      youtube_playlist,
      sponsors,
      created_by,
      signups_opened_at,
      signups_closed_at,
      ended_at
    )
    VALUES (
      $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
    )
    RETURNING *`,
		season.Slug,
		season.TournamentId,
		season.Name,
		season.Abbr,
		season.PublishedAt,
		season.Description,
		season.Rules,
		season.TeamSize,
		season.TeamSizeMax,
		season.Capacity,
		season.Duration,
		season.YoutubePlaylist,
		season.Sponsors,
		season.CreatedBy,
		season.SignupsOpenedAt,
		season.SignupsClosedAt,
		season.EndedAt,
	)
}

func (e *Env) GetSeasonById(id string) (*Season, error) {
	var season Season
	err := e.Db.Get(
		&season, `
    SELECT *
    FROM season
    WHERE id=$1`,
		id,
	)
	return &season, BetterGetterErrors(err)
}

func (e *Env) GetSeasonByBracket(bracket *Bracket) (*Season, error) {
	var season Season
	err := e.Db.Get(
		&season, `
    SELECT season.*
    FROM bracket 
    JOIN stage ON stage.id=bracket.stage_id
    JOIN season ON season.id=stage.season_id
    WHERE bracket.id=$1`,
		bracket.Id,
	)
	return &season, BetterGetterErrors(err)
}

func (e *Env) GetSeasons(modifier *QueryModifier) ([]Season, error) {
	seasons := make([]Season, 0)
	sql, args, err := modifier.ToSql("season", "*")
	if err != nil {
		return seasons, err
	}

	err = e.Db.Select(&seasons, sql, args...)
	return seasons, err
}

func (e *Env) UpdateSeason(season *Season, updatedBy string) error {
	return e.Db.Get(
		season, `
    UPDATE season
    SET
      slug=$2,
      name=$3,
      abbr=$4,
      published_at=$5,
      description=$6,
      rules=$7,
      team_size=$8,
      team_size_max=$9,
      capacity=$10,
      duration=$11,
      youtube_playlist=$12,
      sponsors=$13,
      signups_opened_at=$14,
      signups_closed_at=$15,
      ended_at=$16,
      updated_by=$17
    WHERE id=$1
    RETURNING *`,
		season.Id,
		season.Slug,
		season.Name,
		season.Abbr,
		season.PublishedAt,
		season.Description,
		season.Rules,
		season.TeamSize,
		season.TeamSizeMax,
		season.Capacity,
		season.Duration,
		season.YoutubePlaylist,
		season.Sponsors,
		season.SignupsOpenedAt,
		season.SignupsClosedAt,
		season.EndedAt,
		updatedBy,
	)
}
