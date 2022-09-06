package windscribe

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/openvpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) OpenVPNConfig(connection models.Connection,
	settings settings.OpenVPN, ipv6Supported bool) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ciphers: []string{
			openvpn.AES256gcm,
			openvpn.AES256cbc,
			openvpn.AES128gcm,
		},
		Auth:           openvpn.SHA512,
		Ping:           10, //nolint:gomnd
		VerifyX509Type: "name",
		KeyDirection:   "1",
		RenegDisabled:  true,
		CA:             "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", //nolint:lll
		TLSAuth:        "5801926a57ac2ce27e3dfd1dd6ef82042d82bd4f3f0021296f57734f6f1ea714a6623845541c4b0c3dea0a050fe6746cb66dfab14cda27e5ae09d7c155aa554f399fa4a863f0e8c1af787e5c602a801d3a2ec41e395a978d56729457fe6102d7d9e9119aa83643210b33c678f9d4109e3154ac9c759e490cb309b319cf708cae83ddadc3060a7a26564d1a24411cd552fe6620ea16b755697a4fc5e6e9d0cfc0c5c4a1874685429046a424c026db672e4c2c492898052ba59128d46200b40f880027a8b6610a4d559bdc9346d33a0a6b08e75c7fd43192b162bfd0aef0c716b31584827693f676f9a5047123466f0654eade34972586b31c6ce7e395f4b478cb",                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     //nolint:lll
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
