package hidemyass

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type HideMyAss struct {
	servers    []models.Server
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.Server, randSource rand.Source) *HideMyAss {
	return &HideMyAss{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(providers.HideMyAss),
	}
}
