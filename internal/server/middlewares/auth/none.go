package auth

import "net/http"

type noneMethod struct{}

func newNoneMethod() *noneMethod {
	return &noneMethod{}
}

// equal returns true if another auth checker is equal.
// This is used to deduplicate checkers for a particular route.
func (n *noneMethod) equal(other authorizationChecker) bool {
	_, ok := other.(*noneMethod)
	return ok
}

func (n *noneMethod) isAuthorized(_ http.Header, _ *http.Request) bool {
	return true
}
