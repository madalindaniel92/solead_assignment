package es

import "errors"

var (
	ErrMissingCACertificate = errors.New("missing CA certificate")
	ErrSearchResult         = errors.New("error on search")
	ErrUnexpectedResponse   = errors.New("unexpected response from ES")
)
