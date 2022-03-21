package healthcheck

import (
	"context"
	"net"

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
	dialer  *net.Dialer
	config  settings.Health
	vpn     vpnHealth
}

func NewServer(config settings.Health,
	logger Logger, vpnLooper vpn.Looper) *Server {
	return &Server{
		logger:  logger,
		handler: newHandler(),
		dialer:  &net.Dialer{},
		config:  config,
		vpn: vpnHealth{
			looper:      vpnLooper,
			healthyWait: *config.VPN.Initial,
		},
	}
}
