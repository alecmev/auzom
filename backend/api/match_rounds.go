package api

import (
	"net/http"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) GetMatchRound(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	round, err := e.M.GetMatchRoundById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.IsAdmin {
			return OK(round, c, w)
		}

		report, err := e.M.GetMatchReportById(round.MatchReportId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		}

		match, err := e.M.GetMatchById(report.MatchId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		}

		teamAmLeaderOf, err, status := match.UserIsLeaderOf(e.M, me.Id)
		if err != nil {
			return &Error{E: err, C: status}
		} else if teamAmLeaderOf != nil {
			return OK(round, c, w)
		}
	}

	return OK(round.MatchRoundPublic, c, w)
}

func (e *Env) GetMatchRounds(
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

	reportId, _ := data.Filter["match_report_id"]
	if reportId == "" {
		return &Error{
			C: http.StatusBadRequest, M: "match_report_id filter is mandatory",
		}
	}

	rounds, err := e.M.GetMatchRounds(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"match_report_id"},
		[]string{"id", "match_id", "created_at"},
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
			return OK(rounds, c, w)
		}

		report, err := e.M.GetMatchReportById(reportId)
		if err == utils.ErrNotFound {
			return &Error{C: http.StatusBadRequest, M: "bad report id"}
		} else if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		}

		match, err := e.M.GetMatchById(report.MatchId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		}

		teamAmLeaderOf, err, status := match.UserIsLeaderOf(e.M, me.Id)
		if err != nil {
			return &Error{E: err, C: status}
		} else if teamAmLeaderOf != nil {
			return OK(rounds, c, w)
		}
	}

	public := make([]models.MatchRoundPublic, 0)
	for _, round := range rounds {
		public = append(public, round.MatchRoundPublic)
	}

	return OK(public, c, w)
}
