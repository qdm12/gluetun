package publicip

import "errors"

var (
	ErrBadStatusCode = errors.New("bad HTTP status")
)
