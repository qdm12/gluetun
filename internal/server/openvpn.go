package server

import (
	"encoding/json"
	"net/http"
)

func (s *server) handleGetPortForwarded(w http.ResponseWriter) {
	port := s.getPortForwarded()
	data, err := json.Marshal(struct {
		Port uint16 `json:"port"`
	}{port})
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

