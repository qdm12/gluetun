package protonvpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Protonvpn struct {
	servers    []models.ProtonvpnServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.ProtonvpnServer, randSource rand.Source) *Protonvpn {
	return &Protonvpn{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(providers.Protonvpn),
	}
}
