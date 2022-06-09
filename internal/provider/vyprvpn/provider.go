package vyprvpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/utils"
	"github.com/qdm12/gluetun/internal/provider/vyprvpn/updater"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	utils.NoPortForwarder
	common.Fetcher
}

func New(storage common.Storage, randSource rand.Source,
	unzipper common.Unzipper, updaterWarner common.Warner) *Provider {
	return &Provider{
		storage:         storage,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(providers.Vyprvpn),
		Fetcher:         updater.New(unzipper, updaterWarner),
	}
}

func (p *Provider) Name() string {
	return providers.Vyprvpn
}
