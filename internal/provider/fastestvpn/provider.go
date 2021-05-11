package fastestvpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/models"
)

type Fastestvpn struct {
	servers    []models.FastestvpnServer
	randSource rand.Source
}

func New(servers []models.FastestvpnServer, randSource rand.Source) *Fastestvpn {
	return &Fastestvpn{
		servers:    servers,
		randSource: randSource,
	}
}
