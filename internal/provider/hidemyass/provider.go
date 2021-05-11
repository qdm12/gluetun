package hidemyass

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/models"
)

type HideMyAss struct {
	servers    []models.HideMyAssServer
	randSource rand.Source
}

func New(servers []models.HideMyAssServer, randSource rand.Source) *HideMyAss {
	return &HideMyAss{
		servers:    servers,
		randSource: randSource,
	}
}
