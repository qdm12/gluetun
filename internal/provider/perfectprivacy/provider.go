package perfectprivacy

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Perfectprivacy struct {
	servers    []models.PerfectprivacyServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.PerfectprivacyServer, randSource rand.Source) *Perfectprivacy {
	return &Perfectprivacy{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(constants.Perfectprivacy),
	}
}
