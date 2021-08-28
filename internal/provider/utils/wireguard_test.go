package utils

import (
	"net"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/wireguard"
	"github.com/stretchr/testify/assert"
)

func Test_BuildWireguardSettings(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		connection   models.Connection
		userSettings configuration.Wireguard
		settings     wireguard.Settings
	}{
		"some settings": {
			connection: models.Connection{
				IP:     net.IPv4(1, 2, 3, 4),
				Port:   51821,
				PubKey: "public",
			},
			userSettings: configuration.Wireguard{
				PrivateKey:   "private",
				PreSharedKey: "pre-shared",
				Addresses: []*net.IPNet{
					{IP: net.IPv4(1, 1, 1, 1), Mask: net.IPv4Mask(255, 255, 255, 255)},
					{IP: net.IPv4(2, 2, 2, 2), Mask: net.IPv4Mask(255, 255, 255, 255)},
				},
				Interface: "wg1",
			},
			settings: wireguard.Settings{
				InterfaceName: "wg1",
				PrivateKey:    "private",
				PublicKey:     "public",
				PreSharedKey:  "pre-shared",
				Endpoint: &net.UDPAddr{
					IP:   net.IPv4(1, 2, 3, 4),
					Port: 51821,
				},
				Addresses: []*net.IPNet{
					{IP: net.IPv4(1, 1, 1, 1), Mask: net.IPv4Mask(255, 255, 255, 255)},
					{IP: net.IPv4(2, 2, 2, 2), Mask: net.IPv4Mask(255, 255, 255, 255)},
				},
				RulePriority: 101,
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			settings := BuildWireguardSettings(testCase.connection,
				testCase.userSettings)

			assert.Equal(t, testCase.settings, settings)
		})
	}
}
