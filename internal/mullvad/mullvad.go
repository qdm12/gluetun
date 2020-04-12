package mullvad

import (
	"net"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// Configurator contains methods to download, read and modify the openvpn configuration to connect as a client
type Configurator interface {
	GetOpenVPNConnections(country models.MullvadCountry, city models.MullvadCity, provider models.MullvadProvider, protocol models.NetworkProtocol, customPort uint16, targetIP net.IP) (connections []models.OpenVPNConnection, err error)
	BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher string) (err error)
}

type configurator struct {
	fileManager files.FileManager
	logger      logging.Logger
}

// NewConfigurator returns a new Configurator object
func NewConfigurator(fileManager files.FileManager, logger logging.Logger) Configurator {
	return &configurator{
		fileManager: fileManager,
		logger:      logger.WithPrefix("Mullvad configurator: "),
	}
}
