package privatevpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Privatevpn struct {
	servers    []models.PrivatevpnServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.PrivatevpnServer, randSource rand.Source) *Privatevpn {
	return &Privatevpn{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(providers.Privatevpn),
	}
}
