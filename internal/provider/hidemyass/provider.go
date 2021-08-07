package hidemyass

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type HideMyAss struct {
	servers    []models.HideMyAssServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.HideMyAssServer, randSource rand.Source) *HideMyAss {
	return &HideMyAss{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(constants.HideMyAss),
	}
}
