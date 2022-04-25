package vyprvpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (v *Vyprvpn) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ciphers: []string{
			constants.AES256cbc,
		},
		Auth:      constants.SHA256,
		Ping:      10, //nolint:gomnd
		CA:        constants.VyprvpnCA,
		TLSCipher: "TLS-ECDHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA256:TLS-DHE-RSA-WITH-AES-256-CBC-SHA", //nolint:lll
		ExtraLines: []string{
			"comp-lzo",
		},
		// VerifyX509Name: []string{"lu1.vyprvpn.com","name"},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
