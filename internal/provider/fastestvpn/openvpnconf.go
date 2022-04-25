package fastestvpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (f *Fastestvpn) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	//nolint:gomnd
	providerSettings := utils.OpenVPNProviderSettings{
		AuthUserPass: true,
		Ciphers: []string{
			constants.AES256cbc,
		},
		MssFix:        1450,
		TLSCipher:     "TLS-DHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA256:TLS-DHE-RSA-WITH-CAMELLIA-256-CBC-SHA:TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-CAMELLIA-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA", //nolint:lll
		AuthToken:     true,
		KeyDirection:  "1",
		RenegDisabled: true,
		CA:            constants.FastestvpnCA,
		TLSAuth:       constants.FastestvpnTLSAuth,
		UDPLines: []string{
			"tun-mtu 1500",
			"tun-mtu-extra 32",
			"ping 15",
		},
		ExtraLines: []string{
			"comp-lzo",
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
