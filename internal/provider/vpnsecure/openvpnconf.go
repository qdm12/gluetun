package vpnsecure

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
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ping:          10,
		// note DES-CBC is not added since it's quite unsecure
		Ciphers: []string{openvpn.AES256cbc, openvpn.AES128cbc},
		ExtraLines: []string{
			"comp-lzo",
			"float",
		},
		CA: "MIIEJjCCAw6gAwIBAgIJAMkzh6p4m6XfMA0GCSqGSIb3DQEBCwUAMGkxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJOWTERMA8GA1UEBxMITmV3IFlvcmsxFTATBgNVBAoTDHZwbnNlY3VyZS5tZTEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEB2cG5zZWN1cmUubWUwIBcNMTcwNTA2MTMzMTQyWhgPMjkzODA4MjYxMzMxNDJaMGkxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJOWTERMA8GA1UEBxMITmV3IFlvcmsxFTATBgNVBAoTDHZwbnNlY3VyZS5tZTEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEB2cG5zZWN1cmUubWUwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDiClT1wcZ6oovYjSxUJIQplrBSQRKB44uymC8evohzK7q67x0NE2sLz5Zn9ZiC7RnXQCtEqJfHqjuqjaH5MghjhUDnRbZS/8ElxdGKn9FPvs9b+aTVGSfrQm5KKoVigwAye3ilNiWAyy6MDlBeoKluQ4xW7SGiVZRxLcJbLAmjmfCjBS7eUGbtA8riTkIegFo4WFiy9G76zQWw1V26kDhyzcJNT4xO7USMPUeZthy13g+zi9+rcILhEAnl776sIil6w8UVK8xevFKBlOPk+YyXlo4eZiuppq300ogaS+fX/0mfD7DDE+Gk5/nCeACDNiBlfQ3ol/De8Cm60HWEUtZVAgMBAAGjgc4wgcswHQYDVR0OBBYEFBJyf4mpGT3dIu65/1zAFqCgGxZoMIGbBgNVHSMEgZMwgZCAFBJyf4mpGT3dIu65/1zAFqCgGxZooW2kazBpMQswCQYDVQQGEwJVUzELMAkGA1UECBMCTlkxETAPBgNVBAcTCE5ldyBZb3JrMRUwEwYDVQQKEwx2cG5zZWN1cmUubWUxIzAhBgkqhkiG9w0BCQEWFHN1cHBvcnRAdnBuc2VjdXJlLm1lggkAyTOHqnibpd8wDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEArbTAibGQilY4Lu2RAVPjNx14SfojueBroeN7NIpAFUfbifPQRWvLamzRfxFTO0PXRc2pw/It7oa8yM7BsZj0vOiZY2p1JBHZwKom6tiSUVENDGW6JaYtiaE8XPyjfA5Yhfx4FefmaJ1veDYid18S+VVpt+Y+UIUxNmg1JB3CCUwbjl+dWlcvDBy4+jI+sZ7A1LF3uX64ZucDQ/XrpuopHhvDjw7g1PpKXsRqBYL+cpxUI7GrINBa/rGvXqv/NvFH8bguggknWKxKhd+jyMqkW3Ws258e0OwHz7gQ+tTJ909tR0TxJhZGkHatNSbpwW1Y52A972+9gYJMadSfm4bUHA==", //nolint:lll
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
