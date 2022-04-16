package vyprvpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Vyprvpn struct {
	servers    []models.VyprvpnServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.VyprvpnServer, randSource rand.Source) *Vyprvpn {
	return &Vyprvpn{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(providers.Vyprvpn),
	}
}
