package server

import (
	"context"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/dns"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/updater"
	"github.com/qdm12/golibs/logging"
)

func newHandlerV0(ctx context.Context, logger logging.Logger,
	openvpn openvpn.Looper, dns dns.Looper, updater updater.Looper) http.Handler {
	return &handlerV0{
		ctx:     ctx,
		logger:  logger,
		openvpn: openvpn,
		dns:     dns,
		updater: updater,
	}
}

type handlerV0 struct {
	ctx     context.Context
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
		outcome, _ := h.openvpn.SetStatus(h.ctx, constants.Stopped)
		h.logger.Info("openvpn: %s", outcome)
		outcome, _ = h.openvpn.SetStatus(h.ctx, constants.Running)
		h.logger.Info("openvpn: %s", outcome)
		if _, err := w.Write([]byte("openvpn restarted, please consider using the /v1/ API in the future.")); err != nil {
			h.logger.Warn(err)
		}
	case "/unbound/actions/restart":
		outcome, _ := h.dns.ApplyStatus(h.ctx, constants.Stopped)
		h.logger.Info("dns: %s", outcome)
		outcome, _ = h.dns.ApplyStatus(h.ctx, constants.Running)
		h.logger.Info("dns: %s", outcome)
		if _, err := w.Write([]byte("dns restarted, please consider using the /v1/ API in the future.")); err != nil {
			h.logger.Warn(err)
		}
	case "/openvpn/portforwarded":
		http.Redirect(w, r, "/v1/openvpn/portforwarded", http.StatusPermanentRedirect)
	case "/openvpn/settings":
		http.Redirect(w, r, "/v1/openvpn/settings", http.StatusPermanentRedirect)
	case "/updater/restart":
		outcome, _ := h.updater.SetStatus(h.ctx, constants.Stopped)
		h.logger.Info("updater: %s", outcome)
		outcome, _ = h.updater.SetStatus(h.ctx, constants.Running)
		h.logger.Info("updater: %s", outcome)
		if _, err := w.Write([]byte("updater restarted, please consider using the /v1/ API in the future.")); err != nil {
			h.logger.Warn(err)
		}
	default:
		http.Error(w, "unversioned API: requested URI not found", http.StatusNotFound)
	}
}
