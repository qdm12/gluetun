package routing

import (
	"net"

	"github.com/qdm12/golibs/logging"
)

type Routing interface {
	AddRouteVia(destination net.IPNet, gateway net.IP, iface string, table int) error
	DeleteRouteVia(destination net.IPNet) (err error)
	DefaultRoute() (defaultInterface string, defaultGateway net.IP, err error)
	DefaultIP() (ip net.IP, err error)
	LocalSubnet() (defaultSubnet net.IPNet, err error)
	AssignedIP(interfaceName string) (ip net.IP, err error)
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
