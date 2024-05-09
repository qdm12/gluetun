package torguard

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
			openvpn.AES256gcm, // In case the OpenVPN server accepts it
			openvpn.AES128gcm, // For OpenVPN 2.6, see https://github.com/qdm12/gluetun/issues/2271#issuecomment-2103349935
			openvpn.AES128cbc, // For OpenVPN 2.5, see https://github.com/qdm12/gluetun/issues/2271#issuecomment-2103349935
		},
		Auth:         openvpn.SHA256,
		TunMTUExtra:  32, //nolint:gomnd
		KeyDirection: "1",
		CAs: []string{
			"MIIDMTCCAhmgAwIBAgIJAKnGGJK6qLqSMA0GCSqGSIb3DQEBCwUAMBQxEjAQBgNVBAMMCVRHLVZQTi1DQTAgFw0xOTA1MjExNDIzMTFaGA8yMDU5MDUxMTE0MjMxMVowFDESMBAGA1UEAwwJVEctVlBOLUNBMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlv0UgPD3xVAvhhP6q1HCmeAWbH+9HPkyQ2P6qM5oHY5dntjmq8YT48FZGHWv7+s9O47v6Bv7rEc4UwQx15cc2LByivX2JwmE8JACvNfwEnZXYAPq9WU3ZgRrAGvA09ItuLqK2fQ4A7h8bFhmyxCbSzP1sSIT/zJY6ebuh5rDQSMJRMaoI0t1zorEZ7PlEmh+o0w5GPs0D0vY50UcnEzB4GOdWC9pJREwEqppWYLN7RRdG8JyIqmA59mhARCnQFUo38HWic4trxFe71jtD7YInNV7ShQtg0S0sXo36Rqfz72Jo08qqI70dNs5DN1aGNkQ/tRK9DhL5DLmTkaCw7mEFQIDAQABo4GDMIGAMB0GA1UdDgQWBBR7DcymXBp6u/jAaZOPUjUhEyhXfjBEBgNVHSMEPTA7gBR7DcymXBp6u/jAaZOPUjUhEyhXfqEYpBYwFDESMBAGA1UEAwwJVEctVlBOLUNBggkAqcYYkrqoupIwDAYDVR0TBAUwAwEB/zALBgNVHQ8EBAMCAQYwDQYJKoZIhvcNAQELBQADggEBAE79ngbdSlP7IBbfnJ+2Ju7vqt9/GyhcsYtjibp6gsMUxKlD8HuvlSGj5kNO5wiwN7XXqsjYtJfdhmzzVbXksi8Fnbnfa8GhFl4IAjLJ5cxaWOxjr6wx2AhIs+BVVARjaU7iTK91RXJnl6u7UDHTkQylBTl7wgpMeG6GjhaHfcOL1t7D2w8x23cTO+p+n53P3cBq+9TiAUORdzXJvbCxlPMDSDArsgBjC57W7dtdnZo7gTfQG77JTDFBeSwPwLF7PjBB4S6rzU/4fcYwy83XKP6zDn9tgUJDnpFb/7jJ/PbNkK4BWYJp3XytOtt66v9SEKw+v/fJ+VkjU16vE/9Q3h4=",                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             //nolint:lll
			"MIIFwjCCA6qgAwIBAgIRAPqbeSF13PE019f4UOUhbx8wDQYJKoZIhvcNAQENBQAwPTERMA8GA1UECgwIVG9yR3VhcmQxKDAmBgNVBAMMH1Rvckd1YXJkIFByaXZhdGUgUm9vdCBDQSAxIDIwMjAwIBcNMjMwNjI1MTM0ODU2WhgPMjA1MzA2MTcxMzQ4NTZaMD0xETAPBgNVBAoMCFRvckd1YXJkMSgwJgYDVQQDDB9Ub3JHdWFyZCBQcml2YXRlIFJvb3QgQ0EgMSAyMDIwMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA1Z1zrVEDLc8DUUFsCGz0H3fOi+YVGeuHsmNvIlKDnLpXqPkKjfcFOxs1pwMNYr8fBBkBct9W2oh1G1DxYLfjM1K8hlZNY1fvRs6mRFAX/nj+poK0gT5n0uTD0vQ5j/AqHO2wXCQm1xa2lUb7WrIt0ixKpgglRCeZwTXV2p7f9JZUI+ORX0B1zrV83e1ruefK+RCd3vf2UKurvz+sm0DS8xAC4LBX8xh1kk7MiAsK3a1mTufHpYmjAyS736yi+1rSCDEb7hBI3QXAGVwRFrGofHhR409XfB7aYwJela+bxRW44UD5az0uaeBM0GJcexH1fwi9F7ExAdR0kwWbJYX70S1F8es0Ik1ZpsLo2UEHc2/ueQMfpaLUL4kWfZOKNWWFSSbXR1YxPHitBSH638v4GfyNadBtG8UpVZ0dpsR/3VDoWH+WmowmlwhOAr5S/qt/iXf+/l8aHh4E/5AN4yTM1cCX+5LnKFCfJoWaxShI3TKi6Iw/80JWfAXAV52OKErRRuQ2YM+sQnJu+0vlW3oeNSQD2JwvSs0RD0zMC6Q6kCQXuDXyogS5K9qBlMt7UKDfZgaNnfiYvHjDh1XeQDN2hWUm0fTf14SCz4Lo8uE+CfnJHjU3zwk4GLvF8cs8RXhf8uZ5V/QHLxX9tK7FmLiTD8q1/U2tuzNlHgJURt8beGkCAwEAAaOBujCBtzAOBgNVHQ8BAf8EBAMCAcYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUM5sAY05zw9+2R6IdqRB5uRRmxpkwdQYDVR0jBG4wbIAUM5sAY05zw9+2R6IdqRB5uRRmxpmhQaQ/MD0xETAPBgNVBAoMCFRvckd1YXJkMSgwJgYDVQQDDB9Ub3JHdWFyZCBQcml2YXRlIFJvb3QgQ0EgMSAyMDIwghEA+pt5IXXc8TTX1/hQ5SFvHzANBgkqhkiG9w0BAQ0FAAOCAgEAnHPYMbo5Tf3tCD8HKVoibt4dtd9wEUh/XDFg2RNM8caa9x32gZJXCXSDUatdHYabukrsYqZIIt/XeL0SB8KzCQVyiMIHadCZBKc8Va/ays9lP/Kky6f3jkbrT5t9IhyHYNDWkrXmY/gNXCPoeulRQ55R0I1g5ko/JwvNp6q/V3fwvcpJJaFSh/NTOvBGCPR/pnR8isgmjF/i7KcN/b8gvO4EiqCk4AVl30aDUJBDyjnisCk9JMS4JxAYkJ9MGkqI1wHno3eKqBWoUEtyNe58VFQwxUgSf8cTV+p6DEZaM14qqDXzIQ3kHdGTH5ciqlzok0ocUM3AXvpHyoPbMPIFJ1uNvrYBWyDeP/KT512VNjpW30GtfMzZXJ2sEkcMAxghdqHxeKkOWVSsHHglHhq2qHsGF7eTZO1CFkV6kL0sn8shlPiJ/EE1//0tXycWstBaTe1TpiOYjLiLpwJvu7oMQIrl/YtCi/tXfkl8BLG0hncCLUovsIqQdjpo6jMux8p7D8L7yDV9GuQGxoT542GM53o83/esHhDSEMzDydH/cvpht/b9/YOzBxTMcxdxL8RDOKommtIfro1VE2z0YJ0KURD7jZe9mygV2KXokIBG4V+vhOglb7hT//drKFz6GDZAqs/KKeUIZxUWlpPaNssJygwDq6EjlNdelrxdWIYtR9Y=", //nolint:lll
		},
		TLSAuth: "770e8de5fc56e0248cc7b5aab56be80d0e19cbf003c1b3ed68efbaf08613c3a1a019dac6a4b84f13a6198f73229ffc21fa512394e288f82aa2cf0180f01fb3eb1a71e00a077a20f6d7a83633f5b4f47f27e30617eaf8485dd8c722a8606d56b3c183f65da5d3c9001a8cbdb96c793d936251098b24fe52a6dd2472e98cfccbc466e63520d63ade7a0eacc36208c3142a1068236a52142fbb7b3ed83d785e12a28261bccfb3bcb62a8d2f6d18f5df5f3652e59c5627d8d9c8f7877c4d7b08e19a5c363556ba68d392be78b75152dd55ba0f74d45089e84f77f4492d886524ea6c82b9f4dd83d46528d4f5c3b51cfeaf2838d938bd0597c426b0e440434f2c451f", //nolint:lll
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
