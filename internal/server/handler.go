package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

func newHandler(ctx context.Context, logger infoWarner, logging bool,
	buildInfo models.BuildInformation,
	vpnLooper VPNLooper,
	pf PortForwarding,
	unboundLooper DNSLoop,
	updaterLooper UpdaterLooper,
	publicIPLooper PublicIPLoop,
	storage Storage,
	ipv6Supported bool,
) http.Handler {
	handler := &handler{}

	vpn := newVPNHandler(ctx, vpnLooper, storage, ipv6Supported, logger)
	openvpn := newOpenvpnHandler(ctx, vpnLooper, pf, logger)
	dns := newDNSHandler(ctx, unboundLooper, logger)
	updater := newUpdaterHandler(ctx, updaterLooper, logger)
	publicip := newPublicIPHandler(publicIPLooper, logger)

	handler.v0 = newHandlerV0(ctx, logger, vpnLooper, unboundLooper, updaterLooper)
	handler.v1 = newHandlerV1(logger, buildInfo, vpn, openvpn, dns, updater, publicip)

	handlerWithLog := withLogMiddleware(handler, logger, logging)
	handler.setLogEnabled = handlerWithLog.setEnabled

	return handlerWithLog
}

type handler struct {
	v0            http.Handler
	v1            http.Handler
	setLogEnabled func(enabled bool)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = strings.TrimSuffix(r.RequestURI, "/")
	if !strings.HasPrefix(r.RequestURI, "/v1/") && r.RequestURI != "/v1" {
		h.v0.ServeHTTP(w, r)
		return
	}
	r.RequestURI = strings.TrimPrefix(r.RequestURI, "/v1")
	h.v1.ServeHTTP(w, r)
}
