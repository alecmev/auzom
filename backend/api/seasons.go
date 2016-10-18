package api

import (
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) PostSeason(
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
		Slug            string
		TournamentId    string
		Name            string
		Abbr            string
		Description     string
		Rules           string
		TeamSize        int `json:",string"`
		TeamSizeMax     int `json:",string"`
		Capacity        int `json:",string"`
		Duration        int `json:",string"`
		YoutubePlaylist string
		Sponsors        string
		PublishedAt     *time.Time
		SignupsOpenedAt *time.Time
		SignupsClosedAt *time.Time
		EndedAt         *time.Time
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	} else if data.TeamSize < 1 {
		return &Error{
			E: err, C: http.StatusBadRequest,
			M: "team size must be 1 or more",
		}
	} else if data.TeamSizeMax < data.TeamSize && data.TeamSizeMax != 0 {
		return &Error{
			E: err, C: http.StatusBadRequest,
			M: "max team size must be 0 for unlimited, or more than team size",
		}
	} else if data.Capacity < 2 && data.Capacity != 0 {
		return &Error{
			E: err, C: http.StatusBadRequest,
			M: "capacity must be 0 for unlimited, or more than 1",
		}
	} else if data.Duration < 1 {
		return &Error{
			E: err, C: http.StatusBadRequest,
			M: "duration must be 1 or more",
		}
	}

	season := &models.Season{
		SeasonPublic: models.SeasonPublic{
			Slug:            data.Slug,
			TournamentId:    data.TournamentId,
			Name:            data.Name,
			Abbr:            data.Abbr,
			Description:     data.Description,
			Rules:           data.Rules,
			TeamSize:        data.TeamSize,
			TeamSizeMax:     data.TeamSizeMax,
			Capacity:        data.Capacity,
			Duration:        data.Duration,
			YoutubePlaylist: data.YoutubePlaylist,
			Sponsors:        data.Sponsors,
			PublishedAt:     data.PublishedAt,
			SignupsOpenedAt: data.SignupsOpenedAt,
			SignupsClosedAt: data.SignupsClosedAt,
			EndedAt:         data.EndedAt,
		},
		CreatedBy: me.Id,
	}
	err = e.M.CreateSeason(season)
	if err != nil {
		return &Error{E: err}
	}

	return Created(season, c, w)
}

func (e *Env) GetSeason(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	season, err := e.M.GetSeasonById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.IsAdmin {
			return OK(season, c, w)
		}
	}

	if season.PublishedAt.After(time.Now()) {
		return &Error{E: utils.ErrNotFound}
	}

	return OK(season.SeasonPublic, c, w)
}

func (e *Env) GetSeasons(
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

	seasons, err := e.M.GetSeasons(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"tournament_id", "slug"},
		[]string{"id", "tournament_id", "name", "abbr", "published_at",
			"signups_opened_at", "signups_closed_at", "ended_at"},
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
			return OK(seasons, c, w)
		}
	}

	now := time.Now()
	public := make([]models.SeasonPublic, 0)
	for _, season := range seasons {
		if season.PublishedAt.After(now) {
			continue
		}

		public = append(public, season.SeasonPublic)
	}

	return OK(public, c, w)
}

func (e *Env) PutSeason(
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
		Slug            *string
		Name            *string
		Abbr            *string
		Description     *string
		Rules           *string
		TeamSize        *int `json:",string"`
		TeamSizeMax     *int `json:",string"`
		Capacity        *int `json:",string"`
		Duration        *int `json:",string"`
		YoutubePlaylist *string
		Sponsors        *string
		PublishedAt     *time.Time
		SignupsOpenedAt *time.Time
		SignupsClosedAt *time.Time
		EndedAt         *time.Time
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	var season *models.Season
	err = e.M.Atomic(func(etx *models.Env) error {
		var inerr error
		season, inerr = etx.GetSeasonById(c.URLParams["id"])
		if inerr != nil {
			return inerr
		}

		somethingChanged := false

		if data.Slug != nil && *data.Slug != season.Slug {
			season.Slug = *data.Slug
			somethingChanged = true
		}

		if data.Name != nil && *data.Name != season.Name {
			season.Name = *data.Name
			somethingChanged = true
		}

		if data.Abbr != nil && *data.Abbr != season.Abbr {
			season.Abbr = *data.Abbr
			somethingChanged = true
		}

		if data.Description != nil && *data.Description != season.Description {
			inerr = etx.Diff(
				"season", season.Id, "description", me.Id,
				season.Description, *data.Description,
			)
			if inerr != nil {
				return inerr
			}

			season.Description = *data.Description
			somethingChanged = true
		}

		if data.Rules != nil && *data.Rules != season.Rules {
			inerr = etx.Diff(
				"season", season.Id, "rules", me.Id,
				season.Rules, *data.Rules,
			)
			if inerr != nil {
				return inerr
			}

			season.Rules = *data.Rules
			somethingChanged = true
		}

		if data.TeamSize != nil && *data.TeamSize != season.TeamSize {
			if *data.TeamSize < 1 {
				return &Error{
					E: err, C: http.StatusBadRequest,
					M: "team size must be 1 or more",
				}
			}

			season.TeamSize = *data.TeamSize
			somethingChanged = true
		}

		if data.TeamSizeMax != nil && *data.TeamSizeMax != season.TeamSizeMax {
			season.TeamSizeMax = *data.TeamSizeMax
			somethingChanged = true
		}

		if season.TeamSizeMax < season.TeamSize && season.TeamSizeMax != 0 {
			return &Error{
				E: err, C: http.StatusBadRequest,
				M: "max team size must be 0 for unlimited, or more than team size",
			}
		}

		if data.Capacity != nil && *data.Capacity != season.Capacity {
			if *data.Capacity < 2 && *data.Capacity != 0 {
				return &Error{
					E: err, C: http.StatusBadRequest,
					M: "capacity must be 0 for unlimited, or more than 1",
				}
			}

			season.Capacity = *data.Capacity
			somethingChanged = true
		}

		if data.Duration != nil && *data.Duration != season.Duration {
			if *data.Duration < 1 {
				return &Error{
					E: err, C: http.StatusBadRequest,
					M: "duration must be 2 or more",
				}
			}

			season.Duration = *data.Duration
			somethingChanged = true
		}

		if data.YoutubePlaylist != nil &&
			*data.YoutubePlaylist != season.YoutubePlaylist {
			// intentionally blank line
			season.YoutubePlaylist = *data.YoutubePlaylist
			somethingChanged = true
		}

		if data.Sponsors != nil && *data.Sponsors != season.Sponsors {
			season.Sponsors = *data.Sponsors
			somethingChanged = true
		}

		if data.PublishedAt != nil &&
			(season.PublishedAt == nil || *data.PublishedAt != *season.PublishedAt) {
			// intentionally blank line
			season.PublishedAt = data.PublishedAt
			somethingChanged = true
		}

		if data.SignupsOpenedAt != nil &&
			(season.SignupsOpenedAt == nil ||
				*data.SignupsOpenedAt != *season.SignupsOpenedAt) {
			// intentionally blank line
			season.SignupsOpenedAt = data.SignupsOpenedAt
			somethingChanged = true
		}

		if data.SignupsClosedAt != nil &&
			(season.SignupsClosedAt == nil ||
				*data.SignupsClosedAt != *season.SignupsClosedAt) {
			// intentionally blank line
			season.SignupsClosedAt = data.SignupsClosedAt
			somethingChanged = true
		}

		if data.EndedAt != nil && season.EndedAt == nil {
			return &Error{
				C: http.StatusBadRequest, M: "use end patch action to end this season",
			}
		} else if data.EndedAt != nil &&
			(season.EndedAt == nil ||
				*data.EndedAt != *season.EndedAt) {
			// intentionally blank line
			season.EndedAt = data.EndedAt
			somethingChanged = true
		}

		if !somethingChanged {
			return nil
		}

		return etx.UpdateSeason(season, me.Id)
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	return OK(season, c, w)
}

func (e *Env) PatchSeason(
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
		Action string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	var season *models.Season
	err = e.M.Atomic(func(etx *models.Env) error {
		var inerr error
		season, inerr = etx.GetSeasonById(c.URLParams["id"])
		if inerr != nil {
			return inerr
		} else if season.EndedAt != nil {
			return &Error{
				C: http.StatusBadRequest, M: "this season has already ended",
			}
		}

		now := time.Now()
		if data.Action == "accept" {
			if season.SignupsClosedAt == nil || season.SignupsClosedAt.After(now) {
				return &Error{
					C: http.StatusBadRequest,
					M: "can't accept an application while the signups are open",
				}
			}

			teamSeasonRequests, inerr := etx.GetTeamSeasonRequestsBySeason(season.Id)
			if inerr != nil {
				return inerr
			}

			decision := true
			for _, teamSeasonRequest := range teamSeasonRequests {
				teamSeasonRequest.Decision = &decision
				teamSeasonRequest.DecidedAt = &now
				teamSeasonRequest.DecidedBy = &me.Id
				inerr = PatchTeamSeasonRequestHelper(etx, &teamSeasonRequest)
				if inerr != nil {
					return inerr
				}
			}
		} else if data.Action == "end" {
			season.EndedAt = &now
			inerr = etx.UpdateSeason(season, me.Id)
			if inerr != nil {
				return inerr
			}

			teamSeasons, inerr := etx.GetTeamSeasonsBySeason(season)
			if inerr != nil {
				return inerr
			}

			for _, teamSeason := range teamSeasons {
				if teamSeason.LeftAt != nil {
					continue
				}

				teamSeason.LeftAt = &now
				teamSeason.IsDone = true
				inerr = etx.UpdateTeamSeason(&teamSeason, me.Id)
				if inerr != nil {
					return inerr
				}
			}

			teamSeasonRequests, inerr := etx.GetTeamSeasonRequestsBySeason(season.Id)
			if inerr != nil {
				return inerr
			}

			decision := false
			for _, teamSeasonRequest := range teamSeasonRequests {
				teamSeasonRequest.Decision = &decision
				teamSeasonRequest.DecidedAt = &now
				teamSeasonRequest.DecidedBy = &me.Id
				// patch helper isn't needed for negative decisions
				inerr = etx.UpdateTeamSeasonRequest(&teamSeasonRequest)
				if inerr != nil {
					return inerr
				}
			}

			// TODO: admin-approve all team applications? not sure if it's possible
		} else {
			return &Error{C: http.StatusBadRequest, M: "bad action"}
		}

		return nil
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	return OK(season, c, w)
}
