package models

import (
	"time"
)

type UserGamePublic struct {
	Id                    string     `json:"id"`
	UserId                string     `db:"user_id" json:"userId"`
	GameId                string     `db:"game_id" json:"gameId"`
	Name                  *string    `json:"name"`
	Link                  *string    `json:"link"`
	CreatedAt             time.Time  `db:"created_at" json:"createdAt"`
	VerifiedAt            *time.Time `db:"verified_at" json:"verifiedAt"`
	UpdatedAt             *time.Time `db:"updated_at" json:"updatedAt"`
	DataUpdatedAt         *time.Time `db:"data_updated_at" json:"dataUpdatedAt"`
	DataUpdateRequestedAt *time.Time `db:"data_update_requested_at" json:"dataUpdateRequestedAt"`
	NullifiedAt           *time.Time `db:"nullified_at" json:"nullifiedAt"`
}

type UserGame struct {
	UserGamePublic
	Token       *string `json:"token"`
	Data        JSONMap `json:"data"`
	NullifiedBy *string `db:"nullified_by" json:"nullifiedBy"`
	UpdatedBy   *string `db:"updated_by" json:"-"`
}

func (e *Env) CreateUserGame(userGame *UserGame) error {
	return e.Db.Get(
		userGame, `
    INSERT INTO user_game (user_id, game_id, token, data)
    VALUES ($1, $2, $3, $4)
    RETURNING *`,
		userGame.UserId,
		userGame.GameId,
		userGame.Token,
		userGame.Data,
	)
}

func (e *Env) GetUserGameById(id string) (*UserGame, error) {
	var userGame UserGame
	err := e.Db.Get(
		&userGame, `
    SELECT *
    FROM user_game
    WHERE id=$1`,
		id,
	)
	return &userGame, BetterGetterErrors(err)
}

func (e *Env) GetUserGames(modifier *QueryModifier) ([]UserGame, error) {
	userGames := make([]UserGame, 0)
	sql, args, err := modifier.ToSql("user_game", "*")
	if err != nil {
		return userGames, err
	}

	err = e.Db.Select(&userGames, sql, args...)
	return userGames, BetterGetterErrors(err)
}

func (e *Env) UpdateUserGame(userGame *UserGame) error {
	return e.Db.Get(
		userGame, `
    UPDATE user_game
    SET
      token=$2,
      data=$3,
      name=$4,
      link=$5,
      verified_at=$6,
      data_updated_at=$7,
      data_update_requested_at=$8,
      nullified_at=$9,
      nullified_by=$10
    WHERE id=$1
    RETURNING *`,
		userGame.Id,
		userGame.Token,
		userGame.Data,
		userGame.Name,
		userGame.Link,
		userGame.VerifiedAt,
		userGame.DataUpdatedAt,
		userGame.DataUpdateRequestedAt,
		userGame.NullifiedAt,
		userGame.NullifiedBy,
	)
}

func (e *Env) GetUserGameByUserGame(userId, gameId string) (*UserGame, error) {
	var userGame UserGame
	err := e.Db.Get(
		&userGame, `
    SELECT *
    FROM user_game
    WHERE
      user_id=$1 AND
      game_id=$2 AND
      nullified_at IS NULL
    LIMIT 1`, // verified_at IS NOT NULL AND
		userId,
		gameId,
	)
	return &userGame, BetterGetterErrors(err)
}

func (e *Env) GetUserGameByHandleToken(
	handle, token string,
) (*UserGame, error) {
	var userGame UserGame
	err := e.Db.Get(
		&userGame, `
    SELECT user_game.*
    FROM user_game
    JOIN game ON user_game.game_id=game.id
    WHERE
      game.verification_handle=$1 AND
      user_game.token=$2
    LIMIT 1`, // TODO: research the need in LIMIT
		handle,
		token,
	)
	return &userGame, BetterGetterErrors(err)
}

func (e *Env) GetUserGameByGameData(
	gameId, dataKey, dataValue string,
) (*UserGame, error) {
	var userGame UserGame
	err := e.Db.Get(
		&userGame, `
    SELECT *
    FROM user_game
    WHERE
      game_id=$1 AND
      data->>$2=$3 AND
      nullified_at IS NULL
    LIMIT 1`,
		gameId,
		dataKey,
		dataValue,
	)
	return &userGame, BetterGetterErrors(err)
}

func (e *Env) GetUserGamesToUpdateByHandle(handle string) ([]UserGame, error) {
	userGames := make([]UserGame, 0)
	err := e.Db.Select(
		&userGames, `
    SELECT user_game.*
    FROM user_game
    JOIN game ON user_game.game_id=game.id
    WHERE
      game.verification_handle=$1 AND
      verified_at IS NOT NULL AND
      nullified_at IS NULL AND (
        user_game.data_update_requested_at IS NOT NULL OR
        user_game.data_updated_at < now() - interval '24 hours'
      )`,
		handle,
	)
	return userGames, BetterGetterErrors(err)
}
