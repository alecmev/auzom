package api

import (
	"log"
	"net/http"

	"github.com/getsentry/raven-go"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"app/utils"
)

type Error struct {
	E error       `json:"-"`
	M string      `json:"message"`
	D interface{} `json:"-"`
	C int         `json:"-"`
}

func (err *Error) Error() string {
	if err.M != "" {
		return err.M
	}
	if err.E != nil {
		return err.E.Error()
	}
	if err.C != 0 {
		return http.StatusText(err.C)
	}

	return ""
}

func (err *Error) Handle(
	e *Env,
	c web.C,
	w http.ResponseWriter,
	r *http.Request,
) {
	if err == nil {
		return
	}

	if err.M == "" {
		if err.E != nil {
			err.M = err.E.Error()
		} else {
			err.M = http.StatusText(err.C)
		}
	}

	if err.C == 0 {
		if code, ok := utils.ErrCode(err.E); ok {
			err.C = code
		} else {
			err.C = http.StatusInternalServerError
		}
	}

	packet := raven.NewPacket(err.M, raven.NewHttp(r))
	packet.EventID = middleware.GetReqID(c)
	packet.Extra["data"] = err.D
	packet.Extra["code"] = err.C
	if err.C < http.StatusInternalServerError {
		Respond(err.C, err, c, w)
		packet.Level = raven.WARNING
	} else {
		http.Error(w, http.StatusText(err.C), err.C)
		packet.Level = raven.ERROR
	}

	log.Printf("err.E: %v\n", err.E)
	log.Printf("err.M: %v\n", err.M)
	log.Printf("err.D: %v\n", err.D)

	e.Sentry.Capture(packet, nil)
}

// TODO: implement the "affected user" feature of Sentry

// TODO: consider removing the special treatment of errors defined in utils
// (that is always requiring the error code, and defaulting to 500 if none
// supplied), because it introduces quite a bit of magic into error handling
