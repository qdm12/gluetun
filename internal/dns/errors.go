package dns

import "errors"

var (
	ErrBadStatusCode  = errors.New("bad HTTP status")
	ErrCannotReadBody = errors.New("cannot read response body")
)
