package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/logging"
)

func newOpenvpnHandler(looper openvpn.Looper, logger logging.Logger) http.Handler {
	return &openvpnHandler{
		looper: looper,
		logger: logger,
	}
}

type openvpnHandler struct {
	looper openvpn.Looper
	logger logging.Logger
}

func (h *openvpnHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = strings.TrimPrefix(r.RequestURI, "/openvpn")
	switch r.RequestURI {
	case "/status":
		switch r.Method {
		case http.MethodGet:
			h.getStatus(w)
		case http.MethodPut:
			h.setStatus(w, r)
		default:
			http.Error(w, "", http.StatusNotFound)
		}
	case "/settings":
		switch r.Method {
		case http.MethodGet:
			h.getSettings(w)
		case http.MethodPut:
			h.setSettings(w, r)
		default:
			http.Error(w, "", http.StatusNotFound)
		}
	case "/servers":
		switch r.Method {
		case http.MethodGet:
			h.getServers(w)
		default:
			http.Error(w, "", http.StatusNotFound)
		}
	case "/portforwarded":
		switch r.Method {
		case http.MethodGet:
			h.getPortForwarded(w)
		default:
			http.Error(w, "", http.StatusNotFound)
		}
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}

func (h *openvpnHandler) getStatus(w http.ResponseWriter) {
	status := h.looper.GetStatus()
	encoder := json.NewEncoder(w)
	data := statusWrapper{Status: string(status)}
	if err := encoder.Encode(data); err != nil {
		h.logger.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *openvpnHandler) setStatus(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	decoder := json.NewDecoder(r.Body)
	var data statusWrapper
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	status, err := data.getStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	outcome, err := h.looper.SetStatus(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(outcomeWrapper{Outcome: outcome}); err != nil {
		h.logger.Warn(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *openvpnHandler) getSettings(w http.ResponseWriter) {
	settings := h.looper.GetSettings()
	settings.User = "redacted"
	settings.Password = "redacted"
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(settings); err != nil {
		h.logger.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *openvpnHandler) setSettings(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var settings settings.OpenVPN
	if err := decoder.Decode(&settings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.looper.SetSettings(settings)
	w.WriteHeader(http.StatusOK)
}

func (h *openvpnHandler) getServers(w http.ResponseWriter) {
	servers := h.looper.GetServers()
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(servers); err != nil {
		h.logger.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *openvpnHandler) getPortForwarded(w http.ResponseWriter) {
	port := h.looper.GetPortForwarded()
	encoder := json.NewEncoder(w)
	data := portWrapper{Port: port}
	if err := encoder.Encode(data); err != nil {
		h.logger.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
