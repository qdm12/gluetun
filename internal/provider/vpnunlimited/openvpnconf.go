package vpnunlimited

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) OpenVPNConfig(connection models.Connection,
	settings settings.OpenVPN, ipv6Supported bool) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  false,
		Ping:          5, //nolint:gomnd
		RenegDisabled: true,
		CA:            "MIID7jCCA1CgAwIBAgIQQTT3w3N+5i8OMfe565xaSjAKBggqhkjOPQQDBDCBojELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAk5ZMREwDwYDVQQHDAhOZXcgWW9yazEXMBUGA1UECgwOS2VlcFNvbGlkIEluYy4xGjAYBgNVBAsMEUtlZXBTb2xpZCBSb290IENBMRowGAYDVQQDDBFLZWVwU29saWQgUm9vdCBDQTEiMCAGCSqGSIb3DQEJARYTYWRtaW5Aa2VlcHNvbGlkLmNvbTAeFw0yMDA0MDExNjI3MTRaFw0yNTAzMzExNjI3MTRaMIGgMQswCQYDVQQGEwJVUzELMAkGA1UECAwCTlkxETAPBgNVBAcMCE5ldyBZb3JrMRcwFQYDVQQKDA5LZWVwU29saWQgSW5jLjEVMBMGA1UECwwMS2VlcFNvbGlkIENBMR0wGwYDVQQDDBRPcGVuVlBOIFNlcnZlciBTdWJDQTEiMCAGCSqGSIb3DQEJARYTYWRtaW5Aa2VlcHNvbGlkLmNvbTCBmzAQBgcqhkjOPQIBBgUrgQQAIwOBhgAEAR9nmoZUraRSSPUhYwIQBLSx+phJdIlqU7F7Hszh95ivnWYkwuizKLaUYy6lSISDohlUtQl9URBlRrGroVctOGlOAdpL2ARTljw5gmUcaavc5cvLiAV7fPJ7BFUgVxInmaVcaMlDwGgKLxmjU2Fw85VLROHbWQjYc93x/BTSFcYO/np4o4IBIzCCAR8wDAYDVR0TBAUwAwEB/zAdBgNVHQ4EFgQUrUCjH8xe37lJihyzpqjWwxxNOiswgeIGA1UdIwSB2jCB14AU/LRRnTRaEbxct895Pk9DoymNQIqhgaikgaUwgaIxCzAJBgNVBAYTAlVTMQswCQYDVQQIDAJOWTERMA8GA1UEBwwITmV3IFlvcmsxFzAVBgNVBAoMDktlZXBTb2xpZCBJbmMuMRowGAYDVQQLDBFLZWVwU29saWQgUm9vdCBDQTEaMBgGA1UEAwwRS2VlcFNvbGlkIFJvb3QgQ0ExIjAgBgkqhkiG9w0BCQEWE2FkbWluQGtlZXBzb2xpZC5jb22CFEssZFYAz8WhYnIDxLeDgKTLD8p2MAsGA1UdDwQEAwIBBjAKBggqhkjOPQQDBAOBiwAwgYcCQgGuK8UNnpE8k8hAamnT9gxCSs5APqrgmdLe6BxYSz7AptpF2/MPzLFsXgj4YxC6vJP8Rs8e3Hw9VJ7DF0aYgu8DvQJBeyFWjRnk8kmu2zEU+wF9fkvN9AJ7v0xF0iEaFVsdPKv6sJQP1sAL+AIepJQ7TYvh9Q9G/WaRCfItCtcOAEz3SKA=", //nolint:lll
		ExtraLines: []string{
			"route-metric 1",
		},
	}

	// VPN Unlimited's certificate is sha1WithRSAEncryption and sha1 is now
	// rejected by openssl 3.x.x which is used by OpenVPN >= 2.5.
	// We lower the security level to 0 to allow this algorithm,
	// see https://www.openssl.org/docs/man1.1.1/man3/SSL_CTX_set_security_level.html
	providerSettings.TLSCipher = `"DEFAULT:@SECLEVEL=0"`

	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
