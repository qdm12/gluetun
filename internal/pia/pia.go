package pia

import (
	"net"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/verification"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// Configurator contains methods to download, read and modify the openvpn configuration to connect as a client
type Configurator interface {
	MakeOvpn(encryption constants.PIAEncryption,
		protocol constants.NetworkProtocol, region constants.PIARegion) (err error)
}

type configurator struct {
	client      network.Client
	verifier    verification.Verifier
	fileManager files.FileManager
	lookupIP    func(host string) ([]net.IP, error)
}

// NewConfigurator returns a new Configurator object
func NewConfigurator(client network.Client, fileManager files.FileManager) Configurator {
	return &configurator{client, verification.NewVerifier(), fileManager, net.LookupIP}
}

func (c *configurator) MakeOvpn(encryption constants.PIAEncryption,
	protocol constants.NetworkProtocol, region constants.PIARegion) (err error) {
	lines, err := downloadOvpnConfig(c.client, encryption, protocol, region)
	if err != nil {
		return err
	}
	IPs, port, _, err := parseConfig(lines, c.verifier, c.lookupIP)
	if err != nil {
		return err
	}
	lines, err = modifyLines(lines, IPs, port)
	if err != nil {
		return err
	}
	return c.fileManager.WriteLinesToFile(constants.OpenVPNConf, lines)
}
