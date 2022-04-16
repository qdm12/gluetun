package expressvpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Provider struct {
	servers    []models.ExpressvpnServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.ExpressvpnServer, randSource rand.Source) *Provider {
	return &Provider{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(providers.Expressvpn),
	}
}
