package utils

import (
	"testing"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/stretchr/testify/assert"
)

func Test_GetProtocol(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		selection configuration.ServerSelection
		protocol  string
	}{
		"default": {
			protocol: constants.UDP,
		},
		"OpenVPN UDP": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
			},
			protocol: constants.UDP,
		},
		"OpenVPN TCP": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
				OpenVPN: configuration.OpenVPNSelection{
					TCP: true,
				},
			},
			protocol: constants.TCP,
		},
		"Wireguard": {
			selection: configuration.ServerSelection{
				VPN: constants.Wireguard,
			},
			protocol: constants.UDP,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			protocol := GetProtocol(testCase.selection)

			assert.Equal(t, testCase.protocol, protocol)
		})
	}
}
