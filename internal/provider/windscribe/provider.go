package windscribe

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Windscribe struct {
	servers    []models.WindscribeServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.WindscribeServer, randSource rand.Source) *Windscribe {
	return &Windscribe{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(providers.Windscribe),
	}
}
