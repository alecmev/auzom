package api

import (
	"errors"
	"net/http"

	"github.com/elithrar/simple-scrypt"
	"github.com/zenazn/goji/web"

	"app/models"
	"app/utils"
)

var (
	ErrBadAuth  = errors.New("wrong email and/or password")
	ErrBadToken = errors.New("invalid token")
)

func (e *Env) PostSession(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	var data struct {
		Email    string
		Password string
		Remember bool
		Token    string
	}
	err := Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	var user *models.User
	if data.Token == "" {
		user, err = e.M.GetUserByEmail(data.Email)
		if err == utils.ErrNotFound {
			return &Error{E: ErrBadAuth, C: http.StatusBadRequest}
		} else if err != nil {
			return &Error{E: err}
		}

		// TODO: automatically upgrade hash, if work factors differ from config
		err = scrypt.CompareHashAndPassword(
			user.Password,
			[]byte(data.Password),
		)
		if err == scrypt.ErrMismatchedHashAndPassword {
			return &Error{E: ErrBadAuth, C: http.StatusBadRequest}
		} else if err != nil {
			return &Error{E: err}
		}
	} else {
		err = e.M.Atomic(func(etx *models.Env) error {
			otp, inerr := etx.DeleteOTPByToken(data.Token)
			if inerr == utils.ErrNotFound {
				return &Error{E: ErrBadToken, C: http.StatusBadRequest}
			} else if inerr != nil {
				return &Error{E: inerr}
			}

			user, inerr = e.M.GetUserById(otp.UserId)
			if err != nil {
				return &Error{E: inerr, C: http.StatusInternalServerError}
			}

			user.IsEmailVerified = true
			inerr = etx.UpdateUser(user, user.Id)
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
	}

	// User's email verification status isn't checked because there's no
	// correspondence yet, and if she forgets her password she can always reach
	// out to us.

	session, err := e.M.CreateSession(user.Id, data.Remember, r.RemoteAddr)
	if err != nil {
		return &Error{E: err}
	}

	return Created(session, c, w)
}

func (e *Env) DeleteSession(
	c web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	session, err := e.M.GetSessionByToken(c.URLParams["token"])
	if err != nil {
		return &Error{E: err}
	}

	// There's no need in being authenticated when it comes to removing a session,
	// because if I have a valid session key, that alone authorizes me to do
	// anything I want with it.

	err = e.M.DeleteSession(session)
	if err != nil {
		return &Error{E: err}
	}

	return OK(session, c, w)
}
