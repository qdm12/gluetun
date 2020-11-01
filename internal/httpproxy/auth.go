package httpproxy

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func isAuthorized(responseWriter http.ResponseWriter, request *http.Request,
	username, password string) (authorized bool) {
	basicAuth := request.Header.Get("Proxy-Authorization")
	if len(basicAuth) == 0 {
		responseWriter.WriteHeader(http.StatusProxyAuthRequired)
		return false
	}
	b64UsernamePassword := strings.TrimPrefix(basicAuth, "Basic ")
	b, err := base64.StdEncoding.DecodeString(b64UsernamePassword)
	if err != nil {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return false
	}
	usernamePassword := strings.Split(string(b), ":")
	const expectedFields = 2
	if len(usernamePassword) != expectedFields {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return false
	}
	if username != usernamePassword[0] && password != usernamePassword[1] {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return false
	}
	return true
}
