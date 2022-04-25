package windscribe

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (w *Windscribe) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ciphers: []string{
			constants.AES256gcm,
			constants.AES256cbc,
			constants.AES128gcm,
		},
		Auth:           constants.SHA512,
		Ping:           10, //nolint:gomnd
		VerifyX509Type: "name",
		KeyDirection:   "1",
		RenegDisabled:  true,
		CA:             constants.WindscribeCA,
		TLSAuth:        constants.WindscribeTLSAuth,
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
