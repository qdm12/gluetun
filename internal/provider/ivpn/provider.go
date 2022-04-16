package ivpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Ivpn struct {
	servers    []models.IvpnServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.IvpnServer, randSource rand.Source) *Ivpn {
	return &Ivpn{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(providers.Ivpn),
	}
}
