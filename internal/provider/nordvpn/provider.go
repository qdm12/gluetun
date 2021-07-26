package nordvpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Nordvpn struct {
	servers    []models.NordvpnServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.NordvpnServer, randSource rand.Source) *Nordvpn {
	return &Nordvpn{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(constants.Nordvpn),
	}
}
