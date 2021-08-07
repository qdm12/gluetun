package healthcheck

import (
	"context"
	"net"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/golibs/logging"
)

var _ ServerRunner = (*Server)(nil)

type ServerRunner interface {
	Run(ctx context.Context, done chan<- struct{})
}

type Server struct {
	logger   logging.Logger
	handler  *handler
	resolver *net.Resolver
	config   configuration.Health
	openvpn  openvpnHealth
}

func NewServer(config configuration.Health,
	logger logging.Logger, openvpnLooper openvpn.Looper) *Server {
	return &Server{
		logger:   logger,
		handler:  newHandler(logger),
		resolver: net.DefaultResolver,
		config:   config,
		openvpn: openvpnHealth{
			looper:      openvpnLooper,
			healthyWait: config.OpenVPN.Initial,
		},
	}
}
