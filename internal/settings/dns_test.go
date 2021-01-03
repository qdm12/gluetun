package settings

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DNS_Lines(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		settings DNS
		lines    []string
	}{
		"disabled DOT": {
			settings: DNS{
				PlaintextAddress: net.IP{1, 1, 1, 1},
			},
			lines: []string{
				"DNS over TLS disabled, using plaintext DNS 1.1.1.1",
			},
		},
		"enabled DOT": {
			settings: DNS{
				Enabled: true,
			},
			lines: []string{
				"DNS settings:",
				" |--Unbound:",
				"    |--DNS over TLS provider:",
				"    |--Listening port: 0",
				"    |--Access control:",
				"       |--Allowed:",
				"    |--Caching: disabled",
				"    |--IPv4 resolution: disabled",
				"    |--IPv6 resolution: disabled",
				"    |--Verbosity level: 0/5",
				"    |--Verbosity details level: 0/4",
				"    |--Validation log level: 0/2",
				"    |--Blocked hostnames:",
				"    |--Blocked IP addresses:",
				"    |--Allowed hostnames:",
				" |--Block malicious: disabled",
				" |--Block ads: disabled",
				" |--Block surveillance: disabled",
				" |--Update: deactivated",
				" |--Keep nameserver (disabled blocking): no",
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
