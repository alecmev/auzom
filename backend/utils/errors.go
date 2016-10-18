package utils

import (
	"errors"
	"net/http"
	"strings"
)

var errCode = map[error]int{}

func ErrCode(err error) (code int, ok bool) {
	code, ok = errCode[err]
	return
}

func newErr(message string, code int) error {
	if message == "" {
		message = strings.ToLower(http.StatusText(code))
	}

	err := errors.New(message)
	errCode[err] = code
	return err
}

var (
	ErrTODO = newErr("TODO", http.StatusInternalServerError)

	ErrBadRequest   = newErr("", http.StatusBadRequest)
	ErrUnauthorized = newErr("unauthenticated", http.StatusUnauthorized)
	ErrForbidden    = newErr("", http.StatusForbidden)
	ErrNotFound     = newErr("", http.StatusNotFound)

	ErrInternal = newErr("", http.StatusInternalServerError)
)
