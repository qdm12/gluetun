package expressvpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	//nolint:gomnd
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ciphers: []string{
			constants.AES256cbc,
		},
		Auth:           constants.SHA512,
		CA:             constants.ExpressvpnCA,
		Cert:           constants.ExpressvpnCert,
		RSAKey:         constants.ExpressvpnRSAKey,
		TLSAuth:        constants.ExpressvpnTLSAuth,
		MssFix:         1200,
		FastIO:         true,
		Fragment:       1300,
		SndBuf:         524288,
		RcvBuf:         524288,
		KeyDirection:   "1",
		VerifyX509Type: "name-prefix",
		// Always verify against `Server` x509 name prefix, security hole I guess?
		VerifyX509Name: "Server",
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
