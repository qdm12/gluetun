package auth

import "net/http"

type authorizationChecker interface {
	equal(other authorizationChecker) bool
	isAuthorized(headers http.Header, request *http.Request) bool
}
