package api

import (
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) PostStage(
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
		SeasonId  string
		Slug      string
		Name      string
		Abbr      string
		StartedAt *time.Time
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	stage := &models.Stage{
		StagePublic: models.StagePublic{
			SeasonId:  data.SeasonId,
			Slug:      data.Slug,
			Name:      data.Name,
			Abbr:      data.Abbr,
			StartedAt: data.StartedAt,
		},
		CreatedBy: me.Id,
	}
	err = e.M.CreateStage(stage)
	if err != nil {
		return &Error{E: err}
	}

	return Created(stage, c, w)
}

func (e *Env) GetStage(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	stage, err := e.M.GetStageById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.IsAdmin {
			return OK(stage, c, w)
		}
	}

	return OK(stage.StagePublic, c, w)
}

func (e *Env) GetStages(
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

	stages, err := e.M.GetStages(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"season_id", "slug"},
		[]string{"id", "season_id", "name", "abbr", "started_at"},
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
			return OK(stages, c, w)
		}
	}

	public := make([]models.StagePublic, 0, len(stages))
	for _, stage := range stages {
		public = append(public, stage.StagePublic)
	}

	return OK(public, c, w)
}

func (e *Env) PutStage(
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

	stage, err := e.M.GetStageById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	var data struct {
		Slug      *string
		Name      *string
		Abbr      *string
		StartedAt *time.Time
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	somethingChanged := false

	if data.Slug != nil && *data.Slug != stage.Slug {
		stage.Slug = *data.Slug
		somethingChanged = true
	}

	if data.Name != nil && *data.Name != stage.Name {
		stage.Name = *data.Name
		somethingChanged = true
	}

	if data.Abbr != nil && *data.Abbr != stage.Abbr {
		stage.Abbr = *data.Abbr
		somethingChanged = true
	}

	if data.StartedAt != nil &&
		(stage.StartedAt == nil || *data.StartedAt != *stage.StartedAt) {
		// intentionally blank line
		stage.StartedAt = data.StartedAt
		somethingChanged = true
	}

	if !somethingChanged {
		return OK(stage, c, w)
	}

	err = e.M.UpdateStage(stage, me.Id)
	if err != nil {
		return &Error{E: err}
	}

	return OK(stage, c, w)
}
