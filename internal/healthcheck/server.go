package healthcheck

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/vpn"
)

var _ ServerRunner = (*Server)(nil)

type ServerRunner interface {
	Run(ctx context.Context, done chan<- struct{})
}

type Server struct {
	logger  Logger
	handler *handler
	pinger  Pinger
	config  settings.Health
	vpn     vpnHealth
}

func NewServer(config settings.Health,
	logger Logger, vpnLooper vpn.Looper) *Server {
	return &Server{
		logger:  logger,
		handler: newHandler(),
		pinger:  newPinger(config.AddressToPing),
		config:  config,
		vpn: vpnHealth{
			looper:      vpnLooper,
			healthyWait: *config.VPN.Initial,
		},
	}
}
