package fastestvpn

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
		},
		Auth:          openvpn.SHA256,
		MssFix:        1450,
		TLSCipher:     "TLS-DHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA256:TLS-DHE-RSA-WITH-CAMELLIA-256-CBC-SHA:TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-CAMELLIA-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA", //nolint:lll
		AuthToken:     true,
		KeyDirection:  "1",
		RenegDisabled: true,
		CA:            "MIIFQjCCAyqgAwIBAgIIUfxepT+rr8owDQYJKoZIhvcNAQEMBQAwPzELMAkGA1UEBhMCS1kxEzARBgNVBAoTCkZhc3Rlc3RWUE4xGzAZBgNVBAMTEkZhc3Rlc3RWUE4gUm9vdCBDQTAeFw0xNzA5MTYwMDAxNDZaFw0yNzA5MTQwMDAxNDZaMD8xCzAJBgNVBAYTAktZMRMwEQYDVQQKEwpGYXN0ZXN0VlBOMRswGQYDVQQDExJGYXN0ZXN0VlBOIFJvb3QgQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC1Xj+WfPTozFynFqc+c3CVrggIllaXEl5bY5VgFynXkqCTM6lSrfC4pNjGXUbqWe6RnGJbM4/6kUn+lQDjFSQV1rzP2eDS8+r5+X2WXh4AoeNRUWhvSG+HiHD/B2EFK+Nd5BRSdUjpKWAtsCmT2bBt7nT0jN1OdeNrLJeyF8siAqv/oQzKznF9aIe/N01b2M8ZOFTzoXi2fZAckgGWui8NB/lzkVIJqSkAPRL8qiJLuRCPVOX1PFD8vV//R8/QumtfbcYBMo6vCk2HmWdrh5OQHPxb3KJtbtG+Z1j8x6HGEAe17djYepBiRMyCEQvYgfD6tvFylc4IquhqE9yaP60PJod5TxpWnRQ6HIGSeBm+S+rYSMalTZ8+pUqOOA+IQCYpfpx6EKIJL/VsW2C7cXdvudxDhXPI5lR/QidCb9Ohq3WkfxXaYwzrngdg2avmNqId9R4KESuM9GoHW0dszfyBCh5wYfeaffMElfDam3B92NUwyhZwtIiv623WVXY9PPz+EDjSJsIAu2Vi1vdJyA4nD4k9Lwmx/1zTc/UaYVLsiBqL2WdfvFTeoWoV+dNxQXSEPhB8gwi8x4O4lZW0cwVy/6fa8KMY8gZbcbSTr7U5bRERfW8l+jY+mYKQ/M/ccgpxaHiw1/+4LWfbJQ7VhJJrTyN0C36FQzY1URkSXg+53wIDAQABo0IwQDAPBgNVHRMBAf8EBTADAQH/MA4GA1UdDwEB/wQEAwIBBjAdBgNVHQ4EFgQUmVEL4x6xdCqiqu2OBLs27EA8xGYwDQYJKoZIhvcNAQEMBQADggIBABCpITvO1+R4T9v2+onHiFxU5JjtCZ0zkXqRCMp/Z0UIYbeo1p07pZCPAUjBfGPCkAaR++OiG9sysALdJf8Y6HQKcyuAcWUqQnaIhoZ2JcAP7EKq7uCqsMhcYZD/j3O/3RPtSW5UOx6ItDU+Ua0t9Edho9whNw0VQXmo1JjYoP3FzPjuKoDWTSO1q5eYlZfwcTcs55O2shNkFafPg/6cCm5j6v9nyHrM3sk4LjkrBPUXVx2m/aoz219t8O9Ha9/CdMKXsPO/8gTUzpgnzSgPnGnBmi5xr1nspVN8X4E2f3D+DKqBim3YgslD68NcuFQvJ0/BxZzWVbrr+QXoyzaiCgXuogpIDc2bB6oRXqFnHNz36d4QJmJdWdSaijiS/peQ6EOPgOZ1GuObLWlDCBZLNeQ+N6QaiJxVO4XUj/s22i1IRtwdz84TRHrbWiIpEymsqmb/Ep5r4xV5d6+791axclfOTH7tQrY/SbPtTJI4OEgNekI8YfadQifpelF82MsFFEZuaQn0lj+fvLGtE/zKh3OdLTxRc5TAgBB+0T81+JQosygNr2aFFG0hxar1eyw/gLeG8H+7Ie50pyPvXO4OgB6Key8rSExpilQXlvAT1qX0qS3/K1i/9QkSE9ftIPT6vtwLV2sVQzfyanI4IZgWC6ryhvNLsRn0NFnQclor0+aq", //nolint:lll //nolint:lll
		TLSAuth:       "697fe793b32cb5091d30f2326d5d124a9412e93d0a44ef7361395d76528fcbfc82c3859dccea70a93cfa8fae409709bff75f844cf5ff0c237f426d0c20969233db0e706edb6bdf195ec3dc11b3f76bc807a77e74662d9a800c8cd1144ebb67b7f0d3f1281d1baf522bfe03b7c3f963b1364fc0769400e413b61ca7b43ab19fac9e0f77e41efd4bda7fd77b1de2d7d7855cbbe3e620cecceac72c21a825b243e651f44d90e290e09c3ad650de8fca99c858bc7caad584bc69b11e5c9fd9381c69c505ec487a65912c672d83ed0113b5a74ddfbd3ab33b3683cec593557520a72c4d6cce46111f56f3396cc3ce7183edce553c68ea0796cf6c4375fad00aaa2a42",                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         //nolint:lll
		UDPLines: []string{
			"tun-mtu 1500",
			"tun-mtu-extra 32",
			"ping 15",
		},
		ExtraLines: []string{
			"comp-lzo",
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
