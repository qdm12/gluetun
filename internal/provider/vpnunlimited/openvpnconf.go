package vpnunlimited

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  false,
		Ping:          5, //nolint:gomnd
		RenegDisabled: true,
		CA:            constants.VPNUnlimitedCA,
		ExtraLines: []string{
			"route-metric 1",
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
