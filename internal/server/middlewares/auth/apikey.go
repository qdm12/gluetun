package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

type apiKeyMethod struct {
	apiKeyDigest [32]byte
}

func newAPIKeyMethod(apiKey string) *apiKeyMethod {
	return &apiKeyMethod{
		apiKeyDigest: sha256.Sum256([]byte(apiKey)),
	}
}

// equal returns true if another auth checker is equal.
// This is used to deduplicate checkers for a particular route.
func (a *apiKeyMethod) equal(other authorizationChecker) bool {
	otherTokenMethod, ok := other.(*apiKeyMethod)
	if !ok {
		return false
	}
	return a.apiKeyDigest == otherTokenMethod.apiKeyDigest
}

func (a *apiKeyMethod) isAuthorized(_ http.Header, request *http.Request) bool {
	xAPIKey := request.Header.Get("X-API-Key")
	if xAPIKey == "" {
		xAPIKey = request.URL.Query().Get("api_key")
	}
	xAPIKeyDigest := sha256.Sum256([]byte(xAPIKey))
	return subtle.ConstantTimeCompare(xAPIKeyDigest[:], a.apiKeyDigest[:]) == 1
}
