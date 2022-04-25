package surfshark

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (s *Surfshark) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ciphers: []string{
			constants.AES256gcm,
		},
		Auth:          constants.SHA512,
		RenegDisabled: true,
		KeyDirection:  "1",
		Ping:          15,   //nolint:gomnd
		MssFix:        1450, //nolint:gomnd
		TunMTUExtra:   32,   //nolint:gomnd
		CA:            constants.SurfsharkCA,
		TLSAuth:       constants.SurfsharkTLSAuth,
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
