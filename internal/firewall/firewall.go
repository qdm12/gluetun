package firewall

import (
	"net"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// Configurator allows to change firewall rules and modify network routes
type Configurator interface {
	Clear() error
	BlockAll() error
	CreateGeneralRules() error
	CreateVPNRules(dev models.VPNDevice, serverIPs []net.IP, defaultInterface string,
		port uint16, protocol models.NetworkProtocol) error
	CreateLocalSubnetsRules(subnet net.IPNet, extraSubnets []net.IPNet, defaultInterface string) error
	AddRoutesVia(subnets []net.IPNet, defaultGateway net.IP, defaultInterface string) error
	GetDefaultRoute() (defaultInterface string, defaultGateway net.IP, defaultSubnet net.IPNet, err error)
}

type configurator struct {
	commander   command.Commander
	logger      logging.Logger
	fileManager files.FileManager
}

// NewConfigurator creates a new Configurator instance
func NewConfigurator(logger logging.Logger, fileManager files.FileManager) Configurator {
	return &configurator{
		commander:   command.NewCommander(),
		logger:      logger,
		fileManager: fileManager,
	}
}
