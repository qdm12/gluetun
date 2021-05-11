package purevpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/models"
)

type Purevpn struct {
	servers    []models.PurevpnServer
	randSource rand.Source
}

func New(servers []models.PurevpnServer, randSource rand.Source) *Purevpn {
	return &Purevpn{
		servers:    servers,
		randSource: randSource,
	}
}
