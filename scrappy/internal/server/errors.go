package server

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrMethodNotSupported = errors.New("method not supported")
)
