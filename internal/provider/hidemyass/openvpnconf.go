package hidemyass

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (h *HideMyAss) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	//nolint:gomnd
	providerSettings := utils.OpenVPNProviderSettings{
		AuthUserPass: true,
		Ciphers: []string{
			constants.AES256cbc,
		},
		RemoteCertTLS: true,
		CA:            constants.HideMyAssCA,
		Cert:          constants.HideMyAssCert,
		RSAKey:        constants.HideMyAssRSAKey,
		Ping:          5,
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
