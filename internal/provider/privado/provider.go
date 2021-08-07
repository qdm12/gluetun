package privado

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Privado struct {
	servers    []models.PrivadoServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.PrivadoServer, randSource rand.Source) *Privado {
	return &Privado{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(constants.Privado),
	}
}
