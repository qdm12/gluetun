package custom

import (
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/openvpn/extract"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Provider struct {
	extractor extractor
	utils.NoPortForwarder
	common.Fetcher
}

func New() *Provider {
	return &Provider{
		extractor:       extract.New(),
		NoPortForwarder: utils.NewNoPortForwarding(providers.Custom),
		Fetcher:         utils.NewNoFetcher(providers.Custom),
	}
}

func (p *Provider) Name() string {
	return providers.Custom
}
