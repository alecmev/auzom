package api

import (
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

func (e *Env) PostUserGame(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	var data struct {
		GameId string
	}
	err := Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	game, err := e.M.GetGameById(data.GameId)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest, M: "bad game id"}
	} else if game.VerificationHandle == nil ||
		*game.VerificationHandle != "battlefield-4" { // TODO: constant
		return &Error{
			C: http.StatusBadRequest,
			M: "this game's verification method isn't implemented yet",
		}
	}

	_, err = e.M.GetUserGameByUserGame(session.UserId, game.Id)
	if err == nil {
		return &Error{
			C: http.StatusBadRequest, M: "a valid user-game pair already exists",
		}
	} else if err != utils.ErrNotFound {
		return &Error{E: err, C: http.StatusInternalServerError}
	}

	token, err := utils.GenerateToken(5, true)
	if err != nil {
		return &Error{E: err}
	}

	userGame := &models.UserGame{
		UserGamePublic: models.UserGamePublic{
			UserId: session.UserId,
			GameId: game.Id,
		},
		Token: &token,
	}
	err = e.M.CreateUserGame(userGame)
	if err != nil {
		return &Error{E: err}
	}

	return Created(userGame, c, w)
}

func (e *Env) GetUserGame(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	userGame, err := e.M.GetUserGameById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.IsAdmin || me.Id == userGame.UserId {
			return OK(userGame, c, w)
		}
	}

	return OK(userGame.UserGamePublic, c, w)
}

func (e *Env) GetUserGames(
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
	gameId, _ := data.Filter["game_id"]
	if (userId == "") && (gameId == "") {
		return &Error{
			E: err, C: http.StatusBadRequest,
			M: "filter by user_id or game_id required",
		}
	}

	userGames, err := e.M.GetUserGames(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"user_id", "game_id"},
		[]string{},
	))
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.IsAdmin || me.Id == userId {
			return OK(userGames, c, w)
		}
	}

	public := make([]models.UserGamePublic, 0, len(userGames))
	for _, userGame := range userGames {
		public = append(public, userGame.UserGamePublic)
	}

	return OK(public, c, w)
}

func (e *Env) PatchUserGame(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	userGame, err := e.M.GetUserGameById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	} else if userGame.NullifiedAt != nil {
		return &Error{
			C: http.StatusBadRequest, M: "this user-game pair isn't valid anymore",
		}
	}

	var data struct {
		Action string
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	var me *models.User
	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err = e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		}
	}

	now := time.Now()
	if data.Action == "nullify" {
		if me != nil && me.Id != userGame.UserId && !me.IsAdmin {
			return &Error{E: utils.ErrUnauthorized}
		}

		userGame.Token = nil
		userGame.NullifiedAt = &now
		userGame.NullifiedBy = &me.Id
	} else if data.Action == "update" {
		if now.Before(userGame.DataUpdatedAt.Add(time.Minute)) {
			return &Error{
				C: http.StatusBadRequest,
				M: "data has been updated less than a minute ago",
			}
		} else if userGame.DataUpdateRequestedAt != nil {
			return &Error{
				C: http.StatusBadRequest,
				M: "data update has already been requested",
			}
		}

		userGame.DataUpdateRequestedAt = &now
	} else {
		return &Error{C: http.StatusBadRequest, M: "bad action"}
	}

	err = e.M.UpdateUserGame(userGame)
	if err != nil {
		return &Error{E: err}
	}

	if me != nil && (me.Id == userGame.UserId || me.IsAdmin) {
		return OK(userGame, c, w)
	}

	return OK(userGame.UserGamePublic, c, w)
}
