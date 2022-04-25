package ivpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (i *Ivpn) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	//nolint:gomnd
	providerSettings := utils.OpenVPNProviderSettings{
		AuthUserPass: true,
		Ciphers: []string{
			constants.AES256cbc,
		},
		Ping:           5,
		RemoteCertTLS:  true,
		VerifyX509Type: "name-prefix",
		TLSCipher:      "TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-DHE-DSS-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA",
		CA:             constants.IvpnCA,
		TLSAuth:        constants.IvpnTLSAuth,
		ExtraLines: []string{
			"key-direction 1",
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
