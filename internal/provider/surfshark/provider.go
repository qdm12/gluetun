package surfshark

import (
	"math/rand"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/pkg/updaters/surfshark"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	common.Fetcher
}

func New(storage common.Storage, randSource rand.Source,
	client *http.Client, unzipper common.Unzipper, updaterWarner common.Warner,
	parallelResolver common.ParallelResolver,
) *Provider {
	return &Provider{
		storage:    storage,
		randSource: randSource,
		Fetcher:    surfshark.New(client, unzipper, updaterWarner, parallelResolver),
	}
}

func (p *Provider) Name() string {
	return providers.Surfshark
}
