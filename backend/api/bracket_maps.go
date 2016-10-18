package api

import (
	"net/http"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) PostBracketMap(
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

	var data struct {
		BracketId string
		GameMapId string
		SubPool   int `json:",string"`
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	} else if data.SubPool < 0 || data.SubPool > 9 {
		return &Error{C: http.StatusBadRequest, M: "bad sub-pool, must be [0, 9]"}
	}

	bracket, err := e.M.GetBracketById(data.BracketId)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	stage, err := e.M.GetStageById(bracket.StageId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}

	season, err := e.M.GetSeasonById(stage.SeasonId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}

	tournament, err := e.M.GetTournamentById(season.TournamentId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	}

	gameMap, err := e.M.GetGameMapById(data.GameMapId)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	} else if gameMap.GameId != tournament.GameId {
		return &Error{
			C: http.StatusBadRequest, M: "map game does not match tournament game",
		}
	}

	bracketMap := &models.BracketMap{
		BracketMapPublic: models.BracketMapPublic{
			BracketId: data.BracketId,
			GameMapId: data.GameMapId,
			SubPool:   data.SubPool,
		},
		CreatedBy: me.Id,
	}
	err = e.M.CreateBracketMap(bracketMap)
	if err != nil {
		return &Error{E: err}
	}

	return Created(bracketMap, c, w)
}

func (e *Env) GetBracketMap(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	bracketMap, err := e.M.GetBracketMapById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.IsAdmin {
			return OK(bracketMap, c, w)
		}
	}

	return OK(bracketMap.BracketMapPublic, c, w)
}

func (e *Env) GetBracketMaps(
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

	bracketMaps, err := e.M.GetBracketMaps(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"bracket_id", "game_map_id", "sub_pool", "is_enabled"},
		[]string{"id", "bracket_id", "game_map_id", "sub_pool"},
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
			return OK(bracketMaps, c, w)
		}
	}

	public := make([]models.BracketMapPublic, 0, len(bracketMaps))
	for _, bracketMap := range bracketMaps {
		public = append(public, bracketMap.BracketMapPublic)
	}

	return OK(public, c, w)
}
