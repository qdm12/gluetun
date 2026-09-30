package settings

import (
	"testing"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/stretchr/testify/assert"
)

func Test_OpenVPNSelection_validate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		protocol   string
		customPort uint16
		errMessage string
	}{
		"default_udp_no_custom_port": {},
		"tcp_no_custom_port": {
			protocol: constants.TCP,
		},
		"udp_custom_port_allowed": {
			customPort: 15021,
		},
		"tcp_custom_port_allowed": {
			protocol:   constants.TCP,
			customPort: 443,
		},
		"udp_custom_port_denied": {
			customPort: 5565,
			errMessage: "custom endpoint port is not allowed: for VPN service provider purevpn",
		},
		"tcp_custom_port_denied": {
			protocol:   constants.TCP,
			customPort: 8080,
			errMessage: "custom endpoint port is not allowed: for VPN service provider purevpn",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			selection := OpenVPNSelection{}
			selection.setDefaults(providers.Purevpn)
			if testCase.protocol != "" {
				selection.Protocol = testCase.protocol
			}
			selection.CustomPort = new(testCase.customPort)

			err := selection.validate(providers.Purevpn)

			if testCase.errMessage == "" {
				assert.NoError(t, err)
				return
			}
			assert.ErrorContains(t, err, testCase.errMessage)
		})
	}
}
