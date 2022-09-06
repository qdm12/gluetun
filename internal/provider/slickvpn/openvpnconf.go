package slickvpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/openvpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) OpenVPNConfig(connection models.Connection,
	settings settings.OpenVPN, ipv6Supported bool) (lines []string) {
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
		// Certificate found from https://www.slickvpn.com/tutorials/using-openvpn-configuration-files/
		CA: "MIIESDCCAzCgAwIBAgIJAKHK5bbBPSU2MA0GCSqGSIb3DQEBBQUAMHUxCzAJBgNVBAYTAlVTMQwwCgYDVQQIEwNWUE4xDDAKBgNVBAcTA1ZQTjEMMAoGA1UEChMDVlBOMQwwCgYDVQQLEwNWUE4xDDAKBgNVBAMTA1ZQTjEMMAoGA1UEKRMDVlBOMRIwEAYJKoZIhvcNAQkBFgNWUE4wHhcNMjIwMjE0MjEzNDQwWhcNMzIwMjEyMjEzNDQwWjB1MQswCQYDVQQGEwJVUzEMMAoGA1UECBMDVlBOMQwwCgYDVQQHEwNWUE4xDDAKBgNVBAoTA1ZQTjEMMAoGA1UECxMDVlBOMQwwCgYDVQQDEwNWUE4xDDAKBgNVBCkTA1ZQTjESMBAGCSqGSIb3DQEJARYDVlBOMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwUl1XkfGo3c1uFsvgbO3C3yvu0+cHs9IUSsju5U9cPNCo53mqRHU/qntCC+ldIDKN+dNWn7eURIDszy+flutkgucs0qgETy5fzpXnVMtiKmMiOYWiJDor7j7QivRaxoT/O2fyjxvVCL8gMa60ekWSGBT6isYY8t8BnwTPVP0KvDm36wdaqLmubhf2XGvka/hhNx0SXMmz2x3OJ8BcoypZVLLk/+Qm6DJh1KxyDi4kI+jBC41QuaKKDNwb0kth1304eqZoUeCXtGkzl91y76ODAfdqzXf9WYJdgkXpOm53Cg7FtB42AqPRqHJVwYxDnQyrFwy4a3CWqFJnKtxJM/WlwIDAQABo4HaMIHXMB0GA1UdDgQWBBRSzxAtISfbSRPr0fmhwNZc8kOeKzCBpwYDVR0jBIGfMIGcgBRSzxAtISfbSRPr0fmhwNZc8kOeK6F5pHcwdTELMAkGA1UEBhMCVVMxDDAKBgNVBAgTA1ZQTjEMMAoGA1UEBxMDVlBOMQwwCgYDVQQKEwNWUE4xDDAKBgNVBAsTA1ZQTjEMMAoGA1UEAxMDVlBOMQwwCgYDVQQpEwNWUE4xEjAQBgkqhkiG9w0BCQEWA1ZQToIJAKHK5bbBPSU2MAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAGuKFW765F3D5wax5IFSQbEtr+rVHgjR8jiYTzxOCmbLaU4oj6phOhfQJiTTADQYgCIC/DN0HsAEEqrKkwEn8KdAoNiAWfqCV/eqnK83y7yRDGx6/zfsch+PAzKZouMJLrvR9RYbHq7m3adZv84YLge7FE1JqFk1j6rSa4dUUnGQPrQgr9Sasnz8O8KK45XH6fqKrsd4p485n+BXPDzWVsHl4M5dqQV7qUZTazbzzh4NyP5/RQ6Oh5jqMN7po4qiqWv1pu0EKDxUG6gcECc2cTQwHhIOPeCGdHS7uuI2FlLnHaCUFBgi8zTsZxaeaPuPch5M7Zxbdg0GBhS2SmNi+io=", //nolint:lll
		ExtraLines: []string{
			"redirect-gateway",
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
