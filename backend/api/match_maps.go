package api

import (
	"net/http"

	"github.com/zenazn/goji/web"

	"app/models"
)

func (e *Env) GetMatchMap(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	matchMap, err := e.M.GetMatchMapById(c.URLParams["id"])
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
			return OK(matchMap, c, w)
		}
	}

	return OK(matchMap.MatchMapPublic, c, w)
}

func (e *Env) GetMatchMaps(
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

	matchMaps, err := e.M.GetMatchMaps(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"match_id", "game_map_id", "team_id", "is_ban", "discarded_at"},
		[]string{"id", "created_at"},
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
			return OK(matchMaps, c, w)
		}
	}

	public := make([]models.MatchMapPublic, 0)
	for _, matchMap := range matchMaps {
		public = append(public, matchMap.MatchMapPublic)
	}

	return OK(public, c, w)
}
