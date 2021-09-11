package healthcheck

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/vpn"
	"github.com/qdm12/golibs/logging"
)

var _ ServerRunner = (*Server)(nil)

type ServerRunner interface {
	Run(ctx context.Context, done chan<- struct{})
}

type Server struct {
	logger  logging.Logger
	handler *handler
	pinger  Pinger
	config  configuration.Health
	vpn     vpnHealth
}

func NewServer(config configuration.Health,
	logger logging.Logger, vpnLooper vpn.Looper) *Server {
	return &Server{
		logger:  logger,
		handler: newHandler(logger),
		pinger:  newPinger(config.AddressToPing),
		config:  config,
		vpn: vpnHealth{
			looper:      vpnLooper,
			healthyWait: config.VPN.Initial,
		},
	}
}
