package cyberghost

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (c *Cyberghost) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	//nolint:gomnd
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ciphers: []string{
			constants.AES256gcm,
			constants.AES256cbc,
			constants.AES128gcm,
		},
		Auth: constants.SHA256,
		Ping: 10,
		CA:   constants.CyberghostCA,
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
