package scrappy

import "errors"

var (
	ErrInvalidCSVHeader = errors.New("csv file has invalid header")
)
