package mullvad

import (
	"math/rand"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/pkg/updaters/mullvad"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	common.Fetcher
}

func New(storage common.Storage, randSource rand.Source,
	client *http.Client,
) *Provider {
	return &Provider{
		storage:    storage,
		randSource: randSource,
		Fetcher:    mullvad.New(client),
	}
}

func (p *Provider) Name() string {
	return providers.Mullvad
}
