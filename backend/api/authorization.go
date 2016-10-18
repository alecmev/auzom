package api

import (
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"app/utils"
)

const (
	SessionAgeLong        = time.Hour * 24 * 30 // inspired by Google
	SessionAgeShortWindow = time.Minute * 15
	SessionAgeShortMax    = time.Hour * 24
	SessionAgeSensitive   = time.Minute * 15
)

func (e *Env) Auth(
	c *web.C, w http.ResponseWriter, r *http.Request,
) *Error {
	authorization, ok := r.Header["Authorization"]
	if !ok {
		return nil
	}

	session, err := e.M.GetSessionByToken(authorization[0])
	if err == utils.ErrNotFound {
		return &Error{E: utils.ErrUnauthorized, M: "bad token"}
	} else if err != nil {
		return &Error{E: err}
	}

	now := time.Now()
	isExpired := false
	if session.Remember {
		isExpired = session.CreatedAt.Add(SessionAgeLong).Before(now)
	} else {
		isExpired = session.LastUsedAt.Add(SessionAgeShortWindow).Before(now) ||
			session.CreatedAt.Add(SessionAgeShortMax).Before(now)
	}

	if isExpired {
		err = e.M.DeleteSession(session)
		if err != nil {
			return &Error{E: err}
		}

		return &Error{E: utils.ErrUnauthorized, M: "session expired"}
	}

	session.LastUsedAt = now
	session.LastUsedIp = r.RemoteAddr
	err = e.M.UpdateSessionLastUsed(session)
	if err != nil {
		return &Error{E: err}
	}

	c.Env["session"] = session
	return nil
}
