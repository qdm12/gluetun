package custom

import (
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Provider struct {
	extractor Extractor
	utils.NoPortForwarder
	common.Fetcher
}

func New(extractor Extractor) *Provider {
	return &Provider{
		extractor:       extractor,
		NoPortForwarder: utils.NewNoPortForwarding(providers.Custom),
		Fetcher:         utils.NewNoFetcher(providers.Custom),
	}
}

func (p *Provider) Name() string {
	return providers.Custom
}
