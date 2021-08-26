package wevpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Wevpn struct {
	servers    []models.WevpnServer
	randSource rand.Source
	utils.NoPortForwarder
}

func New(servers []models.WevpnServer, randSource rand.Source) *Wevpn {
	return &Wevpn{
		servers:         servers,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(constants.Wevpn),
	}
}
