package api

import (
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) PostAttentionRequest(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	} else if me.IsAdmin {
		return &Error{C: http.StatusBadRequest, M: "you're an admin yourself..."}
	}

	var data struct {
		Target   string
		TargetId string
		Message  string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	} else if data.Target != "match" {
		return &Error{C: http.StatusBadRequest, M: "bad target"}
	}

	match, err := e.M.GetMatchById(data.TargetId)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest, M: "bad match id"}
	}

	var teamAmLeaderOf *string
	teamsAmLeaderOf, err, status := match.UserIsLeaderOf(e.M, me.Id)
	if err != nil {
		return &Error{E: err, C: status}
	} else if len(teamsAmLeaderOf) == 0 {
		return &Error{E: utils.ErrUnauthorized}
	} else if len(teamsAmLeaderOf) > 1 {
		return &Error{
			C: http.StatusBadRequest,
			M: "you're a leader in both teams, ask somebody else to do this",
		}
	} else {
		for x := range teamsAmLeaderOf { // no better way...
			teamAmLeaderOf = &x
		}
	}

	request := &models.AttentionRequest{
		AttentionRequestPublic: models.AttentionRequestPublic{
			Target:    data.Target,
			TargetId:  data.TargetId,
			Message:   data.Message,
			CreatedBy: me.Id,
			TeamBy:    teamAmLeaderOf,
		},
	}
	err = e.M.CreateAttentionRequest(request)
	if err != nil {
		return &Error{E: err}
	}

	e.Slack.Send("attention-requests", me,
		"*Target:* /matches/"+request.TargetId+"\n"+
			"*Message:* "+request.Message,
	)
	return Created(request.AttentionRequestPublic, c, w)
}

func (e *Env) getAttentionRequestStuff(c web.C) (
	request *models.AttentionRequest, me *models.User, err *Error,
) {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		err = &Error{E: utils.ErrUnauthorized}
		return
	}

	var inerr error
	request, inerr = e.M.GetAttentionRequestById(c.URLParams["id"])
	if inerr != nil {
		err = &Error{E: inerr}
		return
	}

	me, inerr = e.M.GetUserById(session.UserId)
	if inerr != nil {
		err = &Error{E: inerr, C: http.StatusInternalServerError}
		return
	}

	if me.IsAdmin {
		return
	}

	if request.Target == "match" {
		match, inerr := e.M.GetMatchById(request.TargetId)
		if inerr != nil {
			err = &Error{E: inerr, C: http.StatusInternalServerError}
			return
		}

		teamsAmLeaderOf, inerr, status := match.UserIsLeaderOf(e.M, me.Id)
		if inerr != nil {
			err = &Error{E: inerr, C: status}
			return
		} else if _, ok := teamsAmLeaderOf[*request.TeamBy]; !ok {
			err = &Error{E: utils.ErrNotFound}
			return
		}
	}

	err = &Error{
		C: http.StatusInternalServerError, M: "invalid request target",
	}
	return
}

func (e *Env) GetAttentionRequest(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	request, me, outerr := e.getAttentionRequestStuff(c)
	if outerr != nil {
		return outerr
	}

	if me.IsAdmin {
		return OK(request, c, w)
	}

	return OK(request.AttentionRequestPublic, c, w)
}

func (e *Env) GetAttentionRequests(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

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

	requests, err := e.M.GetAttentionRequests(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{
			"target", "target_id", "created_by", "team_by", "claimed_by",
			"resolved_at",
		},
		[]string{}, // TODO
	))
	if err != nil {
		return &Error{E: err}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	} else if me.IsAdmin {
		return OK(requests, c, w)
	}

	public := make([]models.AttentionRequestPublic, 0)
	for _, request := range requests {
		if request.Target == "match" {
			match, err := e.M.GetMatchById(request.TargetId)
			if err != nil {
				return &Error{E: err, C: http.StatusInternalServerError}
			}

			teamsAmLeaderOf, err, status := match.UserIsLeaderOf(e.M, me.Id)
			if err != nil {
				return &Error{E: err, C: status}
			} else if _, ok := teamsAmLeaderOf[*request.TeamBy]; !ok {
				continue
			}
		} else {
			return &Error{
				C: http.StatusInternalServerError, M: "invalid request target",
			}
		}

		public = append(public, request.AttentionRequestPublic)
	}

	return OK(public, c, w)
}

func (e *Env) PutAttentionRequest(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	request, me, outerr := e.getAttentionRequestStuff(c)
	if outerr != nil {
		return outerr
	} else if request.ResolvedAt != nil {
		return &Error{
			C: http.StatusBadRequest, M: "can't modify a resolved request",
		}
	}

	var data struct {
		Message *string
	}
	err := Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	somethingChanged := false

	if data.Message != nil && *data.Message != request.Message {
		request.Message = *data.Message
		somethingChanged = true
	}

	if !somethingChanged {
		return OK(request, c, w)
	}

	err = e.M.UpdateAttentionRequest(request, me.Id)
	if err != nil {
		return &Error{E: err}
	}

	if me.IsAdmin {
		return OK(request, c, w)
	}

	return OK(request.AttentionRequestPublic, c, w)
}

func (e *Env) PatchAttentionRequest(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	request, me, outerr := e.getAttentionRequestStuff(c)
	if outerr != nil {
		return outerr
	} else if !me.IsAdmin {
		return &Error{E: utils.ErrUnauthorized}
	} else if request.ResolvedAt != nil {
		return &Error{
			C: http.StatusBadRequest, M: "can't modify a resolved request",
		}
	}

	var data struct {
		Action string
	}
	err := Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	now := time.Now()
	if data.Action == "claim" {
		request.ClaimedAt = &now
		request.ClaimedBy = &me.Id
		if request.ClaimedFirstAt == nil {
			request.ClaimedFirstAt = &now
		}
	} else if data.Action == "unclaim" {
		if request.ClaimedAt == nil {
			return &Error{
				C: http.StatusBadRequest, M: "can't unclaim an unclaimed request",
			}
		}

		request.ClaimedAt = nil
		request.ClaimedBy = nil
	} else if data.Action == "resolve" {
		if request.ClaimedAt == nil {
			return &Error{
				C: http.StatusBadRequest, M: "can't resolve an unclaimed request",
			}
		}

		request.ResolvedAt = &now
	} else if data.Action == "discard" {
		request.ResolvedAt = &now
		request.IsDiscarded = true
		if request.ClaimedAt == nil {
			request.ClaimedAt = &now
			request.ClaimedBy = &me.Id
			if request.ClaimedFirstAt == nil {
				request.ClaimedFirstAt = &now
			}
		}
	} else {
		return &Error{C: http.StatusBadRequest, M: "bad action"}
	}

	err = e.M.UpdateAttentionRequest(request, me.Id)
	if err != nil {
		return &Error{E: err}
	}

	if me.IsAdmin {
		return OK(request, c, w)
	}

	return OK(request.AttentionRequestPublic, c, w)
}
