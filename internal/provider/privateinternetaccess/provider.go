package privateinternetaccess

import (
	"math/rand"
	"time"

	"github.com/qdm12/gluetun/internal/constants/openvpn"
	"github.com/qdm12/gluetun/internal/models"
)

type Provider struct {
	servers    []models.Server
	randSource rand.Source
	timeNow    func() time.Time
	// Port forwarding
	portForwardPath string
	authFilePath    string
}

func New(servers []models.Server, randSource rand.Source,
	timeNow func() time.Time) *Provider {
	const jsonPortForwardPath = "/gluetun/piaportforward.json"
	return &Provider{
		servers:         servers,
		timeNow:         timeNow,
		randSource:      randSource,
		portForwardPath: jsonPortForwardPath,
		authFilePath:    openvpn.AuthConf,
	}
}
