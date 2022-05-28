// Package server defines an interface to run the HTTP control server.
package server

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/dns"
	"github.com/qdm12/gluetun/internal/httpserver"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/portforward"
	"github.com/qdm12/gluetun/internal/publicip"
	"github.com/qdm12/gluetun/internal/vpn"
)

func New(ctx context.Context, address string, logEnabled bool, logger Logger,
	buildInfo models.BuildInformation, openvpnLooper vpn.Looper,
	pfGetter portforward.Getter, unboundLooper dns.Looper,
	updaterLooper UpdaterLooper, publicIPLooper publicip.Looper) (server httpserver.Runner, err error) {
	handler := newHandler(ctx, logger, logEnabled, buildInfo,
		openvpnLooper, pfGetter, unboundLooper, updaterLooper, publicIPLooper)

	httpServerSettings := httpserver.Settings{
		Address: address,
		Handler: handler,
		Logger:  logger,
	}

	server, err = httpserver.New(httpServerSettings)
	if err != nil {
		return nil, fmt.Errorf("cannot create server: %w", err)
	}

	return server, nil
}
