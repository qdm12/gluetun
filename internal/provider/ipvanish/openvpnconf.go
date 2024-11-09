package ipvanish

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/openvpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) OpenVPNConfig(connection models.Connection,
	settings settings.OpenVPN, ipv6Supported bool,
) (lines []string) {
	//nolint:mnd
	providerSettings := utils.OpenVPNProviderSettings{
		AuthUserPass: true,
		Ciphers: []string{
			openvpn.AES256gcm,
			openvpn.AES256cbc,
		},
		Auth:           openvpn.SHA256,
		VerifyX509Type: "name",
		TLSCipher:      "TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-DHE-DSS-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA",
		CAs:            []string{"MIIErzCCA5egAwIBAgIJAMYKzSS8uPKDMA0GCSqGSIb3DQEBDQUAMIGVMQswCQYDVQQGEwJVUzELMAkGA1UECBMCRkwxFDASBgNVBAcTC1dpbnRlciBQYXJrMREwDwYDVQQKEwhJUFZhbmlzaDEVMBMGA1UECxMMSVBWYW5pc2ggVlBOMRQwEgYDVQQDEwtJUFZhbmlzaCBDQTEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEBpcHZhbmlzaC5jb20wIBcNMjIwNTA5MjAyMDQ1WhgPMjA4MjA0MjQyMDIwNDVaMIGVMQswCQYDVQQGEwJVUzELMAkGA1UECBMCRkwxFDASBgNVBAcTC1dpbnRlciBQYXJrMREwDwYDVQQKEwhJUFZhbmlzaDEVMBMGA1UECxMMSVBWYW5pc2ggVlBOMRQwEgYDVQQDEwtJUFZhbmlzaCBDQTEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEBpcHZhbmlzaC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC30MFY2v8go65jdOYM/nHu9hlHQMbEttdTxPIDMFuNS0UUxuHGUeJdVCtkeaDOKH3jHsGBczu1amYwphVv6A1qox1YTrzRCbec7CaHL926VcOQQcDAPTmL+JPHhlpR21Xa+woHFGDW90LgASLAPtupXgc6LXfFwb3vVpDnkyPUp4J0DRo2+lq3UtbHaONbGx8jyzYu/kWSiLUc7X69OedoSwlmsGACQteki2o/b0uKTf84Ei+QEjGUquGJU+LETmo2IP55I+KuyZE6+zIiiegm25jgPDkrqlw2UrJiLCjUg4VhTdjF9/AUmT5tJbhZUGGx1/l0bGr+44ea7PmB3DELAgMBAAGjgf0wgfowDAYDVR0TBAUwAwEB/zAdBgNVHQ4EFgQUS/0UJYkd58Fwg9f2nxEcJU4Z7q4wgcoGA1UdIwSBwjCBv4AUS/0UJYkd58Fwg9f2nxEcJU4Z7q6hgZukgZgwgZUxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJGTDEUMBIGA1UEBxMLV2ludGVyIFBhcmsxETAPBgNVBAoTCElQVmFuaXNoMRUwEwYDVQQLEwxJUFZhbmlzaCBWUE4xFDASBgNVBAMTC0lQVmFuaXNoIENBMSMwIQYJKoZIhvcNAQkBFhRzdXBwb3J0QGlwdmFuaXNoLmNvbYIJAMYKzSS8uPKDMA0GCSqGSIb3DQEBDQUAA4IBAQCc9JV7IR8BfBrF/BQTXg0SZMZyyMAxR2jfW9qMHKSeJuZVVjfHiqoynEgBCNbn71wZWv3OF/Thu9BJ4GiYJ2Bc9nIa90D1NGYgiOVYLGXfUUqy5FgfrsWh0Go5oYm9l7W9pWfIifwsaZynkY0rTIHn32FF0H3+wZrGrEUzVL6qi+KD8iR3cBbLT+xUzulMTBp4JYaQnxpV4fZNS0ZsNrWKFWz4Iz1SSBcsnvUhfWs1aKx4yOJQx33Pc+KwpUI+meTlMjoh+AoTriooKU2MbOqLQl32y3pR0MP3fX4HDVFRylxdckEc+VryGNHQLUJiIBKBCORih/YiRhtEhpoBxmkw"}, //nolint:lll
		MssFix:         1320,
		ExtraLines: []string{
			"comp-lzo", // Explicitly disable compression
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
