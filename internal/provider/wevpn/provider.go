package wevpn

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/utils"
	"github.com/qdm12/gluetun/internal/provider/wevpn/updater"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	utils.NoPortForwarder
	common.Fetcher
}

func New(storage common.Storage, randSource rand.Source,
	updaterWarner common.Warner) *Provider {
	return &Provider{
		storage:         storage,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(providers.Wevpn),
		Fetcher:         updater.New(updaterWarner),
	}
}

func (p *Provider) Name() string {
	return providers.Wevpn
}
