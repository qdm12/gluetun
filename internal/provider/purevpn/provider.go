package purevpn

import (
	"math/rand"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/purevpn/updater"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	common.Fetcher
}

func New(storage common.Storage, randSource rand.Source,
	client *http.Client,
	ipFetcher common.IPFetcher, unzipper common.Unzipper,
	updaterWarner common.Warner, parallelResolver common.ParallelResolver,
) *Provider {
	return &Provider{
		storage:    storage,
		randSource: randSource,
		Fetcher:    updater.New(client, ipFetcher, unzipper, updaterWarner, parallelResolver),
	}
}

func (p *Provider) Name() string {
	return providers.Purevpn
}
