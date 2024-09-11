package auth

import (
	"crypto/subtle"
	"net/http"
)

type apiKeyMethod struct {
	apiKey string
}

func newAPIKeyMethod(apiKey string) *apiKeyMethod {
	return &apiKeyMethod{
		apiKey: apiKey,
	}
}

// equal returns true if another auth checker is equal.
// This is used to deduplicate checkers for a particular route.
func (a *apiKeyMethod) equal(other authorizationChecker) bool {
	otherTokenMethod, ok := other.(*apiKeyMethod)
	if !ok {
		return false
	}
	return a.apiKey == otherTokenMethod.apiKey
}

func (a *apiKeyMethod) isAuthorized(_ http.Header, request *http.Request) bool {
	xAPIKey := request.Header.Get("X-API-Key")
	if xAPIKey == "" {
		xAPIKey = request.URL.Query().Get("api_key")
	}
	return subtle.ConstantTimeCompare([]byte(xAPIKey), []byte(a.apiKey)) == 1
}
