package ivpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/models"
)

type Ivpn struct {
	servers    []models.IvpnServer
	randSource rand.Source
}

func New(servers []models.IvpnServer, randSource rand.Source) *Ivpn {
	return &Ivpn{
		servers:    servers,
		randSource: randSource,
	}
}
