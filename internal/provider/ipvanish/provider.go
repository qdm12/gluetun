package ipvanish

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Ipvanish struct {
	servers    []models.IpvanishServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.IpvanishServer, randSource rand.Source) *Ipvanish {
	return &Ipvanish{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(constants.Ipvanish),
	}
}
