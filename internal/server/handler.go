package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/server/middlewares/log"
)

func newHandler(ctx context.Context, logger infoWarner, logging bool,
	buildInfo models.BuildInformation,
	vpnLooper VPNLooper,
	pfGetter PortForwardedGetter,
	dnsLooper DNSLoop,
	updaterLooper UpdaterLooper,
	publicIPLooper PublicIPLoop,
	storage Storage,
	ipv6Supported bool,
) (httpHandler http.Handler) {
	handler := &handler{}

	vpn := newVPNHandler(ctx, vpnLooper, storage, ipv6Supported, logger)
	openvpn := newOpenvpnHandler(ctx, vpnLooper, pfGetter, logger)
	dns := newDNSHandler(ctx, dnsLooper, logger)
	updater := newUpdaterHandler(ctx, updaterLooper, logger)
	publicip := newPublicIPHandler(publicIPLooper, logger)

	handler.v0 = newHandlerV0(ctx, logger, vpnLooper, dnsLooper, updaterLooper)
	handler.v1 = newHandlerV1(logger, buildInfo, vpn, openvpn, dns, updater, publicip)

	middlewares := []func(http.Handler) http.Handler{
		log.New(logger, logging),
	}
	httpHandler = handler
	for _, middleware := range middlewares {
		httpHandler = middleware(httpHandler)
	}
	return httpHandler
}

type handler struct {
	v0 http.Handler
	v1 http.Handler
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
