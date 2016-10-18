package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/elithrar/simple-scrypt"
	"github.com/getsentry/raven-go"
	"github.com/goji/param"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"app/mail"
	"app/models"
	"app/slack"
)

type Env struct {
	StaticHost string
	Scrypt     scrypt.Params
	M          *models.Env
	Mail       *mail.Env
	Slack      *slack.Env
	Sentry     *raven.Client
}

func New(
	staticHost string,
	scryptParams scrypt.Params,
	modelsEnv *models.Env,
	mailEnv *mail.Env,
	slackEnv *slack.Env,
	sentry *raven.Client,
) *Env {
	return &Env{staticHost, scryptParams, modelsEnv, mailEnv, slackEnv, sentry}
}

func Decode(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func DecodeQuery(r *http.Request, v interface{}) error {
	return param.Parse(r.URL.Query(), v)
}

func Respond(
	code int, data interface{}, c web.C, w http.ResponseWriter,
) *Error {
	w.Header().Set("Request-Id", middleware.GetReqID(c))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			return &Error{E: err}
		}
	}

	return nil
}

func OK(data interface{}, c web.C, w http.ResponseWriter) *Error {
	return Respond(http.StatusOK, data, c, w)
}

func Created(data interface{}, c web.C, w http.ResponseWriter) *Error {
	return Respond(http.StatusCreated, data, c, w)
}

func NoContent(c web.C, w http.ResponseWriter) *Error {
	return Respond(http.StatusNoContent, nil, c, w)
}

func NotFound(c web.C, w http.ResponseWriter) *Error {
	return Respond(http.StatusNotFound, nil, c, w)
}

// ParseInt turns a string into an int. Returns 0 if the string has no length.
// Reason for parsing is to avoid hitting the DB with an invalid query.
func ParseInt(s string, name string) (int, error) {
	if len(s) == 0 {
		return 0, nil
	}

	result, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New(name + " is not a valid number")
	}

	return result, nil
}
