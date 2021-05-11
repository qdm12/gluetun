package torguard

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/models"
)

type Torguard struct {
	servers    []models.TorguardServer
	randSource rand.Source
}

func New(servers []models.TorguardServer, randSource rand.Source) *Torguard {
	return &Torguard{
		servers:    servers,
		randSource: randSource,
	}
}
