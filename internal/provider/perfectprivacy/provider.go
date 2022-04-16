package perfectprivacy

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Perfectprivacy struct {
	servers    []models.Server
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.Server, randSource rand.Source) *Perfectprivacy {
	return &Perfectprivacy{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(providers.Perfectprivacy),
	}
}
