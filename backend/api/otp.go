package api

import (
	"net/http"

	"github.com/zenazn/goji/web"

	"app/utils"
)

func (e *Env) PostOTP(c web.C, w http.ResponseWriter, r *http.Request) *Error {
	var data struct {
		Email string
	}
	err := Decode(r, &data)
	if err != nil {
		return &Error{E: err, C: http.StatusBadRequest}
	}

	user, err := e.M.GetUserByEmail(data.Email)
	if err == utils.ErrNotFound {
		return &Error{E: err, M: "no user with such email"}
	} else if err != nil {
		return &Error{E: err}
	} else if !user.IsEmailVerified {
		return &Error{
			C: http.StatusBadRequest,
			M: "this email address isn't verified yet; " +
				"contact support@auzom.gg if you can't verify it for some reason",
		}
	}

	otp, err := e.M.CreateOTP(user.Id)
	if err != nil {
		return &Error{E: err}
	}

	err = e.Mail.Send(
		data.Email,
		"Password reset link",
		"[Reset password](https://"+e.StaticHost+"/password-reset/"+otp.Token+")",
	)
	if err != nil {
		return &Error{E: err}
	}

	return NoContent(c, w)
}
