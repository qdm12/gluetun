package server

import (
	"encoding/json"
	"net/http"
)

func (s *server) handleGetVersion(w http.ResponseWriter) {
	data, err := json.Marshal(s.buildInfo)
	if err != nil {
		s.logger.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(data); err != nil {
		s.logger.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
