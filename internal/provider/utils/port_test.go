package utils

import (
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/stretchr/testify/assert"
)

func boolPtr(b bool) *bool       { return &b }
func uint16Ptr(n uint16) *uint16 { return &n }

func Test_GetPort(t *testing.T) {
	t.Parallel()

	const (
		defaultOpenVPNTCP = 443
		defaultOpenVPNUDP = 1194
		defaultWireguard  = 51820
	)

	testCases := map[string]struct {
		selection settings.ServerSelection
		port      uint16
	}{
		"default": {
			port: defaultOpenVPNUDP,
		},
		"OpenVPN UDP": {
			selection: settings.ServerSelection{
				VPN: constants.OpenVPN,
			},
			port: defaultOpenVPNUDP,
		},
		"OpenVPN TCP": {
			selection: settings.ServerSelection{
				VPN: constants.OpenVPN,
				OpenVPN: settings.OpenVPNSelection{
					TCP: boolPtr(true),
				},
			},
			port: defaultOpenVPNTCP,
		},
		"OpenVPN custom port": {
			selection: settings.ServerSelection{
				VPN: constants.OpenVPN,
				OpenVPN: settings.OpenVPNSelection{
					CustomPort: uint16Ptr(1234),
				},
			},
			port: 1234,
		},
		"Wireguard": {
			selection: settings.ServerSelection{
				VPN: constants.Wireguard,
			},
			port: defaultWireguard,
		},
		"Wireguard custom port": {
			selection: settings.ServerSelection{
				VPN: constants.Wireguard,
				Wireguard: settings.WireguardSelection{
					EndpointPort: uint16Ptr(1234),
				},
			},
			port: 1234,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			port := GetPort(testCase.selection,
				defaultOpenVPNTCP, defaultOpenVPNUDP, defaultWireguard)

			assert.Equal(t, testCase.port, port)
		})
	}
}
