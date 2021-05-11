package cyberghost

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/models"
)

type Cyberghost struct {
	servers    []models.CyberghostServer
	randSource rand.Source
}

func New(servers []models.CyberghostServer, randSource rand.Source) *Cyberghost {
	return &Cyberghost{
		servers:    servers,
		randSource: randSource,
	}
}
