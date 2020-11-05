package server

import (
	"encoding/json"
	"net/http"
)

func (h *handler) getPortForwarded(w http.ResponseWriter) {
	port := h.openvpnLooper.GetPortForwarded()
	data, err := json.Marshal(struct {
		Port uint16 `json:"port"`
	}{port})
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

func (h *handler) getOpenvpnSettings(w http.ResponseWriter) {
	settings := h.openvpnLooper.GetSettings()
	data, err := json.Marshal(settings)
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
