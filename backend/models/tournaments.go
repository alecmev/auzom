package models

import (
	"time"
)

type TournamentPublic struct {
	Id          string    `json:"id"`
	Slug        string    `json:"slug"`
	GameId      string    `db:"game_id" json:"gameId"`
	Name        string    `json:"name"`
	Abbr        string    `json:"abbr"`
	FoundedAt   time.Time `db:"founded_at" json:"foundedAt"`
	Description string    `json:"description"`
	Email       string    `json:"email"`
	Twitch      string    `json:"twitch"`
	Youtube     string    `json:"youtube"`
	Twitter     string    `json:"twitter"`
	Facebook    string    `json:"facebook"`
	Discord     string    `json:"discord"`
	Web         string    `json:"web"`
	TwitchLive  string    `db:"twitch_live" json:"twitchLive"`
	Blur        string    `json:"blur"`
	Logo        string    `json:"logo"`
	LogoHasText bool      `db:"logo_has_text" json:"logoHasText"`
}

type Tournament struct {
	TournamentPublic
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy string     `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
}

func (e *Env) CreateTournament(tournament *Tournament) error {
	return e.Db.Get(
		tournament, `
    INSERT INTO tournament (
      slug,
      game_id,
      name,
      abbr,
      founded_at,
      description,
      email,
      twitch,
      youtube,
      twitter,
      facebook,
      discord,
      web,
      twitch_live,
      created_by,

      blur,
      logo,
      logo_has_text
    )
    VALUES (
      $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,

      $16, $17, $18
    )
    RETURNING *`,
		tournament.Slug,
		tournament.GameId,
		tournament.Name,
		tournament.Abbr,
		tournament.FoundedAt,
		tournament.Description,
		tournament.Email,
		tournament.Twitch,
		tournament.Youtube,
		tournament.Twitter,
		tournament.Facebook,
		tournament.Discord,
		tournament.Web,
		tournament.TwitchLive,
		tournament.CreatedBy,
		tournament.Blur,
		tournament.Logo,
		tournament.LogoHasText,
	)
}

func (e *Env) GetTournamentById(id string) (*Tournament, error) {
	var tournament Tournament
	err := e.Db.Get(
		&tournament, `
    SELECT *
    FROM tournament
    WHERE id=$1`,
		id,
	)
	return &tournament, BetterGetterErrors(err)
}

func (e *Env) GetTournamentBySlug(slug string) (*Tournament, error) {
	var tournament Tournament
	err := e.Db.Get(
		&tournament, `
    SELECT *
    FROM tournament
    WHERE slug=$1`,
		slug,
	)
	return &tournament, BetterGetterErrors(err)
}

func (e *Env) GetTournaments(modifier *QueryModifier) ([]Tournament, error) {
	tournaments := make([]Tournament, 0)
	sql, args, err := modifier.ToSql("tournament", "*")
	if err != nil {
		return tournaments, err
	}

	err = e.Db.Select(&tournaments, sql, args...)
	return tournaments, err
}

func (e *Env) UpdateTournament(tournament *Tournament, updatedBy string) error {
	return e.Db.Get(
		tournament, `
    UPDATE tournament
    SET
      slug=$2,
      name=$3,
      abbr=$4,
      founded_at=$5,
      description=$6,
      email=$7,
      twitch=$8,
      youtube=$9,
      twitter=$10,
      facebook=$11,
      discord=$12,
      web=$13,
      twitch_live=$14,
      blur=$15,
      logo=$16,
      logo_has_text=$17,
      updated_by=$18
    WHERE id=$1
    RETURNING *`,
		tournament.Id,
		tournament.Slug,
		tournament.Name,
		tournament.Abbr,
		tournament.FoundedAt,
		tournament.Description,
		tournament.Email,
		tournament.Twitch,
		tournament.Youtube,
		tournament.Twitter,
		tournament.Facebook,
		tournament.Discord,
		tournament.Web,
		tournament.TwitchLive,
		tournament.Blur,
		tournament.Logo,
		tournament.LogoHasText,
		updatedBy,
	)
}
