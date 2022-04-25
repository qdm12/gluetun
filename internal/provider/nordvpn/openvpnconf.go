package nordvpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (n *Nordvpn) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	//nolint:gomnd
	providerSettings := utils.OpenVPNProviderSettings{
		AuthUserPass: true,
		Ciphers: []string{
			constants.AES256cbc,
			constants.AES256gcm,
		},
		Auth:          constants.SHA512,
		Ping:          15,
		RemoteCertTLS: true,
		MssFix:        1450,
		CA:            constants.NordvpnCA,
		TLSAuth:       constants.NordvpnTLSAuth,
		TunMTUExtra:   32,
		RenegDisabled: true,
		KeyDirection:  "1",
		UDPLines: []string{
			"fast-io",
		},
		ExtraLines: []string{
			"comp-lzo no", // Explicitly disable compression
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
