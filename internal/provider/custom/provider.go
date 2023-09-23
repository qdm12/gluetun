package custom

import (
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Provider struct {
	extractor Extractor
	common.Fetcher
}

func New(extractor Extractor) *Provider {
	return &Provider{
		extractor: extractor,
		Fetcher:   utils.NewNoFetcher(providers.Custom),
	}
}

func (p *Provider) Name() string {
	return providers.Custom
}
