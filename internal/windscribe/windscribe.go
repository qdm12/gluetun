package windscribe

import (
	"net"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// Configurator contains methods to download, read and modify the openvpn configuration to connect as a client
type Configurator interface {
	GetOpenVPNConnections(region models.WindscribeRegion, protocol models.NetworkProtocol, customPort uint16, targetIP net.IP) (connections []models.OpenVPNConnection, err error)
	BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string) (err error)
}

type configurator struct {
	fileManager files.FileManager
	lookupIP    func(host string) ([]net.IP, error)
}

// NewConfigurator returns a new Configurator object
func NewConfigurator(fileManager files.FileManager) Configurator {
	return &configurator{fileManager, net.LookupIP}
}
