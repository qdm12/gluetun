package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/dns"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/publicip"
	"github.com/qdm12/gluetun/internal/updater"
	"github.com/qdm12/golibs/logging"
)

func newHandler(ctx context.Context, logger logging.Logger, logging bool,
	buildInfo models.BuildInformation,
	openvpnLooper openvpn.Looper,
	unboundLooper dns.Looper,
	updaterLooper updater.Looper,
	publicIPLooper publicip.Looper,
) http.Handler {
	handler := &handler{}

	openvpn := newOpenvpnHandler(ctx, openvpnLooper, logger)
	dns := newDNSHandler(ctx, unboundLooper, logger)
	updater := newUpdaterHandler(ctx, updaterLooper, logger)
	publicip := newPublicIPHandler(publicIPLooper, logger)

	handler.v0 = newHandlerV0(ctx, logger, openvpnLooper, unboundLooper, updaterLooper)
	handler.v1 = newHandlerV1(logger, buildInfo, openvpn, dns, updater, publicip)

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
