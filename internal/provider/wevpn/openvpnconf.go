package wevpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (w *Wevpn) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ciphers: []string{
			constants.AES256gcm,
		},
		Auth:          constants.SHA512,
		Ping:          30, //nolint:gomnd
		RenegDisabled: true,
		CA:            constants.WevpnCA,
		Cert:          constants.WevpnCert,
		TLSCrypt:      constants.WevpnTLSCrypt,
		ExtraLines: []string{
			"redirect-gateway def1 bypass-dhcp",
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
