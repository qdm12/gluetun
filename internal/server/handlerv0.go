package server

import (
	"context"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants"
)

func newHandlerV0(ctx context.Context, logger infoWarner,
	vpn VPNLooper, dns DNSLoop, updater UpdaterLooper) http.Handler {
	return &handlerV0{
		ctx:     ctx,
		logger:  logger,
		vpn:     vpn,
		dns:     dns,
		updater: updater,
	}
}

type handlerV0 struct {
	ctx     context.Context //nolint:containedctx
	logger  infoWarner
	vpn     VPNLooper
	dns     DNSLoop
	updater UpdaterLooper
}

func (h *handlerV0) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "unversioned API: only supports GET method", http.StatusBadRequest)
		return
	}
	switch r.RequestURI {
	case "/version":
		http.Redirect(w, r, "/v1/version", http.StatusPermanentRedirect)
	case "/openvpn/actions/restart":
		outcome, _ := h.vpn.ApplyStatus(h.ctx, constants.Stopped)
		h.logger.Info("openvpn: " + outcome)
		outcome, _ = h.vpn.ApplyStatus(h.ctx, constants.Running)
		h.logger.Info("openvpn: " + outcome)
		if _, err := w.Write([]byte("openvpn restarted, please consider using the /v1/ API in the future.")); err != nil {
			h.logger.Warn(err.Error())
		}
	case "/unbound/actions/restart":
		outcome, _ := h.dns.ApplyStatus(h.ctx, constants.Stopped)
		h.logger.Info("dns: " + outcome)
		outcome, _ = h.dns.ApplyStatus(h.ctx, constants.Running)
		h.logger.Info("dns: " + outcome)
		if _, err := w.Write([]byte("dns restarted, please consider using the /v1/ API in the future.")); err != nil {
			h.logger.Warn(err.Error())
		}
	case "/openvpn/portforwarded":
		http.Redirect(w, r, "/v1/openvpn/portforwarded", http.StatusPermanentRedirect)
	case "/openvpn/settings":
		http.Redirect(w, r, "/v1/openvpn/settings", http.StatusPermanentRedirect)
	case "/updater/restart":
		outcome, _ := h.updater.SetStatus(h.ctx, constants.Stopped)
		h.logger.Info("updater: " + outcome)
		outcome, _ = h.updater.SetStatus(h.ctx, constants.Running)
		h.logger.Info("updater: " + outcome)
		if _, err := w.Write([]byte("updater restarted, please consider using the /v1/ API in the future.")); err != nil {
			h.logger.Warn(err.Error())
		}
	default:
		http.Error(w, "unversioned API: requested URI not found", http.StatusBadRequest)
	}
}
