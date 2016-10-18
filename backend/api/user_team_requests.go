package api

import (
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) PostUserTeamRequest(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	var data struct {
		UserId string `json:"userId"`
		TeamId string `json:"teamId"`
	}
	err := Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	user, err := e.M.GetUserById(data.UserId)
	if err == utils.ErrNotFound {
		return &Error{C: http.StatusBadRequest, M: "bad user id"}
	} else if err != nil {
		return &Error{E: err}
	}

	team, err := e.M.GetTeamById(data.TeamId)
	if err == utils.ErrNotFound {
		return &Error{C: http.StatusBadRequest, M: "bad team id"}
	} else if err != nil {
		return &Error{E: err}
	} else if team.DisbandedAt != nil {
		return &Error{C: http.StatusBadRequest, M: "this team is disbanded"}
	}

	request := &models.UserTeamRequest{
		UserId: user.Id,
		TeamId: team.Id,
	}
	var needAdmin bool
	err = e.M.Atomic(func(etx *models.Env) error {
		// userTeams are checked along with requests, because it's possible to join
		// a team without a request, e.g. when you're the creator
		_, inerr := etx.GetUserTeamByUserTeam(user.Id, team.Id)
		if inerr == nil {
			return &Error{C: http.StatusBadRequest, M: "already a member"}
		}
		if inerr != utils.ErrNotFound {
			return inerr
		}

		_, inerr = etx.GetUserTeamRequestByUserTeam(user.Id, team.Id)
		if inerr == nil {
			return &Error{C: http.StatusBadRequest, M: "already requested"}
		}
		if inerr != utils.ErrNotFound {
			return inerr
		}

		teamsUserIn, inerr := etx.GetUserTeamRelationsByUser(user)
		if inerr != nil {
			return inerr
		}

		teamsRelated, inerr := etx.GetTeamSeasonRelationsByTeam(team)
		if inerr != nil {
			return inerr
		}

		teamsUserInMap := make(map[string]struct{})
		for _, teamUserIn := range teamsUserIn {
			teamsUserInMap[teamUserIn.TeamId] = struct{}{}
		}

		for _, teamRelated := range teamsRelated {
			_, ok = teamsUserInMap[teamRelated.TeamId]
			if ok {
				return &Error{
					C: http.StatusBadRequest,
					M: "user is involved with team " + teamRelated.TeamId +
						" which is involved in the same season " + teamRelated.SeasonId +
						" as this team",
				}
			}
		}

		teamSeasons, inerr := etx.GetTeamSeasons(models.NewQueryModifier(
			models.QueryBase{0, 0, map[string]string{"team_id": team.Id}, ""},
			[]string{"team_id"},
			[]string{},
		))
		if inerr != nil {
			return inerr
		}

		for _, teamSeason := range teamSeasons {
			if teamSeason.LeftAt == nil {
				needAdmin = true
				break
			}
		}

		tmpTrue := true
		if !needAdmin {
			request.AdminDecision = &tmpTrue
		}

		// even if me.IsAdmin == true, explicit approval is better

		if session.UserId == user.Id {
			request.UserDecision = &tmpTrue
		} else {
			myUserTeam, inerr := etx.GetUserTeamByUserTeam(session.UserId, team.Id)
			if inerr == utils.ErrNotFound || (inerr == nil && !myUserTeam.IsLeader) {
				return &Error{E: utils.ErrUnauthorized, M: "you aren't a leader"}
			}
			if inerr != nil {
				return inerr
			}

			request.LeaderDecision = &tmpTrue
			request.LeaderDecidedBy = &myUserTeam.UserId
		}

		return etx.CreateUserTeamRequest(request)
	})
	if err != nil {
		apierr, ok := err.(*Error)
		if ok {
			return apierr
		}

		return &Error{E: err}
	}

	if needAdmin {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		}

		e.Slack.Send("user-team-requests", me,
			"*Target user:* <https://"+e.StaticHost+"/users/"+user.Id+"|"+
				user.Nickname+">\n"+
				"*Target team:* <https://"+e.StaticHost+"/teams/"+team.Id+"|"+
				team.Name+">",
		)
	}

	return Created(request, c, w)
}

func (e *Env) GetUserTeamRequest(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	request, err := e.M.GetUserTeamRequestById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}
	if me.IsAdmin || me.Id == request.UserId {
		return OK(request, c, w)
	}

	myUserTeam, err := e.M.GetUserTeamByUserTeam(me.Id, request.TeamId)
	if err == utils.ErrNotFound || (err == nil && !myUserTeam.IsLeader) {
		return &Error{E: utils.ErrUnauthorized}
	}
	if err != nil {
		return &Error{E: err}
	}

	return OK(request, c, w)
}

func (e *Env) GetUserTeamRequests(
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

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}

	userId, _ := data.Filter["user_id"]
	teamId, _ := data.Filter["team_id"]
	if (userId == "") && (teamId == "") {
		return &Error{
			C: http.StatusBadRequest, M: "filter by user_id or team_id required",
		}
	}

	amRelated := false
	if me.IsAdmin {
		amRelated = true
	} else if me.Id == userId {
		amRelated = true
	} else if teamId != "" {
		myUserTeam, err := e.M.GetUserTeamByUserTeam(me.Id, teamId)
		if err == nil && myUserTeam.IsLeader {
			amRelated = true
		}
	}

	if !amRelated {
		return &Error{E: utils.ErrUnauthorized}
	}

	requests, err := e.M.GetUserTeamRequests(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"user_id", "team_id", "decision"},
		[]string{"id", "user_id", "team_id", "created_at"},
	))
	if err != nil {
		return &Error{E: err}
	}

	return OK(requests, c, w)
}

func (e *Env) PatchUserTeamRequest(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	request, err := e.M.GetUserTeamRequestById(c.URLParams["id"])
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
	if data.Action == "yes" {
		action = true
	} else if data.Action != "no" {
		return &Error{C: http.StatusBadRequest, M: "bad action"}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}

	now := time.Now()
	if request.IsAdminNeeded() && me.IsAdmin {
		request.AdminDecision = &action
		request.AdminDecidedAt = &now
		request.AdminDecidedBy = &me.Id
	}

	if me.Id == request.UserId {
		request.UserDecision = &action
		request.UserDecidedAt = &now
	} else {
		amAuthorized := true
		myUserTeam, err := e.M.GetUserTeamByUserTeam(me.Id, request.TeamId)
		if err != nil && err != utils.ErrNotFound {
			return &Error{E: err}
		}
		if err == utils.ErrNotFound || !myUserTeam.IsLeader {
			amAuthorized = false
		}

		if amAuthorized {
			request.LeaderDecision = &action
			request.LeaderDecidedAt = &now
			request.LeaderDecidedBy = &me.Id
		} else if !request.IsAdminNeeded() || !me.IsAdmin {
			return &Error{
				E: utils.ErrUnauthorized,
				M: "you aren't involved in this request",
			}
		}
	}

	// unfortunately, this is the cleanest way to not leak decision data and avoid
	// code duplication / unnecessary complexion
	if request.Decision != nil {
		return &Error{C: http.StatusBadRequest, M: "already decided"}
	}

	if !action {
		request.Decision = &action
		request.DecidedAt = &now
	} else if request.UserDecision != nil &&
		request.LeaderDecision != nil &&
		request.AdminDecision != nil &&
		true {
		request.Decision = &action
		request.DecidedAt = &now
		// TODO: in this case, the client probably doesn't need the userTeam
		// itself (it can deduct that it was created from the value of decision),
		// but there needs to be a mechanism for including additional responses
		err = e.M.CreateUserTeam(&models.UserTeam{
			UserTeamPublic: models.UserTeamPublic{
				UserId:    request.UserId,
				TeamId:    request.TeamId,
				RequestId: &request.Id,
			},
		})
		if err != nil {
			return &Error{E: err}
		}
	}

	err = e.M.UpdateUserTeamRequest(request)
	if err != nil {
		return &Error{E: err}
	}

	return OK(request, c, w)
}
