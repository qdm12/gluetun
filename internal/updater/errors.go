package updater

import "errors"

var (
	ErrHTTPStatusCodeNotOK     = errors.New("HTTP status code not OK")
	ErrUnmarshalResponseBody   = errors.New("cannot unmarshal response body")
	ErrUpdateServerInformation = errors.New("failed updating server information")
)
