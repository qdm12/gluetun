package airvpn

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/openvpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) OpenVPNConfig(connection models.Connection,
	settings settings.OpenVPN, ipv6Supported bool) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		AuthUserPass:  true,
		RemoteCertTLS: true,
		KeyDirection:  "1",
		CA:            "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", //nolint:lll
		TLSAuth:       "7bb7a23a0f5f28d01e792df68f1764abf2688719288808bf58e8a2d4f9354ecf132625dfb895fc3f6330ae1e868e4dfac164c0931593d7f9a7da9595cf3534338896e1d0a987a0d19838944af8fea4e5215a3a0c76f4c67d5a4aee6a53be66a4c88b84f850030840fb30f8550ed8068f35c1ef34ee8f40a0ea5862dfb6f8d3c57ab5e27ac2799cf93e8765ff63cd8cd86b391b813925cd373bb202796f64d16f003d042ca828d1b07f18ba1d0cb0323ddf3ee9287e9e084e655699efb3cffa923626946fa372e7beee245e7a95b4c1d87d16cae685218d4b8afc019b22e41083476ee9883fe666d236301e55b20625514d91c8a69467a758293994df1e6fa7ae",                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         //nolint:lll
		ExtraLines: []string{
			"comp-lzo no", // Explicitly disable compression
			"push-peer-info",
		},
	}

	switch settings.Version {
	case openvpn.Openvpn24:
		providerSettings.Ciphers = []string{openvpn.AES256cbc}
	case openvpn.Openvpn25:
		providerSettings.Ciphers = []string{
			openvpn.AES256gcm, openvpn.AES256cbc, openvpn.AES192gcm,
			openvpn.AES192cbc, openvpn.AES128gcm, openvpn.AES128cbc,
			openvpn.Chacha20Poly1305}
	default:
		panic(fmt.Sprintf("openvpn version %q is not implemented", settings.Version))
	}

	providerSettings.SetEnv = map[string]string{"UV_IPV6": "no"}
	if ipv6Supported {
		providerSettings.SetEnv["UV_IPV6"] = "yes"
	}

	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
