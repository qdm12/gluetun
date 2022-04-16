package surfshark

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Surfshark struct {
	servers    []models.SurfsharkServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.SurfsharkServer, randSource rand.Source) *Surfshark {
	return &Surfshark{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(providers.Surfshark),
	}
}
