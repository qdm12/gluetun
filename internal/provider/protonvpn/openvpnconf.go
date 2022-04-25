package protonvpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Protonvpn) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	//nolint:gomnd
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ciphers: []string{
			constants.AES256cbc,
		},
		Auth:          constants.SHA512,
		MssFix:        1450,
		TunMTUExtra:   32,
		RenegDisabled: true,
		KeyDirection:  "1",
		CA:            constants.ProtonvpnCA,
		TLSAuth:       constants.ProtonvpnTLSAuth,
		UDPLines: []string{
			"fast-io",
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
