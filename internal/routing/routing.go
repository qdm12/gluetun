// Package routing defines interfaces to interact with the ip routes using NETLINK.
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
	LocalNetworks() (localNetworks []LocalNetwork, err error)
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

// NewRouting creates a new routing instance.
func NewRouting(logger logging.Logger) Routing {
	return &routing{
		logger:  logger,
		verbose: true,
	}
}

func (r *routing) SetVerbose(verbose bool) {
	r.verbose = verbose
}

func (r *routing) SetDebug() {
	r.debug = true
}
