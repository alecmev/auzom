package api

import (
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) PostTeamSeasonRequest(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	var data struct {
		TeamId   string `json:"teamId"`
		SeasonId string `json:"seasonId"`
	}
	err := Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	team, err := e.M.GetTeamById(data.TeamId)
	if err == utils.ErrNotFound {
		return &Error{C: http.StatusBadRequest, M: "bad team id"}
	}
	if err != nil {
		return &Error{E: err}
	}
	if team.DisbandedAt != nil {
		return &Error{C: http.StatusBadRequest, M: "this team is disbanded"}
	}

	myUserTeam, err := e.M.GetUserTeamByUserTeam(session.UserId, team.Id)
	if err == utils.ErrNotFound || (err == nil && !myUserTeam.IsLeader) {
		return &Error{E: utils.ErrUnauthorized, M: "you aren't a leader"}
	}
	if err != nil {
		return &Error{E: err}
	}

	season, err := e.M.GetSeasonById(data.SeasonId)
	if err == utils.ErrNotFound {
		return &Error{C: http.StatusBadRequest, M: "bad season id"}
	}
	if err != nil {
		return &Error{E: err}
	}

	now := time.Now()
	if season.PublishedAt == nil || season.PublishedAt.After(now) {
		return &Error{C: http.StatusNotFound}
	}
	if season.SignupsOpenedAt == nil || season.SignupsOpenedAt.After(now) {
		return &Error{C: http.StatusBadRequest, M: "signups aren't open yet"}
	}
	if season.SignupsClosedAt != nil && season.SignupsClosedAt.Before(now) {
		return &Error{C: http.StatusBadRequest, M: "signups are over"}
	}

	request := &models.TeamSeasonRequest{
		TeamId:    team.Id,
		SeasonId:  season.Id,
		CreatedBy: session.UserId,
	}
	err = e.M.Atomic(func(etx *models.Env) error {
		_, inerr := etx.GetTeamSeasonByTeamSeason(team.Id, season.Id)
		if inerr == nil {
			return &Error{C: http.StatusBadRequest, M: "already a participant"}
		}
		if inerr != utils.ErrNotFound {
			return inerr
		}

		_, inerr = etx.GetTeamSeasonRequestByTeamSeason(team.Id, season.Id)
		if inerr == nil {
			return &Error{C: http.StatusBadRequest, M: "already requested"}
		}
		if inerr != utils.ErrNotFound {
			return inerr
		}

		teamsInSeason, inerr := etx.GetTeamSeasonRelationsBySeason(season)
		if inerr != nil {
			return inerr
		}

		teamsRelated, inerr := etx.GetUserTeamRelationsByTeam(team)
		if inerr != nil {
			return inerr
		}

		teamsInSeasonMap := make(map[string]struct{})
		for _, teamInSeason := range teamsInSeason {
			teamsInSeasonMap[teamInSeason.TeamId] = struct{}{}
		}

		for _, teamRelated := range teamsRelated {
			_, ok = teamsInSeasonMap[teamRelated.TeamId]
			if ok {
				return &Error{
					C: http.StatusBadRequest,
					M: "user " + teamRelated.UserId +
						" is involved with team " + teamRelated.TeamId +
						" which is involved in this season",
				}
			}
		}

		return etx.CreateTeamSeasonRequest(request)
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	return Created(request, c, w)
}

func (e *Env) GetTeamSeasonRequest(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	request, err := e.M.GetTeamSeasonRequestById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	return OK(request, c, w)
}

func (e *Env) GetTeamSeasonRequests(
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

	requests, err := e.M.GetTeamSeasonRequests(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"team_id", "season_id", "decision"},
		[]string{"id", "team_id", "season_id", "created_at"},
	))
	if err != nil {
		return &Error{E: err}
	}

	return OK(requests, c, w)
}

func (e *Env) PatchTeamSeasonRequest(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	request, err := e.M.GetTeamSeasonRequestById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	var data struct {
		Action string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	var action bool
	now := time.Now()
	if data.Action == "yes" {
		action = true
		season, err := e.M.GetSeasonById(request.SeasonId)
		if err != nil {
			return &Error{E: err}
		}
		if season.SignupsClosedAt == nil || season.SignupsClosedAt.After(now) {
			return &Error{
				E: err, C: http.StatusBadRequest,
				M: "can't accept an application while the signups are open",
			}
		}
	} else if data.Action != "no" && data.Action != "cancel" {
		return &Error{E: err, C: http.StatusBadRequest, M: "bad action"}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}

	request.Decision = &action
	request.DecidedAt = &now
	if data.Action == "cancel" {
		myUserTeam, err := e.M.GetUserTeamByUserTeam(me.Id, request.TeamId)
		if err == utils.ErrNotFound || (err == nil && !myUserTeam.IsLeader) {
			return &Error{E: utils.ErrUnauthorized}
		} else if err != nil {
			return &Error{E: err}
		}

		request.CancelledBy = &me.Id
	} else {
		if !me.IsAdmin {
			return &Error{
				E: err, C: http.StatusBadRequest,
				M: "only admins can accept / decline applications",
			}
		}

		request.DecidedBy = &me.Id
	}

	err = e.M.Atomic(func(etx *models.Env) error {
		if !action {
			return etx.UpdateTeamSeasonRequest(request)
		}

		return PatchTeamSeasonRequestHelper(etx, request)
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	return OK(request, c, w)
}

func PatchTeamSeasonRequestHelper(
	eM *models.Env, request *models.TeamSeasonRequest,
) error {
	userTeamRequests, inerr := eM.GetUserTeamRequests(models.NewQueryModifier(
		models.QueryBase{0, 0, map[string]string{
			"team_id":  request.TeamId,
			"decision": "\x00",
		}, ""},
		[]string{"team_id", "decision"},
		[]string{},
	))
	if inerr != nil {
		return inerr
	}

	for _, userTeamRequest := range userTeamRequests {
		if userTeamRequest.AdminDecidedAt != nil {
			continue
		}

		userTeamRequest.AdminDecision = nil
		inerr = eM.UpdateUserTeamRequest(&userTeamRequest)
		if inerr != nil {
			return inerr
		}
	}

	inerr = eM.CreateTeamSeason(&models.TeamSeason{
		TeamSeasonPublic: models.TeamSeasonPublic{
			TeamId:    request.TeamId,
			SeasonId:  request.SeasonId,
			RequestId: &request.Id,
		},
	})
	if inerr != nil {
		return inerr
	}

	return eM.UpdateTeamSeasonRequest(request)
}
