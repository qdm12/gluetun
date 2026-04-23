package cyberghost

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/pkg/updaters/cyberghost"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	common.Fetcher
}

func New(storage common.Storage, randSource rand.Source,
	updaterWarner common.Warner, parallelResolver common.ParallelResolver,
) *Provider {
	return &Provider{
		storage:    storage,
		randSource: randSource,
		Fetcher:    cyberghost.New(parallelResolver, updaterWarner),
	}
}

func (p *Provider) Name() string {
	return providers.Cyberghost
}
