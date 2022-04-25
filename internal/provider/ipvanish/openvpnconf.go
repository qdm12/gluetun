package ipvanish

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (i *Ipvanish) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		AuthUserPass: true,
		Ciphers: []string{
			constants.AES256cbc,
		},
		Auth:           constants.SHA256,
		VerifyX509Type: "name",
		TLSCipher:      "TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-DHE-DSS-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA",
		CA:             constants.IpvanishCA,
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
