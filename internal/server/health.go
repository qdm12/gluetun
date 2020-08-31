package server

import (
	"fmt"
	"net/http"
)

func (s *server) handleHealth(w http.ResponseWriter) {
	// TODO option to disable
	// TODO use mullvad API if current provider is Mullvad
	ips, err := s.lookupIP("github.com")
	var errorMessage string
	switch {
	case err != nil:
		errorMessage = fmt.Sprintf("cannot resolve github.com (%s)", err)
	case len(ips) == 0:
		errorMessage = "resolved no IP addresses for github.com"
	default: // success
		w.WriteHeader(http.StatusOK)
		return
	}
	s.logger.Warn(errorMessage)
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := w.Write([]byte(errorMessage)); err != nil {
		s.logger.Warn(err)
	}
}
