package privado

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
		Auth:           openvpn.SHA256,
		Ping:           10,
		TLSCipher:      "TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-DHE-DSS-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA",
		VerifyX509Type: "name",
		CA:             "MIIFKDCCAxCgAwIBAgIJAMtrmqZxIV/OMA0GCSqGSIb3DQEBDQUAMBIxEDAOBgNVBAMMB1ByaXZhZG8wHhcNMjAwMTA4MjEyODQ1WhcNMzUwMTA5MjEyODQ1WjASMRAwDgYDVQQDDAdQcml2YWRvMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAxPwOgiwNJzZTnKIXwAB0TSu/Lu2qt2U2I8obtQjwhi/7OrfmbmYykSdro70al2XPhnwAGGdCxW6LDnp0UN/IOhD11mgBPo14f5CLkBQjSJ6VN5miPbvK746LsNZl9H8rQGvDuPo4CG9BfPZMiDRGlsMxij/jztzgT1gmuxQ7WHfFRcNzBas1dHa9hV/d3TU6/t47x4SE/ljdcCtJiu7Zn6ODKQoys3mB7Luz2ngqUJWvkqsg+E4+3eJ0M8Hlbn5TPaRJBID7DAdYo6Vs6xGCYr981ThFcmoIQ10js10yANrrfGAzd03b3TnLAgko0uQMHjliMZL6L8sWOPHxyxJI0us88SFh4UgcFyRHKHPKux7w24SxAlZUYoUcTHp9VjG5XvDKYxzgV2RdM4ulBGbQRQ3y3/CyddsyQYMvA55Ets0LfPaBvDIcct70iXijGsdvlX1du3ArGpG7Vaje/RU4nbbGT6HYRdt5YyZfof288ukMOSj20nVcmS+c/4tqsxSerRb1aq5LOi1IemSkTMeC5gCbexk+L1vl7NT/58sxjGmu5bXwnvev/lIItfi2AlITrfUSEv19iDMKkeshwn/+sFJBMWYyluP+yJ56yR+MWoXvLlSWphLDTqq19yx3BZn0P1tgbXoR0g8PTdJFcz8z3RIb7myVLYulV1oGG/3rka0CAwEAAaOBgDB+MB0GA1UdDgQWBBTFtJkZCVDuDAD6k5bJzefjJdO3DTBCBgNVHSMEOzA5gBTFtJkZCVDuDAD6k5bJzefjJdO3DaEWpBQwEjEQMA4GA1UEAwwHUHJpdmFkb4IJAMtrmqZxIV/OMAwGA1UdEwQFMAMBAf8wCwYDVR0PBAQDAgEGMA0GCSqGSIb3DQEBDQUAA4ICAQB7MUSXMeBb9wlSv4sUaT1JHEwE26nlBw+TKmezfuPU5pBlY0LYr6qQZY95DHqsRJ7ByUzGUrGo17dNGXlcuNc6TAaQQEDRPo6y+LVh2TWMk15TUMI+MkqryJtCret7xGvDigKYMJgBy58HN3RAVr1B7cL9youwzLgc2Y/NcFKvnQJKeiIYAJ7g0CcnJiQvgZTS7xdwkEBXfsngmUCIG320DLPEL+Ze0HiUrxwWljMRya6i40AeH3Zu2i532xX1wV5+cjA4RJWIKg6ri/Q54iFGtZrA9/nc6y9uoQHkmz8cGyVUmJxFzMrrIICVqUtVRxLhkTMe4UzwRWTBeGgtW4tS0yq1QonAKfOyjgRw/CeY55D2UGvnAFZdTadtYXS4Alu2P9zdwoEk3fzHiVmDjqfJVr5wz9383aABUFrPI3nz6ed/Z6LZflKh1k+DUDEp8NxU4klUULWsSOKoa5zGX51G8cdHxwQLImXvtGuN5eSR8jCTgxFZhdps/xes4KkyfIz9FMYG748M+uOTgKITf4zdJ9BAyiQaOufVQZ8WjhWzWk9YHec9VqPkzpWNGkVjiRI5ewuXwZzZ164tMv2hikBXSuUCnFz37/ZNwGlDi0oBdDszCk2GxccdFHHaCSmpjU5MrdJ+5IhtTKGeTx+US2hTIVHQFIO99DmacxSYvLNcSQ==", //nolint:lll
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
