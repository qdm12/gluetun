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
				OpenVPN: OpenVPN{
					Provider: Provider{
						Name: constants.Mullvad,
					},
				},
			},
			lines: []string{
				"Settings summary below:",
				"|--OpenVPN:",
				"   |--Verbosity level: 0",
				"   |--Provider:",
				"      |--Mullvad settings:",
				"         |--Network protocol: ",
				"|--DNS:",
				"|--Firewall: disabled ⚠️",
				"|--System:",
				"   |--Process user ID: 0",
				"   |--Process group ID: 0",
				"   |--Timezone: NOT SET ⚠️ - it can cause time related issues",
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
