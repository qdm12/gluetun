// Package firewall defines a configurator used to change the state
// of the firewall as well as do some light routing changes.
package firewall

import (
	"context"
	"net"
	"sync"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
)

// Configurator allows to change firewall rules and modify network routes.
type Configurator interface {
	SetEnabled(ctx context.Context, enabled bool) (err error)
	SetVPNConnection(ctx context.Context, connection models.OpenVPNConnection) (err error)
	SetAllowedPort(ctx context.Context, port uint16, intf string) (err error)
	SetOutboundSubnets(ctx context.Context, subnets []net.IPNet) (err error)
	RemoveAllowedPort(ctx context.Context, port uint16) (err error)
	// SetNetworkInformation is meant to be called only once
	SetNetworkInformation(defaultInterface string, defaultGateway net.IP,
		localNetworks []routing.LocalNetwork, localIP net.IP)
}

type configurator struct { //nolint:maligned
	commander        command.Commander
	logger           logging.Logger
	routing          routing.Routing
	iptablesMutex    sync.Mutex
	ip6tablesMutex   sync.Mutex
	defaultInterface string
	defaultGateway   net.IP
	localNetworks    []routing.LocalNetwork
	localIP          net.IP
	networkInfoMutex sync.Mutex

	// Fixed state
	ip6Tables       bool
	customRulesPath string

	// State
	enabled           bool
	vpnConnection     models.OpenVPNConnection
	outboundSubnets   []net.IPNet
	allowedInputPorts map[uint16]string // port to interface mapping
	stateMutex        sync.Mutex
}

// NewConfigurator creates a new Configurator instance.
func NewConfigurator(logger logging.Logger, routing routing.Routing) Configurator {
	commander := command.NewCommander()
	return &configurator{
		commander:         commander,
		logger:            logger,
		routing:           routing,
		allowedInputPorts: make(map[uint16]string),
		ip6Tables:         ip6tablesSupported(context.Background(), commander),
		customRulesPath:   "/iptables/post-rules.txt",
	}
}

func (c *configurator) SetNetworkInformation(
	defaultInterface string, defaultGateway net.IP, localNetworks []routing.LocalNetwork, localIP net.IP) {
	c.networkInfoMutex.Lock()
	defer c.networkInfoMutex.Unlock()
	c.defaultInterface = defaultInterface
	c.defaultGateway = defaultGateway
	c.localNetworks = localNetworks
	c.localIP = localIP
}
