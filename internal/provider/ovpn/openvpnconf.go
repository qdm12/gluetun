package ovpn

import (
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/openvpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) OpenVPNConfig(connection models.Connection,
	settings settings.OpenVPN, ipv6Supported bool,
) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		AuthUserPass:  true,
		RemoteCertTLS: true,
		Ciphers: []string{
			openvpn.AES256gcm,
			openvpn.AES256cbc,
			openvpn.AES128gcm,
			openvpn.Chacha20Poly1305,
		},
		CAs: []string{
			"MIIEfTCCA2WgAwIBAgIJAK2aIWqpLj1/MA0GCSqGSIb3DQEBBQUAMIGFMQswCQYDVQQGEwJTRTESMBAGA1UECBMJU3RvY2tob2xtMRIwEAYDVQQHEwlTdG9ja2hvbG0xHDAaBgNVBAsTE0Zpcm1hIERhdmlkIFdpYmVyZ2gxEzARBgNVBAMTCm92cG4uc2UgY2ExGzAZBgkqhkiG9w0BCQEWDGluZm9Ab3Zwbi5zZTAeFw0xNDA4MTcxODIxMjlaFw0zNDA4MTIxODIxMjlaMIGFMQswCQYDVQQGEwJTRTESMBAGA1UECBMJU3RvY2tob2xtMRIwEAYDVQQHEwlTdG9ja2hvbG0xHDAaBgNVBAsTE0Zpcm1hIERhdmlkIFdpYmVyZ2gxEzARBgNVBAMTCm92cG4uc2UgY2ExGzAZBgkqhkiG9w0BCQEWDGluZm9Ab3Zwbi5zZTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMR+aP4GTuZwurZuOA2NYzMfqKyZi/TJcLEPlGTB/b4CWA9bTd8f0pHPrDAZsXIEayxxB58BIFNDNiybnbO15JN/QwlsqmA+aZX6mCSkScs/rRwasM6LDo8iGx+KmYEqAgzziONGbCMnlO+OaarXte7LhZ9X6Z/bryu4xq/i1v3raak13kXsrogtu4iDzxqJE/QhbNOi0yhCdlm5RYQjmlKGdPB9pNTgcakVI4HcngRYMzBlrGin0YkvWCdpx5FrDNeld7BSWrJMNYyvd+buaid0Fu1T9/P/Srj/8AiabKoaDyiGFbZdTnGfK+04lWRvwAmvazpqbUt5Omw634jJDuMCAwEAAaOB7TCB6jAdBgNVHQ4EFgQUEvJcHHcTiDtu7bAyZw+xaqg+xdIwgboGA1UdIwSBsjCBr4AUEvJcHHcTiDtu7bAyZw+xaqg+xdKhgYukgYgwgYUxCzAJBgNVBAYTAlNFMRIwEAYDVQQIEwlTdG9ja2hvbG0xEjAQBgNVBAcTCVN0b2NraG9sbTEcMBoGA1UECxMTRmlybWEgRGF2aWQgV2liZXJnaDETMBEGA1UEAxMKb3Zwbi5zZSBjYTEbMBkGCSqGSIb3DQEJARYMaW5mb0BvdnBuLnNlggkArZohaqkuPX8wDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQUFAAOCAQEAJmID6OyBJbV7ayPPgquojF+FICuDdOfGVKP828cyISxcbVA04VpD0QLYVb0k9pFUx0NbgX2SvRTiFhP7LcyS1HV9s+XLCb2WItPPsrdRTwtqU2n3TlCEzWA3WOcOCtT6JSkv1eelmx1JnP0gYJrDvDvRYBFctwWhtE0bineSQkZwN6980zkknADLAiHpeZSu/AMx7CGTwA6SmoFvpNBmHXDcfe/9ZqbbYfUfyPNe+0JbMrcv1elKi+6wlEkHFaEBphiZwGEbOX1CjUMcQFgW/cIp3n50Eiyx6ktuqimhyb59P4Nw8gqH452tTtE4MM/brA5y0Q0WFBRBojfZIbGWWQ==", //nolint:lll
		},
		TLSAuth:      "81782767e4d59c4464cc5d1896f1cf6015017d53ac62e2e3b94b889e00b2c69ddc01944fe1c6d895b4d80540502eb71910b8d785c9efa9e3182343532adffe1cfbb7bb6eae39c502da2748edf0fb89b8a20b0a1085cc1f06135037881bc0c4ad8f2c0f4f72d2ab466fb54af3d8264c5fddeb0f21aa0ca41863678f5fc4c44de4ca0926b36dfddc42c6f2fabd1694bdc8215b2d223b9c21dc6734c2c778093187afb8c33403b228b9af68b540c284f6d183bcc88bd41d47bd717996e499ce1cbbfa768a9723c19c58314c4d19cfed82e543ee92e73d38ad26d4fbec231c0f9f3b30773a5c87792e9bc7c34e8d7611002ebedd044e48a0f1f96527bfdcc940aa09", //nolint:lll
		KeyDirection: "1",
		ExtraLines: []string{
			"allow-compression asym",
			"replay-window 256",
		},
	}

	if strings.HasSuffix(connection.Hostname, "singapore.ovpn.com") {
		providerSettings.TLSCrypt = providerSettings.TLSAuth
		providerSettings.TLSAuth = ""
		providerSettings.KeyDirection = ""
	}

	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
