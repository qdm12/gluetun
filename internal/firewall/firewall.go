package firewall

import (
	"context"
	"net"
	"sync"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/routing"
)

// Configurator allows to change firewall rules and modify network routes
type Configurator interface {
	Version(ctx context.Context) (string, error)
	SetEnabled(ctx context.Context, enabled bool) (err error)
	SetVPNConnections(ctx context.Context, connections []models.OpenVPNConnection) (err error)
	SetAllowedSubnets(ctx context.Context, subnets []net.IPNet) (err error)
	SetAllowedPort(ctx context.Context, port uint16) error
	RemoveAllowedPort(ctx context.Context, port uint16) (err error)
	SetPortForward(ctx context.Context, port uint16) (err error)
	SetDebug()
}

type configurator struct { //nolint:maligned
	commander     command.Commander
	logger        logging.Logger
	routing       routing.Routing
	fileManager   files.FileManager // for custom iptables rules
	iptablesMutex sync.Mutex
	debug         bool

	// State
	enabled        bool
	vpnConnections []models.OpenVPNConnection
	allowedSubnets []net.IPNet
	allowedPorts   map[uint16]struct{}
	portForwarded  uint16
	stateMutex     sync.Mutex
}

// NewConfigurator creates a new Configurator instance
func NewConfigurator(logger logging.Logger, routing routing.Routing, fileManager files.FileManager) Configurator {
	return &configurator{
		commander:    command.NewCommander(),
		logger:       logger.WithPrefix("firewall: "),
		routing:      routing,
		fileManager:  fileManager,
		allowedPorts: make(map[uint16]struct{}),
	}
}

func (c *configurator) SetDebug() {
	c.debug = true
}
