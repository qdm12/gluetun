package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

type basicAuthMethod struct {
	authDigest [32]byte
}

func newBasicAuthMethod(username, password string) *basicAuthMethod {
	return &basicAuthMethod{
		authDigest: sha256.Sum256([]byte(username + password)),
	}
}

// equal returns true if another auth checker is equal.
// This is used to deduplicate checkers for a particular route.
func (a *basicAuthMethod) equal(other authorizationChecker) bool {
	otherBasicMethod, ok := other.(*basicAuthMethod)
	if !ok {
		return false
	}
	return a.authDigest == otherBasicMethod.authDigest
}

func (a *basicAuthMethod) isAuthorized(headers http.Header, request *http.Request) bool {
	username, password, ok := request.BasicAuth()
	if !ok {
		headers.Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		return false
	}
	requestAuthDigest := sha256.Sum256([]byte(username + password))
	return subtle.ConstantTimeCompare(a.authDigest[:], requestAuthDigest[:]) == 1
}
