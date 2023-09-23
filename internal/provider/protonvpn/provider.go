package protonvpn

import (
	"github.com/qdm12/gluetun/internal/provider/utils"
	"math/rand"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/protonvpn/updater"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	utils.NoWireguardConfigurator
	common.Fetcher
}

func New(storage common.Storage, randSource rand.Source,
	client *http.Client, updaterWarner common.Warner) *Provider {
	return &Provider{
		storage:    storage,
		randSource: randSource,
		Fetcher:    updater.New(client, updaterWarner),
	}
}

func (p *Provider) Name() string {
	return providers.Protonvpn
}
