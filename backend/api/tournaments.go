package api

import (
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) PostTournament(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	} else if !me.IsAdmin {
		return &Error{E: utils.ErrUnauthorized}
	}

	var data struct {
		Slug        string
		GameId      string
		Name        string
		Abbr        string
		FoundedAt   time.Time
		Description string
		Email       string
		Twitch      string
		Youtube     string
		Twitter     string
		Facebook    string
		Discord     string
		Web         string
		TwitchLive  string
		Blur        string
		Logo        string
		LogoHasText bool
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	tournament := &models.Tournament{
		TournamentPublic: models.TournamentPublic{
			Slug:        data.Slug,
			GameId:      data.GameId,
			Name:        data.Name,
			Abbr:        data.Abbr,
			FoundedAt:   utils.TimeToDate(data.FoundedAt),
			Description: data.Description,
			Email:       data.Email,
			Twitch:      data.Twitch,
			Youtube:     data.Youtube,
			Twitter:     data.Twitter,
			Facebook:    data.Facebook,
			Discord:     data.Discord,
			Web:         data.Web,
			TwitchLive:  data.TwitchLive,
			Blur:        data.Blur,
			Logo:        data.Logo,
			LogoHasText: data.LogoHasText,
		},
		CreatedBy: me.Id,
	}
	err = e.M.CreateTournament(tournament)
	if err != nil {
		return &Error{E: err}
	}

	return Created(tournament, c, w)
}

func (e *Env) GetTournament(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	tournament, err := e.M.GetTournamentById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.IsAdmin {
			return OK(tournament, c, w)
		}
	}

	return OK(tournament.TournamentPublic, c, w)
}

func (e *Env) GetTournaments(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	var data struct {
		Offset uint64            `param:"offset"`
		Limit  uint64            `param:"count"`
		Filter map[string]string `param:"filter"`
		Sort   string            `param:"sort"`
	}
	err := DecodeQuery(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	tournaments, err := e.M.GetTournaments(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"game_id", "slug"},
		[]string{"id", "game_id", "name", "abbr", "founded_at"},
	))
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.IsAdmin {
			return OK(tournaments, c, w)
		}
	}

	public := make([]models.TournamentPublic, 0)
	for _, tournament := range tournaments {
		public = append(public, tournament.TournamentPublic)
	}

	return OK(public, c, w)
}

func (e *Env) PutTournament(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	} else if !me.IsAdmin {
		return &Error{E: utils.ErrUnauthorized}
	}

	var data struct {
		Slug        *string
		Name        *string
		Abbr        *string
		FoundedAt   *time.Time
		Description *string
		Email       *string
		Twitch      *string
		Youtube     *string
		Twitter     *string
		Facebook    *string
		Discord     *string
		Web         *string
		TwitchLive  *string
		Blur        *string
		Logo        *string
		LogoHasText *bool
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	var tournament *models.Tournament
	err = e.M.Atomic(func(etx *models.Env) error {
		var inerr error
		tournament, inerr = etx.GetTournamentById(c.URLParams["id"])
		if inerr != nil {
			return inerr
		}

		somethingChanged := false

		if data.Slug != nil && *data.Slug != tournament.Slug {
			tournament.Slug = *data.Slug
			somethingChanged = true
		}

		if data.Name != nil && *data.Name != tournament.Name {
			tournament.Name = *data.Name
			somethingChanged = true
		}

		if data.Abbr != nil && *data.Abbr != tournament.Abbr {
			tournament.Abbr = *data.Abbr
			somethingChanged = true
		}

		if data.FoundedAt != nil && *data.FoundedAt != tournament.FoundedAt {
			tournament.FoundedAt = utils.TimeToDate(*data.FoundedAt)
			somethingChanged = true
		}

		if data.Description != nil && *data.Description != tournament.Description {
			inerr = etx.Diff(
				"tournament", tournament.Id, "description", me.Id,
				tournament.Description, *data.Description,
			)
			if inerr != nil {
				return inerr
			}

			tournament.Description = *data.Description
			somethingChanged = true
		}

		if data.Email != nil && *data.Email != tournament.Email {
			tournament.Email = *data.Email
			somethingChanged = true
		}

		if data.Twitch != nil && *data.Twitch != tournament.Twitch {
			tournament.Twitch = *data.Twitch
			somethingChanged = true
		}

		if data.Youtube != nil && *data.Youtube != tournament.Youtube {
			tournament.Youtube = *data.Youtube
			somethingChanged = true
		}

		if data.Twitter != nil && *data.Twitter != tournament.Twitter {
			tournament.Twitter = *data.Twitter
			somethingChanged = true
		}

		if data.Facebook != nil && *data.Facebook != tournament.Facebook {
			tournament.Facebook = *data.Facebook
			somethingChanged = true
		}

		if data.Discord != nil && *data.Discord != tournament.Discord {
			tournament.Discord = *data.Discord
			somethingChanged = true
		}

		if data.Web != nil && *data.Web != tournament.Web {
			tournament.Web = *data.Web
			somethingChanged = true
		}

		if data.TwitchLive != nil && *data.TwitchLive != tournament.TwitchLive {
			tournament.TwitchLive = *data.TwitchLive
			somethingChanged = true
		}

		if data.Blur != nil && *data.Blur != tournament.Blur {
			tournament.Blur = *data.Blur
			somethingChanged = true
		}

		if data.Logo != nil && *data.Logo != tournament.Logo {
			tournament.Logo = *data.Logo
			somethingChanged = true
		}

		if data.LogoHasText != nil && *data.LogoHasText != tournament.LogoHasText {
			tournament.LogoHasText = *data.LogoHasText
			somethingChanged = true
		}

		if !somethingChanged {
			return nil
		}

		return etx.UpdateTournament(tournament, me.Id)
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	return OK(tournament, c, w)
}
