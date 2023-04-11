package vpnunlimited

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
		AuthUserPass:  false,
		Ping:          5, //nolint:gomnd
		RenegDisabled: true,
		CA:            "MIIEjjCCA/egAwIBAgIJAKsVbHBdakXuMA0GCSqGSIb3DQEBBQUAMIHgMQswCQYDVQQGEwJVUzELMAkGA1UECBMCTlkxETAPBgNVBAcTCE5ldyBZb3JrMR8wHQYDVQQKExZTaW1wbGV4IFNvbHV0aW9ucyBJbmMuMRYwFAYDVQQLEw1WcG4gVW5saW1pdGVkMSMwIQYDVQQDExpzZXJ2ZXIudnBudW5saW1pdGVkYXBwLmNvbTEjMCEGA1UEKRMac2VydmVyLnZwbnVubGltaXRlZGFwcC5jb20xLjAsBgkqhkiG9w0BCQEWH3N1cHBvcnRAc2ltcGxleHNvbHV0aW9uc2luYy5jb20wHhcNMTMxMjE2MTM1OTQ0WhcNMjMxMjE0MTM1OTQ0WjCB4DELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAk5ZMREwDwYDVQQHEwhOZXcgWW9yazEfMB0GA1UEChMWU2ltcGxleCBTb2x1dGlvbnMgSW5jLjEWMBQGA1UECxMNVnBuIFVubGltaXRlZDEjMCEGA1UEAxMac2VydmVyLnZwbnVubGltaXRlZGFwcC5jb20xIzAhBgNVBCkTGnNlcnZlci52cG51bmxpbWl0ZWRhcHAuY29tMS4wLAYJKoZIhvcNAQkBFh9zdXBwb3J0QHNpbXBsZXhzb2x1dGlvbnNpbmMuY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDADUzS8QWGvdhLFKsEzAiq5+b0ukKjBza0k6/dmCeYTvCVqGKg/h1IAtQdVVLAkmEp0zvGH7PuOhXm7zZrCouBr/RiG4tEcoRHwc6AJmowkYERlY7+xGx3OuNrD00QceNTsan0bn78jwt0zhFNmHdoTtFjgK3oqmQYSAtbEVWYgwIDAQABo4IBTDCCAUgwHQYDVR0OBBYEFKClmYP+tMNgWagUJCCHjtaui2k+MIIBFwYDVR0jBIIBDjCCAQqAFKClmYP+tMNgWagUJCCHjtaui2k+oYHmpIHjMIHgMQswCQYDVQQGEwJVUzELMAkGA1UECBMCTlkxETAPBgNVBAcTCE5ldyBZb3JrMR8wHQYDVQQKExZTaW1wbGV4IFNvbHV0aW9ucyBJbmMuMRYwFAYDVQQLEw1WcG4gVW5saW1pdGVkMSMwIQYDVQQDExpzZXJ2ZXIudnBudW5saW1pdGVkYXBwLmNvbTEjMCEGA1UEKRMac2VydmVyLnZwbnVubGltaXRlZGFwcC5jb20xLjAsBgkqhkiG9w0BCQEWH3N1cHBvcnRAc2ltcGxleHNvbHV0aW9uc2luYy5jb22CCQCrFWxwXWpF7jAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBBQUAA4GBALkWhfw7SSV7it+ZYZmT+cQbExjlYgQ40zae2J2CqIYACRcfsDHvh7Q+fiwSocevv2NE0dWi6WB2H6SiudYjvDvubAX/QbzfBxtbxCCoAIlfPCm8xOnWFN7TUJAzWwHJkKgEnu29GZHu2x8J+7VeDbKH5RTMHHe8FkSxh6Zz/BMN", //nolint:lll
		ExtraLines: []string{
			"route-metric 1",
		},
	}

	if settings.Version != openvpn.Openvpn24 {
		// VPN Unlimited's certificate is sha1WithRSAEncryption and sha1 is now
		// rejected by openssl 3.x.x which is used by OpenVPN >= 2.5.
		// We lower the security level to 0 to allow this algorithm,
		// see https://www.openssl.org/docs/man1.1.1/man3/SSL_CTX_set_security_level.html
		providerSettings.TLSCipher = `"DEFAULT:@SECLEVEL=0"`
	}

	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
