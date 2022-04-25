package torguard

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (t *Torguard) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ciphers: []string{
			constants.AES256gcm,
		},
		Auth:          constants.SHA256,
		MssFix:        1450,   //nolint:gomnd
		TunMTUExtra:   32,     //nolint:gomnd
		SndBuf:        393216, //nolint:gomnd
		RcvBuf:        393216, //nolint:gomnd
		Ping:          5,      //nolint:gomnd
		RenegDisabled: true,
		KeyDirection:  "1",
		CA:            constants.TorguardCA,
		TLSAuth:       constants.TorguardTLSAuth,
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
