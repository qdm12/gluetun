package pia

import (
	"net"

	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/verification"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// Configurator contains methods to download, read and modify the openvpn configuration to connect as a client
type Configurator interface {
	DownloadOvpnConfig(encryption constants.PIAEncryption,
		protocol constants.NetworkProtocol, region constants.PIARegion) (lines []string, err error)
	ParseConfig(lines []string) (IPs []net.IP, port uint16, device models.VPNDevice, err error)
	ModifyLines(lines []string, IPs []net.IP, port uint16) (modifiedLines []string, err error)
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
