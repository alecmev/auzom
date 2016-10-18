package api

import (
	"net/http"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) GetBracketRound(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	bracketRound, err := e.M.GetBracketRoundById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		}
		if me.IsAdmin {
			return OK(bracketRound, c, w)
		}
	}

	return OK(bracketRound.BracketRoundPublic, c, w)
}

func (e *Env) GetBracketRounds(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	var data models.QueryBase
	err := DecodeQuery(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	bracketRounds, err := e.M.GetBracketRounds(models.NewQueryModifier(
		data,
		[]string{"bracket_id", "number"},
		[]string{"id", "bracket_id", "number"},
	))
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		}
		if me.IsAdmin {
			return OK(bracketRounds, c, w)
		}
	}

	public := make([]models.BracketRoundPublic, 0)
	for _, bracketRound := range bracketRounds {
		public = append(public, bracketRound.BracketRoundPublic)
	}

	return OK(public, c, w)
}

func (e *Env) PutBracketRound(
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

	bracketRound, err := e.M.GetBracketRoundById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	var data struct {
		Name             *string
		Description      *string
		MapVetoProcedure *string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	somethingChanged := false

	if data.Name != nil && *data.Name != bracketRound.Name {
		bracketRound.Name = *data.Name
		somethingChanged = true
	}

	if data.Description != nil && *data.Description != bracketRound.Description {
		bracketRound.Description = *data.Description
		somethingChanged = true
	}

	// TODO: disallow changing veto procedure if one is in progress (deep issue)
	if data.MapVetoProcedure != nil {
		procedure, apierr := parseMapVetoProcedure(*data.MapVetoProcedure)
		if apierr != nil {
			return apierr
		} else if len(procedure) == 0 {
			*data.MapVetoProcedure = ""
		}

		if *data.MapVetoProcedure != bracketRound.MapVetoProcedure {
			bracketRound.MapVetoProcedure = *data.MapVetoProcedure
			somethingChanged = true
		}
	}

	if !somethingChanged {
		return OK(bracketRound, c, w)
	}

	err = e.M.UpdateBracketRound(bracketRound, me.Id)
	if err != nil {
		return &Error{E: err}
	}

	return OK(bracketRound, c, w)
}
