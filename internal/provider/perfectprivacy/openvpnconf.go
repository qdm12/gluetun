package perfectprivacy

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Perfectprivacy) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	//nolint:gomnd
	providerSettings := utils.OpenVPNProviderSettings{
		AuthUserPass: true,
		Ciphers: []string{
			constants.AES256cbc,
			constants.AES256gcm,
		},
		Auth:         constants.SHA512,
		MssFix:       1450,
		Ping:         5,
		CA:           constants.PerfectprivacyCA,
		Cert:         constants.PerfectprivacyCert,
		Key:          constants.PerfectprivacyKey,
		TLSCrypt:     constants.PerfectprivacyTLSCrypt,
		TLSCipher:    "TLS_CHACHA20_POLY1305_SHA256:TLS-DHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-DHE-RSA-WITH-AES-128-GCM-SHA256:TLS-DHE-RSA-WITH-AES-128-CBC-SHA:TLS_AES_256_GCM_SHA384:TLS-RSA-WITH-AES-256-CBC-SHA", //nolint:lll
		TunMTU:       1500,
		TunMTUExtra:  32,
		RenegSec:     3600,
		KeyDirection: "1",
		IPv6Lines: []string{
			"redirect-gateway def1",
			`pull-filter ignore "redirect-gateway def1 ipv6"`,
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
