package vyprvpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/openvpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) OpenVPNConfig(connection models.Connection,
	settings settings.OpenVPN, ipv6Supported bool) (lines []string) {
	//nolint:gomnd
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ciphers: []string{
			openvpn.AES256cbc,
		},
		Auth:      openvpn.SHA256,
		Ping:      10,
		CA:        "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", //nolint:lll
		TLSCipher: "TLS-ECDHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA256:TLS-DHE-RSA-WITH-AES-256-CBC-SHA",                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               //nolint:lll
		ExtraLines: []string{
			"comp-lzo",
		},
		// VerifyX509Name: []string{"lu1.vyprvpn.com","name"},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
