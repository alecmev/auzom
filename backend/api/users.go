package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/elithrar/simple-scrypt"
	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

const (
	pwdMinLen = 8
)

func (e *Env) PostUser(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	var data struct {
		Email    string
		Password string
	}
	err := Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}
	if len(data.Password) < pwdMinLen {
		return &Error{C: http.StatusBadRequest, M: "password too short"}
	}

	data.Email = strings.TrimSpace(data.Email)
	_, err = e.M.GetUserByEmail(data.Email)
	if err == nil {
		return &Error{C: http.StatusBadRequest, M: "email taken"}
	}
	if err != utils.ErrNotFound {
		return &Error{E: err}
	}

	password, err := scrypt.GenerateFromPassword(
		[]byte(data.Password),
		e.Scrypt,
	)
	if err != nil {
		return &Error{E: err}
	}

	user := &models.User{
		Email:    data.Email,
		Password: password,
	}
	err = e.M.Atomic(func(etx *models.Env) error {
		inerr := etx.CreateUser(user)
		if inerr != nil {
			return &Error{E: inerr}
		}

		otp, inerr := etx.CreateOTP(user.Id)
		if inerr != nil {
			return &Error{E: inerr}
		}

		inerr = e.Mail.Send(
			data.Email,
			"Please verify your email address",
			"[Verify email address](https://"+e.StaticHost+"/verify/"+otp.Token+")",
		)
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

	return Created(user, c, w)
}

func (e *Env) GetUser(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	user, err := e.M.GetUserById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	session, ok := c.Env["session"].(*models.Session)
	if ok {
		me, err := e.M.GetUserById(session.UserId)
		if err != nil {
			return &Error{E: err, C: http.StatusInternalServerError}
		} else if me.Id == user.Id || me.IsAdmin {
			return OK(user, c, w)
		}
	}

	return OK(user.UserPublic, c, w)
}

func (e *Env) GetUsers(
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

	users, err := e.M.GetUsers(models.NewQueryModifier(
		models.QueryBase{data.Offset, data.Limit, data.Filter, data.Sort},
		[]string{"is_admin", "is_email_verified"},
		[]string{"id", "email", "nickname", "fullname", "is_admin",
			"is_email_verified", "created_at"},
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
			return OK(users, c, w)
		}
	}

	public := make([]models.UserPublic, 0)
	for _, user := range users {
		public = append(public, user.UserPublic)
	}

	return OK(public, c, w)
}

func (e *Env) PutUser(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, ok := c.Env["session"].(*models.Session)
	if !ok {
		return &Error{E: utils.ErrUnauthorized}
	}

	me, err := e.M.GetUserById(session.UserId)
	if err != nil {
		return &Error{E: err, C: http.StatusInternalServerError}
	} else if c.URLParams["id"] != session.UserId && !me.IsAdmin {
		return &Error{E: utils.ErrUnauthorized}
	}

	user, err := e.M.GetUserById(c.URLParams["id"])
	if err != nil {
		return &Error{E: err}
	}

	var data struct {
		Email         *string
		Password      *string
		Nickname      *string
		Fullname      *string
		GravatarEmail *string
		IsAdmin       *bool
	}
	err = Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	} else if data.Email != nil || data.Password != nil {
		if session.CreatedAt.Add(SessionAgeSensitive).Before(time.Now()) {
			return &Error{
				E: utils.ErrUnauthorized,
				M: "please, re-login, so you can change your email/password",
			}
		}
	}

	somethingChanged := false
	emailChanged := false

	if data.Email != nil {
		*data.Email = strings.TrimSpace(*data.Email)
		if *data.Email != user.Email {
			_, err = e.M.GetUserByEmail(*data.Email)
			if err == nil {
				return &Error{C: http.StatusBadRequest, M: "email taken"}
			} else if err != utils.ErrNotFound {
				return &Error{E: err}
			}

			user.Email = *data.Email
			user.IsEmailVerified = false
			somethingChanged = true
			emailChanged = true
		}
	}

	if data.Password != nil {
		if len(*data.Password) < pwdMinLen {
			return &Error{C: http.StatusBadRequest, M: "password too short"}
		}

		password, err := scrypt.GenerateFromPassword(
			[]byte(*data.Password),
			e.Scrypt,
		)
		if err != nil {
			return &Error{E: err}
		}

		user.Password = password
		somethingChanged = true
	}

	if data.Nickname != nil {
		*data.Nickname = strings.TrimSpace(*data.Nickname)
		if *data.Nickname != user.Nickname {
			if len(*data.Nickname) == 0 { // redundant, but explicit
				return &Error{C: http.StatusBadRequest, M: "nickname too short"}
			}

			user.Nickname = *data.Nickname
			somethingChanged = true
		}
	}

	if data.Fullname != nil {
		*data.Fullname = strings.TrimSpace(*data.Fullname)
		if *data.Fullname != user.Fullname {
			user.Fullname = *data.Fullname
			somethingChanged = true
		}
	}

	if data.GravatarEmail != nil {
		*data.GravatarEmail = strings.TrimSpace(*data.GravatarEmail)
		if *data.GravatarEmail != user.GravatarEmail {
			user.GravatarEmail = *data.GravatarEmail
			somethingChanged = true
		}
	}

	if data.IsAdmin != nil && *data.IsAdmin != user.IsAdmin {
		if !me.IsAdmin {
			return &Error{E: utils.ErrUnauthorized}
		}

		user.IsAdmin = *data.IsAdmin
		somethingChanged = true
	}

	if !somethingChanged {
		return OK(user, c, w)
	}

	err = e.M.Atomic(func(etx *models.Env) error {
		inerr := etx.UpdateUser(user, me.Id)
		if inerr != nil {
			return &Error{E: inerr}
		} else if !emailChanged {
			return nil
		}

		otp, inerr := etx.CreateOTP(user.Id)
		if inerr != nil {
			return &Error{E: inerr}
		}

		inerr = e.Mail.Send(
			user.Email,
			"Please verify your email address",
			"[Verify email address](https://"+e.StaticHost+"/verify/"+otp.Token+")",
		)
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

	return OK(user, c, w)
}
