package csv

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidCSVHeader = errors.New("csv file has invalid header")
	ErrEmptyCSV         = errors.New("empty CSV")
	ErrMissingURLHost   = errors.New("missing URL host")
)

// We don't want to spam stdout with more than maxInvalidLines in case too many URLs
// from the CSV are invalid, so only this many errors will be shown.
const MaxInvalidCSVLines = 20

type InvalidCSVLine struct {
	Index int
	Line  string
	Err   error
}

type ErrInvalidCSVLines []InvalidCSVLine

func (e ErrInvalidCSVLines) Error() string {
	return fmt.Sprintf("%d invalid CSV lines", len(e))
}

// Allows error type to be compared using errors.Is
func (e ErrInvalidCSVLines) Is(target error) bool {
	_, ok := target.(ErrInvalidCSVLines)
	return ok
}

// Gather invalid lines, but only the first MaxInvalidCSVLines
func (e ErrInvalidCSVLines) Append(err error, line string, index int) ErrInvalidCSVLines {
	if len(e) >= MaxInvalidCSVLines {
		return e
	}

	return append(e, InvalidCSVLine{Err: err, Line: line, Index: index})
}
