package routing

import (
	"errors"
)

var (
	ErrLinkByIndex         = errors.New("cannot obtain link by index")
	ErrLinkDefaultNotFound = errors.New("default link not found")
	ErrRoutesList          = errors.New("cannot list routes")
)
