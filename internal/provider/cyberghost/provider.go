package cyberghost

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Cyberghost struct {
	servers    []models.CyberghostServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.CyberghostServer, randSource rand.Source) *Cyberghost {
	return &Cyberghost{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(providers.Cyberghost),
	}
}
