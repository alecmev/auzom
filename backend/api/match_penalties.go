package api

import (
	"net/http"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) GetMatchPenalty(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	penalty, err := e.M.GetMatchPenaltyById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.IsAdmin {
			return OK(penalty, c, w)
		}

		report, err := e.M.GetMatchReportById(penalty.MatchReportId)
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
			return OK(penalty, c, w)
		}
	}

	return OK(penalty.MatchPenaltyPublic, c, w)
}

func (e *Env) GetMatchPenalties(
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

	penalties, err := e.M.GetMatchPenalties(models.NewQueryModifier(
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
			return OK(penalties, c, w)
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
			return OK(penalties, c, w)
		}
	}

	public := make([]models.MatchPenaltyPublic, 0)
	for _, penalty := range penalties {
		public = append(public, penalty.MatchPenaltyPublic)
	}

	return OK(public, c, w)
}
