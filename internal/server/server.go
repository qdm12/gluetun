package server

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/httpserver"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/server/middlewares/auth"
)

func New(ctx context.Context, address string, logEnabled bool, logger Logger,
	authSettings auth.Settings, buildInfo models.BuildInformation, openvpnLooper VPNLooper,
	pfGetter PortForwardedGetter, dnsLooper DNSLoop,
	updaterLooper UpdaterLooper, publicIPLooper PublicIPLoop, storage Storage,
	ipv6Supported bool) (
	server *httpserver.Server, err error) {
	handler, err := newHandler(ctx, logger, logEnabled, authSettings, buildInfo,
		openvpnLooper, pfGetter, dnsLooper, updaterLooper, publicIPLooper,
		storage, ipv6Supported)
	if err != nil {
		return nil, fmt.Errorf("creating handler: %w", err)
	}

	httpServerSettings := httpserver.Settings{
		Address: address,
		Handler: handler,
		Logger:  logger,
	}

	server, err = httpserver.New(httpServerSettings)
	if err != nil {
		return nil, fmt.Errorf("creating server: %w", err)
	}

	return server, nil
}
