package slickvpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/openvpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) OpenVPNConfig(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	const pingSeconds = 10
	const bufSize = 393216
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ciphers: []string{
			openvpn.AES256cbc,
		},
		Ping:   pingSeconds,
		SndBuf: bufSize,
		RcvBuf: bufSize,
		CA:     "MIIDQDCCAqmgAwIBAgIJAM8Brk2pUr0KMA0GCSqGSIb3DQEBBQUAMHQxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDQTEMMAoGA1UEBxMDVlBOMQwwCgYDVQQKEwNWUE4xDDAKBgNVBAsTA1ZQTjEMMAoGA1UEAxMDVlBOMQwwCgYDVQQpEwNWUE4xEjAQBgkqhkiG9w0BCQEWA1ZQTjAeFw0xMjAzMDMwMjExNDJaFw0yMjAzMDEwMjExNDJaMHQxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDQTEMMAoGA1UEBxMDVlBOMQwwCgYDVQQKEwNWUE4xDDAKBgNVBAsTA1ZQTjEMMAoGA1UEAxMDVlBOMQwwCgYDVQQpEwNWUE4xEjAQBgkqhkiG9w0BCQEWA1ZQTjCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAwY2K08N7or1Br/EsD9XBon7gs7dKflWYuymgMLJfeMFWuJloNdsn+3GARIhYBbN6zhvFGFE214qKPqAydW1WmIIK7KoC0sgndr+Vk/au9gssFzVmmvr6+WN/nfo2L9KvvBMoYLrMAiyw/D4cRapZi2pXJLcMDfC+p1VWAX8TYWkCAwEAAaOB2TCB1jAdBgNVHQ4EFgQUmyvO4rTnu5/ABnp0FngU+SdR8WAwgaYGA1UdIwSBnjCBm4AUmyvO4rTnu5/ABnp0FngU+SdR8WCheKR2MHQxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDQTEMMAoGA1UEBxMDVlBOMQwwCgYDVQQKEwNWUE4xDDAKBgNVBAsTA1ZQTjEMMAoGA1UEAxMDVlBOMQwwCgYDVQQpEwNWUE4xEjAQBgkqhkiG9w0BCQEWA1ZQToIJAM8Brk2pUr0KMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAoB0kOuGvrzPBTIRXIDHCCxBMdny+3sKAOllmH4+51j2aWhAJ4Pyc/yBTYyQGNoriABjmNzp+R05oiaxAD3vTgR80juKDPtQb8LoGLBF18gL7Vtc3+hJXcJasXZaDSSoyh5f+TtGvytIT+eceJWIrKnFXzlHOvKlyLkcZn15gwKQ=", //nolint:lll
		ExtraLines: []string{
			"redirect-gateway",
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
