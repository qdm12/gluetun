package custom

import (
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/openvpn/extract"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

type Provider struct {
	extractor extract.Interface
	utils.NoPortForwarder
}

func New() *Provider {
	return &Provider{
		extractor:       extract.New(),
		NoPortForwarder: utils.NewNoPortForwarding(constants.Custom),
	}
}
