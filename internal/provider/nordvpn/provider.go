package nordvpn

import (
	"math/rand"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/pkg/updaters/nordvpn"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	common.Fetcher
}

func New(storage common.Storage, randSource rand.Source,
	client *http.Client, updaterWarner common.Warner,
) *Provider {
	return &Provider{
		storage:    storage,
		randSource: randSource,
		Fetcher:    nordvpn.New(client, updaterWarner),
	}
}

func (p *Provider) Name() string {
	return providers.Nordvpn
}
