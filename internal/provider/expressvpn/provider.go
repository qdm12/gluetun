package expressvpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/expressvpn/updater"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	utils.NoPortForwarder
	utils.NoWireguardConfigurator
	common.Fetcher
}

func New(storage common.Storage, randSource rand.Source,
	unzipper common.Unzipper, updaterWarner common.Warner,
	parallelResolver common.ParallelResolver) *Provider {
	return &Provider{
		storage:         storage,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(providers.Expressvpn),
		Fetcher:         updater.New(unzipper, updaterWarner, parallelResolver),
	}
}

func (p *Provider) Name() string {
	return providers.Expressvpn
}
