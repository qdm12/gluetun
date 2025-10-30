package protonvpn

import (
	"math/rand"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/protonvpn/updater"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	common.Fetcher
	portForwarded uint16
}

func New(storage common.Storage, randSource rand.Source,
	client *http.Client, updaterWarner common.Warner,
	username, password string,
) *Provider {
	return &Provider{
		storage:    storage,
		randSource: randSource,
		Fetcher:    updater.New(client, updaterWarner, username, password),
	}
}

func (p *Provider) Name() string {
	return providers.Protonvpn
}
