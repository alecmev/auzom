package api

import (
	"net/http"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) PostGameMap(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: utils.ErrUnauthorized}
	} else if !me.IsAdmin {
		return &Error{E: utils.ErrUnauthorized}
	}

	var data struct {
		GameId    string
		Name      string
		Abbr      string
		SideX     string
		SideXAbbr string
		SideY     string
		SideYAbbr string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	gameMap := &models.GameMap{
		GameMapPublic: models.GameMapPublic{
			GameId:    data.GameId,
			Name:      data.Name,
			Abbr:      data.Abbr,
			SideX:     data.SideX,
			SideXAbbr: data.SideXAbbr,
			SideY:     data.SideY,
			SideYAbbr: data.SideYAbbr,
		},
		CreatedBy: me.Id,
	}
	err = e.M.CreateGameMap(gameMap)
	if err != nil {
		return &Error{E: err}
	}

	return Created(gameMap, c, w)
}

func (e *Env) GetGameMap(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	gameMap, err := e.M.GetGameMapById(c.URLParams["id"])
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
			return OK(gameMap, c, w)
		}
	}

	return OK(gameMap.GameMapPublic, c, w)
}

func (e *Env) GetGameMaps(
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

	gameMaps, err := e.M.GetGameMaps(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"game_id"},
		[]string{"id", "game_id", "name", "abbr"},
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
			return OK(gameMaps, c, w)
		}
	}

	public := make([]models.GameMapPublic, 0)
	for _, gameMap := range gameMaps {
		public = append(public, gameMap.GameMapPublic)
	}

	return OK(public, c, w)
}

func (e *Env) PutGameMap(
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

	gameMap, err := e.M.GetGameMapById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	var data struct {
		Name      *string
		Abbr      *string
		SideX     *string
		SideXAbbr *string
		SideY     *string
		SideYAbbr *string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	somethingChanged := false

	if data.Name != nil && *data.Name != gameMap.Name {
		gameMap.Name = *data.Name
		somethingChanged = true
	}

	if data.Abbr != nil && *data.Abbr != gameMap.Abbr {
		gameMap.Abbr = *data.Abbr
		somethingChanged = true
	}

	if data.SideX != nil && *data.SideX != gameMap.SideX {
		gameMap.SideX = *data.SideX
		somethingChanged = true
	}

	if data.SideXAbbr != nil && *data.SideXAbbr != gameMap.SideXAbbr {
		gameMap.SideXAbbr = *data.SideXAbbr
		somethingChanged = true
	}

	if data.SideY != nil && *data.SideY != gameMap.SideY {
		gameMap.SideY = *data.SideY
		somethingChanged = true
	}

	if data.SideYAbbr != nil && *data.SideYAbbr != gameMap.SideYAbbr {
		gameMap.SideYAbbr = *data.SideYAbbr
		somethingChanged = true
	}

	if !somethingChanged {
		return OK(gameMap, c, w)
	}

	err = e.M.UpdateGameMap(gameMap, me.Id)
	if err != nil {
		return &Error{E: err}
	}

	return OK(gameMap, c, w)
}
