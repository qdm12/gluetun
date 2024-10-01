package api

import "errors"

var (
	ErrTokenNotValid   = errors.New("token is not valid")
	ErrTooManyRequests = errors.New("too many requests sent for this month")
	ErrBadHTTPStatus   = errors.New("bad HTTP status received")
	ErrServiceLimited  = errors.New("service is limited")
)
