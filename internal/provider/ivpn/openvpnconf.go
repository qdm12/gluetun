package ivpn

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
		AuthUserPass: true,
		Ciphers: []string{
			openvpn.AES256cbc,
		},
		Ping:           5,
		RemoteCertTLS:  true,
		VerifyX509Type: "name-prefix",
		TLSCipher:      "TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-DHE-DSS-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA",
		CA:             "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", //nolint:lll
		TLSAuth:        "ac470c93ff9f5602a8aab37dee84a52814d10f20490ad23c47d5d82120c1bf859e93d0696b455d4a1b8d55d40c2685c41ca1d0aef29a3efd27274c4ef09020a3978fe45784b335da6df2d12db97bbb838416515f2a96f04715fd28949c6fe296a925cfada3f8b8928ed7fc963c1563272f5cf46e5e1d9c845d7703ca881497b7e6564a9d1dea9358adffd435295479f47d5298fabf5359613ff5992cb57ff081a04dfb81a26513a6b44a9b5490ad265f8a02384832a59cc3e075ad545461060b7bcab49bac815163cb80983dd51d5b1fd76170ffd904d8291071e96efc3fb777856c717b148d08a510f5687b8a8285dcffe737b98916dd15ef6235dee4266d3b",                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             //nolint:lll
		ExtraLines: []string{
			"key-direction 1",
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
