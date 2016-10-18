package api

import (
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) GetTeamSeason(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	teamSeason, err := e.M.GetTeamSeasonById(c.URLParams["id"])
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
			return OK(teamSeason, c, w)
		}

		myUserTeam, err := e.M.GetUserTeamByUserTeam(me.Id, teamSeason.TeamId)
		if err == nil && myUserTeam.IsLeader {
			return OK(teamSeason, c, w)
		}
		if err != nil && err != utils.ErrNotFound {
			return &Error{E: err, C: http.StatusInternalServerError}
		}
	}

	return OK(teamSeason.TeamSeasonPublic, c, w)
}

func (e *Env) GetTeamSeasons(
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

	teamId, _ := data.Filter["team_id"]
	seasonId, _ := data.Filter["season_id"]
	requestId, _ := data.Filter["request_id"]
	if (teamId == "") && (seasonId == "") && (requestId == "") {
		return &Error{
			E: err, C: http.StatusBadRequest,
			M: "filter by team_id, season_id or request_id required",
		}
	}

	teamSeasons, err := e.M.GetTeamSeasons(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"team_id", "season_id", "is_done"},
		[]string{"id", "team_id", "season_id", "is_done", "created_at", "left_at"},
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
			return OK(teamSeasons, c, w)
		}

		// TODO: return w/private info, if filtered by team_id and I'm a leader
	}

	public := make([]models.TeamSeasonPublic, 0)
	for _, teamSeason := range teamSeasons {
		public = append(public, teamSeason.TeamSeasonPublic)
	}

	return OK(public, c, w)
}

func (e *Env) PatchTeamSeason(
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
		Action string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	teamSeason, err := e.M.GetTeamSeasonById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}
	if teamSeason.LeftAt != nil {
		return &Error{C: http.StatusBadRequest, M: "this team isn't a participant"}
	}

	now := time.Now()
	teamSeason.LeftAt = &now
	if data.Action == "leave" {
		//
	} else if data.Action == "done" {
		teamSeason.IsDone = true
	} else if data.Action == "kick" {
		teamSeason.KickedBy = &me.Id
	} else {
		return &Error{E: err, C: http.StatusBadRequest, M: "bad action"}
	}

	err = e.M.UpdateTeamSeason(teamSeason, me.Id)
	if err != nil {
		return &Error{E: err}
	}

	return OK(teamSeason, c, w)
}
