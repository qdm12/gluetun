package server

import (
	"net/http"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/dns"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/updater"
	"github.com/qdm12/golibs/logging"
)

func newHandlerV0(logger logging.Logger,
	openvpn openvpn.Looper, dns dns.Looper, updater updater.Looper) http.Handler {
	return &handlerV0{
		logger:  logger,
		openvpn: openvpn,
		dns:     dns,
		updater: updater,
	}
}

type handlerV0 struct {
	logger  logging.Logger
	openvpn openvpn.Looper
	dns     dns.Looper
	updater updater.Looper
}

func (h *handlerV0) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "unversioned API: only supports GET method", http.StatusNotFound)
		return
	}
	switch r.RequestURI {
	case "/version":
		http.Redirect(w, r, "/v1/version", http.StatusPermanentRedirect)
	case "/openvpn/actions/restart":
		message := h.openvpn.SetStatus(constants.Stopped)
		h.logger.Info("openvpn: %s", message)
		message = h.openvpn.SetStatus(constants.Running)
		h.logger.Info("openvpn: %s", message)
		if _, err := w.Write([]byte("openvpn restarted, please consider using the /v1/ API in the future.")); err != nil {
			h.logger.Warn(err)
		}
	case "/unbound/actions/restart":
		// message := h.dns.SetStatus(constants.Stopped)
		// h.logger.Info("dns: %s", message)
		// message = h.dns.SetStatus(constants.Running)
		// h.logger.Info("dns: %s", message)
		// if _, err := w.Write([]byte("dns restarted, please consider using the /v1/ API in the future.")); err != nil {
		// 	h.logger.Warn(err)
		// }
	case "/openvpn/portforwarded":
		http.Redirect(w, r, "/v1/openvpn/portforwarded", http.StatusPermanentRedirect)
	case "/openvpn/settings":
		http.Redirect(w, r, "/v1/openvpn/settings", http.StatusPermanentRedirect)
	case "/updater/restart":
		// message := h.updater.SetStatus(constants.Stopped)
		// h.logger.Info("updater: %s", message)
		// message = h.updater.SetStatus(constants.Running)
		// h.logger.Info("updater: %s", message)
		// if _, err := w.Write([]byte("updater restarted, please consider using the /v1/ API in the future.")); err != nil {
		// 	h.logger.Warn(err)
		// }
	default:
		http.Error(w, "unversioned API: requested URI not found", http.StatusNotFound)
	}
}
