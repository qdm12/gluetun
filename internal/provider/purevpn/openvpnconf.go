package purevpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Purevpn) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ciphers: []string{
			constants.AES256gcm,
		},
		Ping:    10, //nolint:gomnd
		CA:      constants.PurevpnCA,
		Cert:    constants.PurevpnCert,
		Key:     constants.PurevpnKey,
		TLSAuth: constants.PurevpnTLSAuth,
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
