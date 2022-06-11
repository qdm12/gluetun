package healthcheck

import (
	"context"
	"net"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
)

type Server struct {
	logger  Logger
	handler *handler
	dialer  *net.Dialer
	config  settings.Health
	vpn     vpnHealth
}

func NewServer(config settings.Health,
	logger Logger, vpnLoop StatusApplier) *Server {
	return &Server{
		logger:  logger,
		handler: newHandler(),
		dialer:  &net.Dialer{},
		config:  config,
		vpn: vpnHealth{
			loop:        vpnLoop,
			healthyWait: *config.VPN.Initial,
		},
	}
}

type StatusApplier interface {
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
}
