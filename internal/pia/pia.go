package pia

import (
	"net"

	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/verification"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// Configurator contains methods to download, read and modify the openvpn configuration to connect as a client
type Configurator interface {
	DownloadOvpnConfig(encryption models.PIAEncryption,
		protocol models.NetworkProtocol, region models.PIARegion) (lines []string, err error)
	ParseConfig(lines []string) (IPs []net.IP, port uint16, device models.VPNDevice, err error)
	ModifyLines(lines []string, IPs []net.IP, port uint16) (modifiedLines []string, err error)
	PortForward(filepath models.Filepath)
}

type configurator struct {
	client     network.Client
	verifyPort func(port string) error
	lookupIP   func(host string) ([]net.IP, error)
}

// NewConfigurator returns a new Configurator object
func NewConfigurator(client network.Client) Configurator {
	return &configurator{client, verification.NewVerifier().VerifyPort, net.LookupIP}
}
