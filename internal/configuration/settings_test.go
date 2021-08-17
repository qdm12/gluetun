package configuration

import (
	"testing"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/stretchr/testify/assert"
)

func Test_Settings_lines(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		settings Settings
		lines    []string
	}{
		"default settings": {
			settings: Settings{
				VPN: VPN{
					Type: constants.OpenVPN,
					Provider: Provider{
						Name: constants.Mullvad,
					},
					OpenVPN: OpenVPN{
						Version: constants.Openvpn25,
					},
				},
			},
			lines: []string{
				"Settings summary below:",
				"|--VPN:",
				"   |--Type: openvpn",
				"   |--OpenVPN:",
				"      |--Version: 2.5",
				"      |--Verbosity level: 0",
				"   |--Mullvad settings:",
				"      |--OpenVPN selection:",
				"         |--Protocol: udp",
				"|--DNS:",
				"|--Firewall: disabled ⚠️",
				"|--System:",
				"   |--Process user ID: 0",
				"   |--Process group ID: 0",
				"   |--Timezone: NOT SET ⚠️ - it can cause time related issues",
				"|--Health:",
				"   |--Server address: ",
				"   |--OpenVPN:",
				"      |--Initial duration: 0s",
				"|--HTTP control server:",
				"   |--Listening port: 0",
				"|--Public IP getter: disabled",
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			lines := testCase.settings.lines()

			assert.Equal(t, testCase.lines, lines)
		})
	}
}
