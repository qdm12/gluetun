// Package routing defines interfaces to interact with the ip routes using NETLINK.
package routing

import (
	"net"
	"sync"

	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/qdm12/golibs/logging"
)

type ReadWriter interface {
	Reader
	Writer
}

type Reader interface {
	DefaultRouteGetter
	DefaultIPGetter
	LocalSubnetGetter
	LocalNetworksGetter
	VPNGetter
}

type VPNGetter interface {
	VPNDestinationIPGetter
	VPNLocalGatewayIPGetter
}

type Writer interface {
	Setuper
	TearDowner
	OutboundRoutesSetter
}

type Routing struct {
	netLinker       netlink.NetLinker
	logger          logging.Logger
	outboundSubnets []net.IPNet
	stateMutex      sync.RWMutex
}

// New creates a new routing instance.
func New(netLinker netlink.NetLinker,
	logger logging.Logger) *Routing {
	return &Routing{
		netLinker: netLinker,
		logger:    logger,
	}
}
