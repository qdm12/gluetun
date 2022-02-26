package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/portforward"
	"github.com/qdm12/gluetun/internal/vpn"
)

func newOpenvpnHandler(ctx context.Context, looper vpn.Looper,
	pfGetter portforward.Getter, w warner) http.Handler {
	return &openvpnHandler{
		ctx:    ctx,
		looper: looper,
		pf:     pfGetter,
		warner: w,
	}
}

type openvpnHandler struct {
	ctx    context.Context //nolint:containedctx
	looper vpn.Looper
	pf     portforward.Getter
	warner warner
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
		h.warner.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *openvpnHandler) setStatus(w http.ResponseWriter, r *http.Request) {
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
	outcome, err := h.looper.ApplyStatus(h.ctx, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(outcomeWrapper{Outcome: outcome}); err != nil {
		h.warner.Warn(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *openvpnHandler) getSettings(w http.ResponseWriter) {
	vpnSettings := h.looper.GetSettings()
	settings := vpnSettings.OpenVPN
	settings.User = "redacted"
	settings.Password = "redacted"
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(settings); err != nil {
		h.warner.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *openvpnHandler) getPortForwarded(w http.ResponseWriter) {
	port := h.pf.GetPortForwarded()
	encoder := json.NewEncoder(w)
	data := portWrapper{Port: port}
	if err := encoder.Encode(data); err != nil {
		h.warner.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
