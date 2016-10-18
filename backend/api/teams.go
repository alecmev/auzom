package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

const ( // inspired by Steam
	nameMinLen = 2
	nameMaxLen = 64
	abbrMinLen = 1
	abbrMaxLen = 12
)

func (e *Env) PostTeam(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	var data struct {
		Name string
		Abbr string
		Logo string
	}
	err := Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	data.Name = strings.TrimSpace(data.Name)
	if len(data.Name) < nameMinLen || len(data.Name) > nameMaxLen {
		return &Error{C: http.StatusBadRequest, M: "invalid name length"}
	}

	data.Abbr = strings.TrimSpace(data.Abbr)
	if len(data.Abbr) < abbrMinLen || len(data.Abbr) > abbrMaxLen {
		return &Error{C: http.StatusBadRequest, M: "invalid abbreviation length"}
	}

	teams, err := e.M.GetTeams(models.NewQueryModifier(
		models.QueryBase{0, 0, map[string]string{
			"created_by":   session.UserId,
			"disbanded_at": "\x00",
		}, ""},
		[]string{"created_by", "disbanded_at"},
		[]string{},
	))
	if err != nil {
		return &Error{E: err}
	}
	if len(teams) >= 2 {
		return &Error{C: http.StatusBadRequest, M: "can't create more than 2 teams"}
	}

	team := &models.Team{
		TeamPublic: models.TeamPublic{
			Name:      data.Name,
			Abbr:      data.Abbr,
			Logo:      data.Logo,
			CreatedBy: session.UserId,
		},
	}
	err = e.M.Atomic(func(etx *models.Env) error {
		inerr := etx.CreateTeam(team)
		if inerr != nil {
			return &Error{E: inerr}
		}

		inerr = etx.CreateUserTeam(&models.UserTeam{
			UserTeamPublic: models.UserTeamPublic{
				UserId:   session.UserId,
				TeamId:   team.Id,
				IsLeader: true,
			},
		})
		if inerr != nil {
			return &Error{E: inerr}
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

	return Created(team, c, w)
}

func (e *Env) GetTeam(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	team, err := e.M.GetTeamById(c.URLParams["id"])
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
			return OK(team, c, w)
		}

		myUserTeam, err := e.M.GetUserTeamByUserTeam(me.Id, team.Id)
		if err == nil && myUserTeam.IsLeader {
			return OK(team, c, w)
		}
		if err != nil && err != utils.ErrNotFound {
			return &Error{E: err, C: http.StatusInternalServerError}
		}
	}

	return OK(team.TeamPublic, c, w)
}

func (e *Env) GetTeams(
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

	teams, err := e.M.GetTeams(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"created_by", "disbanded_at"},
		[]string{"id", "name", "abbr", "created_at", "created_by"},
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
			return OK(teams, c, w)
		}
	}

	public := make([]models.TeamPublic, 0)
	for _, team := range teams {
		public = append(public, team.TeamPublic)
	}

	return OK(public, c, w)
}

func (e *Env) PutTeam(
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

	team, err := e.M.GetTeamById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}
	if team.DisbandedAt != nil {
		return &Error{C: http.StatusBadRequest, M: "this team is disbanded"}
	}

	if !me.IsAdmin {
		myUserTeam, err := e.M.GetUserTeamByUserTeam(me.Id, team.Id)
		if err != nil && err != utils.ErrNotFound {
			return &Error{E: err, C: http.StatusInternalServerError}
		}
		if err == utils.ErrNotFound || !myUserTeam.IsLeader {
			return &Error{E: utils.ErrUnauthorized}
		}
	}

	var data struct {
		Name *string
		Abbr *string
		Logo *string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	somethingChanged := false

	if data.Name != nil {
		*data.Name = strings.TrimSpace(*data.Name)
		if *data.Name != team.Name {
			if len(*data.Name) < nameMinLen || len(*data.Name) > nameMaxLen {
				return &Error{C: http.StatusBadRequest, M: "invalid name length"}
			}

			team.Name = *data.Name
			somethingChanged = true
		}
	}

	if data.Abbr != nil {
		*data.Abbr = strings.TrimSpace(*data.Abbr)
		if *data.Abbr != team.Abbr {
			if len(*data.Abbr) < abbrMinLen || len(*data.Abbr) > abbrMaxLen {
				return &Error{
					C: http.StatusBadRequest, M: "invalid abbreviation length",
				}
			}

			team.Abbr = *data.Abbr
			somethingChanged = true
		}
	}

	if data.Logo != nil && *data.Logo != team.Logo {
		team.Logo = *data.Logo
		somethingChanged = true
	}

	if !somethingChanged {
		return OK(team, c, w)
	}

	err = e.M.UpdateTeam(team, me.Id)
	if err != nil {
		return &Error{E: err}
	}

	return OK(team, c, w)
}

func (e *Env) PatchTeam(
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
		return &Error{
			E: utils.ErrUnauthorized,
			M: "only admins can disband teams at the moment",
		}
	}

	team, err := e.M.GetTeamById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}
	if team.DisbandedAt != nil {
		return &Error{C: http.StatusBadRequest, M: "this team is disbanded"}
	}

	var data struct {
		Action string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}
	if data.Action != "disband" {
		return &Error{E: err, C: http.StatusBadRequest, M: "bad action"}
	}

	now := time.Now()
	team.DisbandedAt = &now
	team.DisbandedBy = &me.Id
	err = e.M.Atomic(func(etx *models.Env) error {
		teamSeasons, inerr := etx.GetTeamSeasons(models.NewQueryModifier(
			models.QueryBase{0, 0, map[string]string{
				"team_id": team.Id,
				"left_at": "\x00",
			}, ""},
			[]string{"team_id", "left_at"},
			[]string{},
		))
		if inerr != nil {
			return inerr
		}
		if len(teamSeasons) > 0 {
			return &Error{
				E: err, C: http.StatusBadRequest,
				M: "this team is in an active season",
			}
		}

		teamSeasonRequests, inerr := etx.GetTeamSeasonRequests(
			models.NewQueryModifier(
				models.QueryBase{0, 0, map[string]string{
					"team_id":  team.Id,
					"decision": "\x00",
				}, ""},
				[]string{"team_id", "decision"},
				[]string{},
			),
		)
		if inerr != nil {
			return inerr
		}

		decision := false
		for _, teamSeasonRequest := range teamSeasonRequests {
			teamSeasonRequest.Decision = &decision
			teamSeasonRequest.DecidedAt = &now
			teamSeasonRequest.DecidedBy = &me.Id
			inerr = etx.UpdateTeamSeasonRequest(&teamSeasonRequest)
			if inerr != nil {
				return inerr
			}
		}

		userTeams, inerr := etx.GetUserTeams(models.NewQueryModifier(
			models.QueryBase{0, 0, map[string]string{
				"team_id": team.Id,
				"left_at": "\x00",
			}, ""},
			[]string{"team_id", "left_at"},
			[]string{},
		))
		if inerr != nil {
			return inerr
		}

		for _, userTeam := range userTeams {
			userTeam.IsLeader = false
			userTeam.LeftAt = &now
			userTeam.KickedBy = &me.Id
			inerr = etx.UpdateUserTeam(&userTeam, me.Id)
			if inerr != nil {
				return inerr
			}
		}

		userTeamRequests, inerr := etx.GetUserTeamRequests(models.NewQueryModifier(
			models.QueryBase{0, 0, map[string]string{
				"team_id":  team.Id,
				"decision": "\x00",
			}, ""},
			[]string{"team_id", "decision"},
			[]string{},
		))
		if inerr != nil {
			return inerr
		}

		for _, userTeamRequest := range userTeamRequests {
			userTeamRequest.Decision = &decision
			userTeamRequest.DecidedAt = &now
			userTeamRequest.AdminDecision = &decision
			userTeamRequest.AdminDecidedAt = &now
			userTeamRequest.AdminDecidedBy = &me.Id
			inerr = etx.UpdateUserTeamRequest(&userTeamRequest)
			if inerr != nil {
				return inerr
			}
		}

		inerr = etx.UpdateTeam(team, me.Id)
		if inerr != nil {
			return inerr
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

	return OK(team, c, w)
}
