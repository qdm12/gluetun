package vpnunlimited

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/models"
)

type Provider struct {
	servers    []models.VPNUnlimitedServer
	randSource rand.Source
}

func New(servers []models.VPNUnlimitedServer, randSource rand.Source) *Provider {
	return &Provider{
		servers:    servers,
		randSource: randSource,
	}
}
