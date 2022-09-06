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
			openvpn.AES256gcm,
		},
		Auth:          openvpn.SHA256,
		MssFix:        1450,   //nolint:gomnd
		TunMTUExtra:   32,     //nolint:gomnd
		SndBuf:        393216, //nolint:gomnd
		RcvBuf:        393216, //nolint:gomnd
		Ping:          5,      //nolint:gomnd
		RenegDisabled: true,
		KeyDirection:  "1",
		CA:            "MIIDMTCCAhmgAwIBAgIJAKnGGJK6qLqSMA0GCSqGSIb3DQEBCwUAMBQxEjAQBgNVBAMMCVRHLVZQTi1DQTAgFw0xOTA1MjExNDIzMTFaGA8yMDU5MDUxMTE0MjMxMVowFDESMBAGA1UEAwwJVEctVlBOLUNBMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlv0UgPD3xVAvhhP6q1HCmeAWbH+9HPkyQ2P6qM5oHY5dntjmq8YT48FZGHWv7+s9O47v6Bv7rEc4UwQx15cc2LByivX2JwmE8JACvNfwEnZXYAPq9WU3ZgRrAGvA09ItuLqK2fQ4A7h8bFhmyxCbSzP1sSIT/zJY6ebuh5rDQSMJRMaoI0t1zorEZ7PlEmh+o0w5GPs0D0vY50UcnEzB4GOdWC9pJREwEqppWYLN7RRdG8JyIqmA59mhARCnQFUo38HWic4trxFe71jtD7YInNV7ShQtg0S0sXo36Rqfz72Jo08qqI70dNs5DN1aGNkQ/tRK9DhL5DLmTkaCw7mEFQIDAQABo4GDMIGAMB0GA1UdDgQWBBR7DcymXBp6u/jAaZOPUjUhEyhXfjBEBgNVHSMEPTA7gBR7DcymXBp6u/jAaZOPUjUhEyhXfqEYpBYwFDESMBAGA1UEAwwJVEctVlBOLUNBggkAqcYYkrqoupIwDAYDVR0TBAUwAwEB/zALBgNVHQ8EBAMCAQYwDQYJKoZIhvcNAQELBQADggEBAE79ngbdSlP7IBbfnJ+2Ju7vqt9/GyhcsYtjibp6gsMUxKlD8HuvlSGj5kNO5wiwN7XXqsjYtJfdhmzzVbXksi8Fnbnfa8GhFl4IAjLJ5cxaWOxjr6wx2AhIs+BVVARjaU7iTK91RXJnl6u7UDHTkQylBTl7wgpMeG6GjhaHfcOL1t7D2w8x23cTO+p+n53P3cBq+9TiAUORdzXJvbCxlPMDSDArsgBjC57W7dtdnZo7gTfQG77JTDFBeSwPwLF7PjBB4S6rzU/4fcYwy83XKP6zDn9tgUJDnpFb/7jJ/PbNkK4BWYJp3XytOtt66v9SEKw+v/fJ+VkjU16vE/9Q3h4=", //nolint:lll
		TLSAuth:       "770e8de5fc56e0248cc7b5aab56be80d0e19cbf003c1b3ed68efbaf08613c3a1a019dac6a4b84f13a6198f73229ffc21fa512394e288f82aa2cf0180f01fb3eb1a71e00a077a20f6d7a83633f5b4f47f27e30617eaf8485dd8c722a8606d56b3c183f65da5d3c9001a8cbdb96c793d936251098b24fe52a6dd2472e98cfccbc466e63520d63ade7a0eacc36208c3142a1068236a52142fbb7b3ed83d785e12a28261bccfb3bcb62a8d2f6d18f5df5f3652e59c5627d8d9c8f7877c4d7b08e19a5c363556ba68d392be78b75152dd55ba0f74d45089e84f77f4492d886524ea6c82b9f4dd83d46528d4f5c3b51cfeaf2838d938bd0597c426b0e440434f2c451f",                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         //nolint:lll
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
