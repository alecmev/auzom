package api

import (
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) GetUserTeam(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	userTeam, err := e.M.GetUserTeamById(c.URLParams["id"])
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
			return OK(userTeam, c, w)
		}

		myUserTeam, err := e.M.GetUserTeamByUserTeam(me.Id, userTeam.TeamId)
		if err == nil && myUserTeam.IsLeader {
			return OK(userTeam, c, w)
		}
		if err != nil && err != utils.ErrNotFound {
			return &Error{E: err, C: http.StatusInternalServerError}
		}
	}

	return OK(userTeam.UserTeamPublic, c, w)
}

func (e *Env) GetUserTeams(
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

	userId, _ := data.Filter["user_id"]
	teamId, _ := data.Filter["team_id"]
	requestId, _ := data.Filter["request_id"]

	if (userId == "") && (teamId == "") && (requestId == "") {
		// It is not allowed to provide no filter.
		return &Error{
			E: err, C: http.StatusBadRequest,
			M: "filter by user_id, team_id or request_id required",
		}
	}

	userTeams, err := e.M.GetUserTeams(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"user_id", "team_id", "is_leader", "request_id", "left_at"},
		[]string{"id", "user_id", "team_id", "is_leader", "created_at", "left_at"},
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
			return OK(userTeams, c, w)
		}
	}

	public := make([]models.UserTeamPublic, 0)
	for _, userTeam := range userTeams {
		public = append(public, userTeam.UserTeamPublic)
	}

	return OK(public, c, w)
}

func (e *Env) PatchUserTeam(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	userTeam, err := e.M.GetUserTeamById(c.URLParams["id"])
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

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}

	if data.Action == "leave" {
		if userTeam.UserId != me.Id { // yep, the admin can only kick
			return &Error{E: utils.ErrUnauthorized}
		}
	} else if !me.IsAdmin { // it's okay that this is a "catch-all", I think
		myUserTeam, err := e.M.GetUserTeamByUserTeam(me.Id, userTeam.TeamId)
		if err == utils.ErrNotFound || (err == nil && !myUserTeam.IsLeader) {
			return &Error{E: utils.ErrUnauthorized}
		}
		if err != nil {
			return &Error{E: err}
		}
	}

	// this is checked only now as to not give up confidential info
	if userTeam.LeftAt != nil {
		return &Error{C: http.StatusBadRequest, M: "this user isn't a member"}
	}

	var isLastLeader bool
	err = e.M.Atomic(func(etx *models.Env) error {
		if userTeam.IsLeader {
			userTeams, inerr := etx.GetUserTeams(models.NewQueryModifier(
				models.QueryBase{0, 0, map[string]string{
					"team_id":   userTeam.TeamId,
					"is_leader": "true",
					"left_at":   "\x00",
				}, ""},
				[]string{"team_id", "is_leader", "left_at"},
				[]string{},
			))
			if inerr != nil {
				return inerr
			}

			isLastLeader = len(userTeams) == 1
		}

		now := time.Now()
		if data.Action == "leave" {
			if isLastLeader {
				return &Error{
					C: http.StatusBadRequest,
					M: "last leader can't leave",
				}
			}

			userTeam.IsLeader = false // just in case
			userTeam.LeftAt = &now
		} else if data.Action == "promote" {
			if userTeam.IsLeader {
				return &Error{
					C: http.StatusBadRequest,
					M: "this member is already a leader",
				}
			}

			userTeam.IsLeader = true
		} else if data.Action == "demote" {
			if !userTeam.IsLeader {
				return &Error{
					C: http.StatusBadRequest,
					M: "this member isn't a leader",
				}
			}
			if isLastLeader {
				return &Error{
					C: http.StatusBadRequest,
					M: "last leader can't be demoted",
				}
			}

			userTeam.IsLeader = false
		} else if data.Action == "kick" {
			if isLastLeader {
				return &Error{
					C: http.StatusBadRequest,
					M: "last leader can't be kicked",
				}
			}

			userTeam.IsLeader = false // just in case
			userTeam.LeftAt = &now
			userTeam.KickedBy = &me.Id
		} else {
			return &Error{C: http.StatusBadRequest, M: "bad action"}
		}

		return etx.UpdateUserTeam(userTeam, me.Id)
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	return OK(userTeam, c, w)
}
