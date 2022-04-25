package mullvad

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (m *Mullvad) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	//nolint:gomnd
	providerSettings := utils.OpenVPNProviderSettings{
		AuthUserPass: true,
		Ciphers: []string{
			constants.AES256cbc,
			constants.AES128gcm,
		},
		Ping:          10,
		RemoteCertTLS: true,
		TLSCipher:     "TLS-DHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA",
		SndBuf:        524288,
		RcvBuf:        524288,
		CA:            constants.MullvadCA,
		UDPLines:      []string{"fast-io"},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
