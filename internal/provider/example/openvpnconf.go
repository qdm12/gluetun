package example

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/openvpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) OpenVPNConfig(connection models.Connection,
	settings settings.OpenVPN, ipv6Supported bool) (lines []string) {
	// TODO: Set the necessary fields in `providerSettings` to
	// generate the right OpenVPN configuration file.
	//nolint:gomnd
	providerSettings := utils.OpenVPNProviderSettings{
		AuthUserPass: true,
		Ciphers: []string{
			openvpn.AES256gcm,
		},
		Ping:          5,
		RemoteCertTLS: true,
		CA:            "MIIDZzCCAk+gAwIBAgIUVwHEFE6geihigDSNkBppm2Zamx0wDQYJKoZIhvcNAQELBQAwQzELMAkGA1UEBhMCQ0ExDzANBgNVBAgMBlF1ZWJlYzERMA8GA1UEBwwITW9udHJlYWwxEDAOBgNVBAoMB0dsdWV0dW4wHhcNMjIwNzAxMTY1MzE5WhcNMjcwNjMwMTY1MzE5WjBDMQswCQYDVQQGEwJDQTEPMA0GA1UECAwGUXVlYmVjMREwDwYDVQQHDAhNb250cmVhbDEQMA4GA1UECgwHR2x1ZXR1bjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALmJRhTUr+87NFkHL2PWjIz7efHqQgrWuDQt8oOBHvl0Hm72N+ckO+5Q0zG4XtqlpBjFjGUSjfNUWSrRztbXlMmzDcjHKjYHUPepJpoF100fK2q3XPiFRl6sEXzYeOdFgpaTdmGHS6DL9aWeCoYA/k6NV8YqHXujr14gOYOAWG6cRimpTJf8DtEDcxtp1w6fOEoN0b5PvV7dcpLiva8LYyZKPvFYBzlc5BZxOLvq6bvhQm54R6zoHFpaKOf7FeqhxI6KOQu4IPwU12YBlOP5CbkMAQ1cWWVQ4pnh0Hwh71Sjm848jS/OcammNzsp4xWaKt/pzcix3fpJt/MDP/9fxA8CAwEAAaNTMFEwHQYDVR0OBBYEFCIQ9l28Yy1/3qJvFarXjhKdG9tVMB8GA1UdIwQYMBaAFCIQ9l28Yy1/3qJvFarXjhKdG9tVMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAKLPmLTppXYTTOOxHhTMyHI0oTl7ID2PQfJsref+jDshB3hib98BC17b9ESpLnwx7ugg17NRl7RYutxjuVw/CK/gwAnTMg3D3mdAnKkMRr3UxnD89KprLIpf7WQCmyJaxalsD5jjgl3kuGM7jf2FJNxQz5RrXBGlQHa465ouov+Rp5v/K5Umyt6wsCZXEbOF0SdUhEGU3nxVbFsoPimNYSHHwc29USnQmyW1O/drFDoTcOK4GdHFEVkrHQgqwU8ay1fYGYfIUDhsV/5AAWgQC41r9FWr+VQgyJC94qmDg0c46RE123dL/YifVUl2DKuJ0aOY+OkSgwknKZ+FQd+8d6k=", //nolint:lll
		TLSAuth:       "bc470c93ff9f5602a8abb27dee84a52814d10f20490ad23c47d5d82120c1bf859e93d0696b455d4a1b8d55d40c2685c41ca1d0aef29a3efd27274c4ef09020a3978fe45784b335da6df2d12db97bbb838416515f2a96f04715fd28949c6fe296a925cfada3f8b8928ed7fc963c1563272f5cf46e5e1d9c845d7703ca881497b7e6564a9d1dea9358adffd435295479f47d5298fabf5359613ff5992cb57ff081a04dfb81a26513a6b44a9b5490ad265f8a02384832a59cc3e075ad545461060b7bcab49bac815163cb80983dd51d5b1fd76170ffd904d8291071e96efc3fb777856c717b148d08a510f5687b8a8285dcffe737b98916dd15ef6235dee4266d3b",                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 //nolint:lll
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
