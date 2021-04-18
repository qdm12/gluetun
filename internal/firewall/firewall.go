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
	"github.com/qdm12/golibs/os"
)

// Configurator allows to change firewall rules and modify network routes.
type Configurator interface {
	Version(ctx context.Context) (string, error)
	SetEnabled(ctx context.Context, enabled bool) (err error)
	SetVPNConnection(ctx context.Context, connection models.OpenVPNConnection) (err error)
	SetAllowedPort(ctx context.Context, port uint16, intf string) (err error)
	SetOutboundSubnets(ctx context.Context, subnets []net.IPNet) (err error)
	RemoveAllowedPort(ctx context.Context, port uint16) (err error)
	SetDebug()
	// SetNetworkInformation is meant to be called only once
	SetNetworkInformation(defaultInterface string, defaultGateway net.IP,
		localNetworks []routing.LocalNetwork, localIP net.IP)
}

type configurator struct { //nolint:maligned
	commander        command.Commander
	logger           logging.Logger
	routing          routing.Routing
	openFile         os.OpenFileFunc // for custom iptables rules
	iptablesMutex    sync.Mutex
	ip6tablesMutex   sync.Mutex
	debug            bool
	defaultInterface string
	defaultGateway   net.IP
	localNetworks    []routing.LocalNetwork
	localIP          net.IP
	networkInfoMutex sync.Mutex

	// State
	enabled           bool
	vpnConnection     models.OpenVPNConnection
	outboundSubnets   []net.IPNet
	allowedInputPorts map[uint16]string // port to interface mapping
	stateMutex        sync.Mutex
}

// NewConfigurator creates a new Configurator instance.
func NewConfigurator(logger logging.Logger, routing routing.Routing, openFile os.OpenFileFunc) Configurator {
	return &configurator{
		commander:         command.NewCommander(),
		logger:            logger.NewChild(logging.SetPrefix("firewall: ")),
		routing:           routing,
		openFile:          openFile,
		allowedInputPorts: make(map[uint16]string),
	}
}

func (c *configurator) SetDebug() {
	c.debug = true
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
