package pia

import (
	"net"

	"github.com/qdm12/golibs/crypto/random"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/verification"
	"github.com/qdm12/private-internet-access-docker/internal/firewall"
)

type pia struct {
	client      network.Client
	fileManager files.FileManager
	firewall    firewall.Configurator
	random      random.Random
	verifyPort  func(port string) error
	lookupIP    func(host string) ([]net.IP, error)
}

func New(client network.Client, fileManager files.FileManager, firewall firewall.Configurator) pia {
	return &configurator{
		client:      client,
		fileManager: fileManager,
		firewall:    firewall,
		random:      random.NewRandom(),
		verifyPort:  verification.NewVerifier().VerifyPort,
		lookupIP:    net.LookupIP}
}
