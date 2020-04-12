package firewall

import (
	"net"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

const logPrefix = "firewall configurator"

// Configurator allows to change firewall rules and modify network routes
type Configurator interface {
	Version() (string, error)
	AcceptAll() error
	Clear() error
	BlockAll() error
	CreateGeneralRules() error
	CreateVPNRules(dev models.VPNDevice, defaultInterface string, connections []models.OpenVPNConnection) error
	CreateLocalSubnetsRules(subnet net.IPNet, extraSubnets []net.IPNet, defaultInterface string) error
	AllowInputTrafficOnPort(device models.VPNDevice, port uint16) error
	AllowAnyIncomingOnPort(port uint16) error
}

type configurator struct {
	commander command.Commander
	logger    logging.Logger
}

// NewConfigurator creates a new Configurator instance
func NewConfigurator(logger logging.Logger) Configurator {
	return &configurator{
		commander: command.NewCommander(),
		logger:    logger,
	}
}
