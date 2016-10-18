package api

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

type handlerFunc func(
	c web.C,
	w http.ResponseWriter,
	r *http.Request,
) *Error

type handlerEnv struct {
	e *Env
	h handlerFunc
}

func (e *handlerEnv) handle(
	c web.C,
	w http.ResponseWriter,
	r *http.Request,
) {
	e.h(c, w, r).Handle(e.e, c, w, r)
}

func (e *Env) NewHandler(h handlerFunc) web.HandlerType {
	return (&handlerEnv{e, h}).handle
}
