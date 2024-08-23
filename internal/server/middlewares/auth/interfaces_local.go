package auth

import "net/http"

type authorizationChecker interface {
	equal(other authorizationChecker) bool
	isAuthorized(request *http.Request) bool
}
