package api

import (
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) PostGame(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}
	if !me.IsAdmin {
		return &Error{E: utils.ErrUnauthorized}
	}

	var data struct {
		Slug               string
		Name               string
		Abbr               string
		ReleasedAt         *time.Time
		Cover              string
		Summary            string
		VerificationHandle string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	if data.ReleasedAt != nil {
		*data.ReleasedAt = utils.TimeToDate(*data.ReleasedAt)
	}

	vh := &data.VerificationHandle
	if *vh == "" {
		vh = nil
	}

	game := &models.Game{
		GamePublic: models.GamePublic{
			Slug:       data.Slug,
			Name:       data.Name,
			Abbr:       data.Abbr,
			ReleasedAt: data.ReleasedAt,
			Cover:      data.Cover,
			Summary:    data.Summary,
		},
		CreatedBy:          me.Id,
		VerificationHandle: vh,
	}
	err = e.M.CreateGame(game)
	if err != nil {
		return &Error{E: err}
	}

	return Created(game, c, w)
}

func (e *Env) GetGame(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	game, err := e.M.GetGameById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.IsAdmin {
			return OK(game, c, w)
		}
	}

	return OK(game.GamePublic, c, w)
}

func (e *Env) GetGames(
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

	games, err := e.M.GetGames(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"franchise_id", "slug"},
		[]string{"id", "franchise_id", "name", "abbr", "released_at"},
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
			return OK(games, c, w)
		}
	}

	public := make([]models.GamePublic, 0)
	for _, game := range games {
		public = append(public, game.GamePublic)
	}

	return OK(public, c, w)
}

func (e *Env) PutGame(
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
		Slug       *string
		Name       *string
		Abbr       *string
		ReleasedAt *time.Time
		Cover      *string
		Summary    *string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	var game *models.Game
	err = e.M.Atomic(func(etx *models.Env) error {
		var inerr error
		game, inerr = etx.GetGameById(c.URLParams["id"])
		if inerr != nil {
			return inerr
		}

		somethingChanged := false

		if data.Slug != nil && *data.Slug != game.Slug {
			game.Slug = *data.Slug
			somethingChanged = true
		}

		if data.Name != nil && *data.Name != game.Name {
			game.Name = *data.Name
			somethingChanged = true
		}

		if data.Abbr != nil && *data.Abbr != game.Abbr {
			game.Abbr = *data.Abbr
			somethingChanged = true
		}

		if data.ReleasedAt != nil &&
			(game.ReleasedAt == nil || *data.ReleasedAt != *game.ReleasedAt) {
			// intentionally blank line
			*game.ReleasedAt = utils.TimeToDate(*data.ReleasedAt)
			somethingChanged = true
		}

		if data.Cover != nil && *data.Cover != game.Cover {
			game.Cover = *data.Cover
			somethingChanged = true
		}

		if data.Summary != nil && *data.Summary != game.Summary {
			inerr = etx.Diff(
				"game", game.Id, "summary", me.Id,
				game.Summary, *data.Summary,
			)
			if inerr != nil {
				return inerr
			}

			game.Summary = *data.Summary
			somethingChanged = true
		}

		if !somethingChanged {
			return nil
		}

		return etx.UpdateGame(game, me.Id)
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	return OK(game, c, w)
}
