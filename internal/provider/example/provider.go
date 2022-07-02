package example

import (
	"math/rand"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/example/updater"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	utils.NoPortForwarder
	common.Fetcher
}

// TODO: remove unneeded arguments once the updater is implemented.
func New(storage common.Storage, randSource rand.Source,
	updaterWarner common.Warner, client *http.Client,
	unzipper common.Unzipper, parallelResolver common.ParallelResolver) *Provider {
	return &Provider{
		storage:         storage,
		randSource:      randSource,
		NoPortForwarder: utils.NewNoPortForwarding(providers.Example),
		Fetcher:         updater.New(updaterWarner, unzipper, client, parallelResolver),
	}
}

func (p *Provider) Name() string {
	// TODO: update the constant to be the right provider name.
	return providers.Example
}
