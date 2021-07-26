package fastestvpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Fastestvpn struct {
	servers    []models.FastestvpnServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.FastestvpnServer, randSource rand.Source) *Fastestvpn {
	return &Fastestvpn{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(constants.Fastestvpn),
	}
}
