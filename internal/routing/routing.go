// Package routing defines interfaces to interact with the ip routes using NETLINK.
package routing

import (
	"net"
	"sync"

	"github.com/qdm12/golibs/logging"
)

type Routing interface {
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

type routing struct {
	logger          logging.Logger
	outboundSubnets []net.IPNet
	stateMutex      sync.RWMutex
}

// NewRouting creates a new routing instance.
func NewRouting(logger logging.Logger) Routing {
	return &routing{
		logger: logger,
	}
}
