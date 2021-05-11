package mullvad

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/models"
)

type Mullvad struct {
	servers    []models.MullvadServer
	randSource rand.Source
}

func New(servers []models.MullvadServer, randSource rand.Source) *Mullvad {
	return &Mullvad{
		servers:    servers,
		randSource: randSource,
	}
}
