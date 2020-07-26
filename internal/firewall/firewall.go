package firewall

import (
	"context"
	"net"
	"sync"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

// Configurator allows to change firewall rules and modify network routes
type Configurator interface {
	Version(ctx context.Context) (string, error)
	SetEnabled(ctx context.Context, enabled bool) (err error)
	SetVPNConnections(ctx context.Context, connections []models.OpenVPNConnection) (err error)
	SetAllowedSubnets(ctx context.Context, subnets []net.IPNet) (err error)
	SetAllowedPort(ctx context.Context, port uint16, intf string) (err error)
	RemoveAllowedPort(ctx context.Context, port uint16) (err error)
	SetDebug()
	// SetNetworkInformation is meant to be called only once
	SetNetworkInformation(defaultInterface string, defaultGateway net.IP, localSubnet net.IPNet)
}

type configurator struct { //nolint:maligned
	commander        command.Commander
	logger           logging.Logger
	routing          routing.Routing
	fileManager      files.FileManager // for custom iptables rules
	iptablesMutex    sync.Mutex
	debug            bool
	defaultInterface string
	defaultGateway   net.IP
	localSubnet      net.IPNet
	networkInfoMutex sync.Mutex

	// State
	enabled           bool
	vpnConnections    []models.OpenVPNConnection
	allowedSubnets    []net.IPNet
	allowedInputPorts map[uint16]string // port to interface mapping
	stateMutex        sync.Mutex
}

// NewConfigurator creates a new Configurator instance
func NewConfigurator(logger logging.Logger, routing routing.Routing, fileManager files.FileManager) Configurator {
	return &configurator{
		commander:         command.NewCommander(),
		logger:            logger.WithPrefix("firewall: "),
		routing:           routing,
		fileManager:       fileManager,
		allowedInputPorts: make(map[uint16]string),
	}
}

func (c *configurator) SetDebug() {
	c.debug = true
}

func (c *configurator) SetNetworkInformation(defaultInterface string, defaultGateway net.IP, localSubnet net.IPNet) {
	c.networkInfoMutex.Lock()
	defer c.networkInfoMutex.Unlock()
	c.defaultInterface = defaultInterface
	c.defaultGateway = defaultGateway
	c.localSubnet = localSubnet
}
