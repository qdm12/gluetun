package routing

import (
	"net"
	"sync"

	"github.com/qdm12/golibs/logging"
)

type Routing interface {
	// Mutations
	Setup() (err error)
	TearDown() error
	SetOutboundRoutes(outboundSubnets []net.IPNet) error

	// Read only
	DefaultRoute() (defaultInterface string, defaultGateway net.IP, err error)
	LocalSubnet() (defaultSubnet net.IPNet, err error)
	DefaultIP() (defaultIP net.IP, err error)
	VPNDestinationIP() (ip net.IP, err error)
	VPNLocalGatewayIP() (ip net.IP, err error)

	// Internal state
	SetVerbose(verbose bool)
	SetDebug()
}

type routing struct {
	logger          logging.Logger
	verbose         bool
	debug           bool
	outboundSubnets []net.IPNet
	stateMutex      sync.RWMutex
}

// NewConfigurator creates a new Configurator instance.
func NewRouting(logger logging.Logger) Routing {
	return &routing{
		logger:  logger.WithPrefix("routing: "),
		verbose: true,
	}
}

func (c *routing) SetVerbose(verbose bool) {
	c.verbose = verbose
}

func (c *routing) SetDebug() {
	c.debug = true
}
