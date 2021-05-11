package privatevpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/models"
)

type Privatevpn struct {
	servers    []models.PrivatevpnServer
	randSource rand.Source
}

func New(servers []models.PrivatevpnServer, randSource rand.Source) *Privatevpn {
	return &Privatevpn{
		servers:    servers,
		randSource: randSource,
	}
}
