package privatevpn

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
			openvpn.AES128gcm,
		},
		Auth:    openvpn.SHA256,
		CA:      "MIIErTCCA5WgAwIBAgIJAPp3HmtYGCIOMA0GCSqGSIb3DQEBCwUAMIGVMQswCQYDVQQGEwJTRTELMAkGA1UECBMCQ0ExEjAQBgNVBAcTCVN0b2NraG9sbTETMBEGA1UEChMKUHJpdmF0ZVZQTjEWMBQGA1UEAxMNUHJpdmF0ZVZQTiBDQTETMBEGA1UEKRMKUHJpdmF0ZVZQTjEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEBwcml2YXR2cG4uc2UwHhcNMTcwNTI0MjAxNTM3WhcNMjcwNTIyMjAxNTM3WjCBlTELMAkGA1UEBhMCU0UxCzAJBgNVBAgTAkNBMRIwEAYDVQQHEwlTdG9ja2hvbG0xEzARBgNVBAoTClByaXZhdGVWUE4xFjAUBgNVBAMTDVByaXZhdGVWUE4gQ0ExEzARBgNVBCkTClByaXZhdGVWUE4xIzAhBgkqhkiG9w0BCQEWFHN1cHBvcnRAcHJpdmF0dnBuLnNlMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwjqTWbKk85WN8nd1TaBgBnBHceQWosp8mMHr4xWMTLagWRcq2Modfy7RPnBo9kyn5j/ZZwL/21gLWJbxidurGyZZdEV9Wb5KQl3DUNxa19kwAbkkEchdES61e99MjmQlWq4vGPXAHjEuDxOZ906AXglCyAvQoXcYW0mNm9yybWllVp1aBrCaZQrNYr7eoFvolqJXdQQ3FFsTBCYa5bHJcKQLBfsiqdJ/BAxhNkQtcmWNSgLy16qoxQpCsxNCxAcYnasuL4rwOP+RazBkJTPXA/2neCJC5rt+sXR9CSfiXdJGwMpYso5m31ZEd7JL2+is0FeAZ6ETrKMnEZMsTpTkdwIDAQABo4H9MIH6MB0GA1UdDgQWBBRCkBlC94zCY6VNncMnK36JxT7bazCBygYDVR0jBIHCMIG/gBRCkBlC94zCY6VNncMnK36JxT7ba6GBm6SBmDCBlTELMAkGA1UEBhMCU0UxCzAJBgNVBAgTAkNBMRIwEAYDVQQHEwlTdG9ja2hvbG0xEzARBgNVBAoTClByaXZhdGVWUE4xFjAUBgNVBAMTDVByaXZhdGVWUE4gQ0ExEzARBgNVBCkTClByaXZhdGVWUE4xIzAhBgkqhkiG9w0BCQEWFHN1cHBvcnRAcHJpdmF0dnBuLnNlggkA+ncea1gYIg4wDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAayugvExKDHar7t1zyYn99Vt1NMf46J8x4Dt9TNjBml5mR9nKvWmreMUuuOhLaO8Da466KGdXeDFNLcBYZd/J2iTawE6/3fmrML9H2sa+k/+E4uU5nQ84ZGOwCinCkMalVjM8EZ0/H2RZvLAVUnvPuUz2JfJhmiRkbeE75fVuqpAm9qdE+/7lg3oICYzxa6BJPxT+Imdjy3Q/FWdsXqX6aallhohPAZlMZgZL4eXECnV8rAfzyjOJggkMDZQt3Flc0Y4iDMfzrEhSOWMkNFBFwjK0F/dnhsX+fPX6GGRpUZgZcCt/hWvypqc05/SnrdKM/vV/jV/yZe0NVzY7S8Ur5g==", //nolint:lll
		TLSAuth: "f035a3acaeffb5aedb5bc920bca26ca7ac701da88249008e03563eba6af6d2625ac8ba1e5e0921f76be004c24ae4fd43e42caf0f84269ad44d8d4c14ba45b1386f251c7330d8cc56afd16d516835645651ef7e87a723ac78ae0d49da5b2f2d78ceafcff7a6367d0712628a6547e5fc8fef93c87f7bcd6107c7b1ae68396e944aadae50111d01a5d0c67223d667bdbf1bf434bdef03644ecc5386e102724eef3872f66547eb66dc0fea8286069cb082a41c89083b28fe9f4cec25d48017f26c4fd85b25ddf2ae5448dd2bccf3eef2aacf42ef1e88c3248c689423d0b05a641e9e79dd6b9b5c40f0cc21ffdc891b9eee951477b537261cb56a958a4f490d961ecb",                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     //nolint:lll
		UDPLines: []string{
			"key-direction 1",
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
