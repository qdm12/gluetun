package mullvad

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
		AuthUserPass: true,
		Ciphers: []string{
			openvpn.AES256cbc,
			openvpn.AES128gcm,
		},
		Ping:          10,
		RemoteCertTLS: true,
		TLSCipher:     "TLS-DHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA",
		SndBuf:        524288,
		RcvBuf:        524288,
		CA:            "MIIGIzCCBAugAwIBAgIJAK6BqXN9GHI0MA0GCSqGSIb3DQEBCwUAMIGfMQswCQYDVQQGEwJTRTERMA8GA1UECAwIR290YWxhbmQxEzARBgNVBAcMCkdvdGhlbmJ1cmcxFDASBgNVBAoMC0FtYWdpY29tIEFCMRAwDgYDVQQLDAdNdWxsdmFkMRswGQYDVQQDDBJNdWxsdmFkIFJvb3QgQ0EgdjIxIzAhBgkqhkiG9w0BCQEWFHNlY3VyaXR5QG11bGx2YWQubmV0MB4XDTE4MTEwMjExMTYxMVoXDTI4MTAzMDExMTYxMVowgZ8xCzAJBgNVBAYTAlNFMREwDwYDVQQIDAhHb3RhbGFuZDETMBEGA1UEBwwKR290aGVuYnVyZzEUMBIGA1UECgwLQW1hZ2ljb20gQUIxEDAOBgNVBAsMB011bGx2YWQxGzAZBgNVBAMMEk11bGx2YWQgUm9vdCBDQSB2MjEjMCEGCSqGSIb3DQEJARYUc2VjdXJpdHlAbXVsbHZhZC5uZXQwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCifDn75E/Zdx1qsy31rMEzuvbTXqZVZp4bjWbmcyyXqvnayRUHHoovG+lzc+HDL3HJV+kjxKpCMkEVWwjY159lJbQbm8kkYntBBREdzRRjjJpTb6haf/NXeOtQJ9aVlCc4dM66bEmyAoXkzXVZTQJ8h2FE55KVxHi5Sdy4XC5zm0wPa4DPDokNp1qm3A9Xicq3HsflLbMZRCAGuI+Jek6caHqiKjTHtujn6Gfxv2WsZ7SjerUAk+mvBo2sfKmB7octxG7yAOFFg7YsWL0AxddBWqgq5R/1WDJ9d1Cwun9WGRRQ1TLvzF1yABUerjjKrk89RCzYISwsKcgJPscaDqZgO6RIruY/xjuTtrnZSv+FXs+Woxf87P+QgQd76LC0MstTnys+AfTMuMPOLy9fMfEzs3LP0Nz6v5yjhX8ff7+3UUI3IcMxCvyxdTPClY5IvFdW7CCmmLNzakmx5GCItBWg/EIg1K1SG0jU9F8vlNZUqLKz42hWy/xB5C4QYQQ9ILdu4araPnrXnmd1D1QKVwKQ1DpWhNbpBDfE776/4xXD/tGM5O0TImp1NXul8wYsDi8g+e0pxNgY3Pahnj1yfG75Yw82spZanUH0QSNoMVMWnmV2hXGsWqypRq0pH8mPeLzeKa82gzsAZsouRD1k8wFlYA4z9HQFxqfcntTqXuwQcQIDAQABo2AwXjAdBgNVHQ4EFgQUfaEyaBpGNzsqttiSMETq+X/GJ0YwHwYDVR0jBBgwFoAUfaEyaBpGNzsqttiSMETq+X/GJ0YwCwYDVR0PBAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggIBADH5izxu4V8Javal8EA4DxZxIHUsWCg5cuopB28PsyJYpyKipsBoI8+RXqbtrLLue4WQfNPZHLXlKi+A3GTrLdlnenYzXVipPd+n3vRZyofaB3Jtb03nirVWGa8FG21Xy/f4rPqwcW54lxrnnh0SA0hwuZ+b2yAWESBXPxrzVQdTWCqoFI6/aRnN8RyZn0LqRYoW7WDtKpLmfyvshBmmu4PCYSh/SYiFHgR9fsWzVcxdySDsmX8wXowuFfp8V9sFhD4TsebAaplaICOuLUgj+Yin5QzgB0F9Ci3Zh6oWwl64SL/OxxQLpzMWzr0lrWsQrS3PgC4+6JC4IpTXX5eUqfSvHPtbRKK0yLnd9hYgvZUBvvZvUFR/3/fW+mpBHbZJBu9+/1uux46M4rJ2FeaJUf9PhYCPuUj63yu0Grn0DreVKK1SkD5V6qXN0TmoxYyguhfsIPCpI1VsdaSWuNjJ+a/HIlKIU8vKp5iN/+6ZTPAg9Q7s3Ji+vfx/AhFtQyTpIYNszVzNZyobvkiMUlK+eUKGlHVQp73y6MmGIlbBbyzpEoedNU4uFu57mw4fYGHqYZmYqFaiNQv4tVrGkg6p+Ypyu1zOfIHF7eqlAOu/SyRTvZkt9VtSVEOVH7nDIGdrCC9U/g1Lqk8Td00Oj8xesyKzsG214Xd8m7/7GmJ7nXe5", //nolint:lll
		UDPLines:      []string{"fast-io"},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
