package routing

import (
	"net"

	"github.com/qdm12/golibs/logging"
)

type Routing interface {
	Setup() (err error)
	TearDown() error
	DefaultRoute() (defaultInterface string, defaultGateway net.IP, err error)
	LocalSubnet() (defaultSubnet net.IPNet, err error)
	VPNDestinationIP() (ip net.IP, err error)
	VPNLocalGatewayIP() (ip net.IP, err error)
	SetDebug()
}

type routing struct {
	logger logging.Logger
	debug  bool
}

// NewConfigurator creates a new Configurator instance.
func NewRouting(logger logging.Logger) Routing {
	return &routing{
		logger: logger.WithPrefix("routing: "),
	}
}

func (c *routing) SetDebug() {
	c.debug = true
}
