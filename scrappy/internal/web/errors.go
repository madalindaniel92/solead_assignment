package web

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrInvalidURL = errors.New("invalid URL")
)
