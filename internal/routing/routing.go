package routing

import (
	"context"
	"net"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
)

type Routing interface {
	AddRouteVia(ctx context.Context, subnet net.IPNet, defaultGateway net.IP, defaultInterface string) error
	DeleteRouteVia(ctx context.Context, subnet net.IPNet) (err error)
	DefaultRoute() (defaultInterface string, defaultGateway net.IP, err error)
	LocalSubnet() (defaultSubnet net.IPNet, err error)
	VPNDestinationIP() (ip net.IP, err error)
	VPNLocalGatewayIP() (ip net.IP, err error)
	SetDebug()
}

type routing struct {
	commander command.Commander
	logger    logging.Logger
	debug     bool
}

// NewConfigurator creates a new Configurator instance.
func NewRouting(logger logging.Logger) Routing {
	return &routing{
		commander: command.NewCommander(),
		logger:    logger.WithPrefix("routing: "),
	}
}

func (c *routing) SetDebug() {
	c.debug = true
}
