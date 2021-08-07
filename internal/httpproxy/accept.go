package httpproxy

import (
	"fmt"
	"net/http"
)

func (h *handler) isAccepted(responseWriter http.ResponseWriter, request *http.Request) bool {
	// Not compatible with HTTP < 1.0 or HTTP >= 2.0 (see https://github.com/golang/go/issues/14797#issuecomment-196103814)
	const (
		minimalMajorVersion = 1
		minimalMinorVersion = 0
		maximumMajorVersion = 2
		maximumMinorVersion = 0
	)
	if !request.ProtoAtLeast(minimalMajorVersion, minimalMinorVersion) ||
		request.ProtoAtLeast(maximumMajorVersion, maximumMinorVersion) {
		message := fmt.Sprintf("http version not supported: %s", request.Proto)
		h.logger.Info(message + ", from " + request.RemoteAddr)
		http.Error(responseWriter, message, http.StatusBadRequest)
		return false
	}
	return true
}
