package server

import (
	"encoding/json"
	"net/http"
)

func (h *handler) getVersion(w http.ResponseWriter) {
	data, err := json.Marshal(h.buildInfo)
	if err != nil {
		h.logger.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(data); err != nil {
		h.logger.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
