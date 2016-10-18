package api

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

type middlewareFunc func(
	c *web.C,
	w http.ResponseWriter,
	r *http.Request,
) *Error

type middlewareEnv struct {
	e  *Env
	mw middlewareFunc
}

func (e *middlewareEnv) handle(c *web.C, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := e.mw(c, w, r)
		if err != nil {
			err.Handle(e.e, *c, w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func (e *Env) NewMiddleware(mw middlewareFunc) web.MiddlewareType {
	return (&middlewareEnv{e, mw}).handle
}
