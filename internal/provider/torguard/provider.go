package torguard

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Torguard struct {
	servers    []models.TorguardServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.TorguardServer, randSource rand.Source) *Torguard {
	return &Torguard{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(constants.Torguard),
	}
}
