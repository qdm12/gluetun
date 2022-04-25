package privatevpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Privatevpn) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ciphers: []string{
			constants.AES128gcm,
		},
		Auth:    constants.SHA256,
		CA:      constants.PrivatevpnCA,
		TLSAuth: constants.PrivatevpnTLSAuth,
		UDPLines: []string{
			"key-direction 1",
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
