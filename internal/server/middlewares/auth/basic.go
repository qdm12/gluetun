package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"strings"
)

type basicAuthMethod struct {
	userpasshash [32]byte
}

func newBasicAuthMethod(user string, pass string) *basicAuthMethod {
	var sb strings.Builder

	sb.WriteString(user)
	sb.WriteString(pass)
	var matchbytes = sha256.Sum256([]byte(sb.String()))
	return &basicAuthMethod{
		userpasshash: matchbytes,
	}
}

// equal returns true if another auth checker is equal.
// This is used to deduplicate checkers for a particular route.
func (a *basicAuthMethod) equal(other authorizationChecker) bool {
	otherBasicMethod, ok := other.(*basicAuthMethod)
	if !ok {
		return false
	}
	return a.userpasshash == otherBasicMethod.userpasshash
}

func (a *basicAuthMethod) isAuthorized(r *http.Request) bool {
	authsuccess := false
	var inpSb strings.Builder
	// Get Inputs from http request
	username, password, ok := r.BasicAuth()
	if ok {
		inpSb.WriteString(username)
		inpSb.WriteString(password)
		inputhash := sha256.Sum256([]byte(inpSb.String()))
		authsuccess = (subtle.ConstantTimeCompare(a.userpasshash[:], inputhash[:]) == 1)
	}
	if authsuccess {
		return true
	}

	return false
}
