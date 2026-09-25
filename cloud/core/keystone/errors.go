package keystone

import "errors"

var (
	ErrUnauthorized = errors.New("keystone: unauthorized")
	ErrForbidden    = errors.New("keystone: forbidden")
	ErrNotFound     = errors.New("keystone: not found")
	ErrBadRequest   = errors.New("keystone: bad request")
	ErrServerError  = errors.New("keystone: server error")
)
