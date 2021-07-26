package vpnunlimited

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Provider struct {
	servers    []models.VPNUnlimitedServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.VPNUnlimitedServer, randSource rand.Source) *Provider {
	return &Provider{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(constants.VPNUnlimited),
	}
}
