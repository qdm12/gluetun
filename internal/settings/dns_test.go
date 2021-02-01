package settings

import (
	"net"
	"testing"
	"time"

	"github.com/qdm12/dns/pkg/models"
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
				" |--DNS:",
				"     |--Plaintext address: 1.1.1.1",
				"     |--Keep nameserver (disabled blocking): no",
				"     |--DNS over TLS: disabled",
			},
		},
		"enabled DOT": {
			settings: DNS{
				Enabled:        true,
				KeepNameserver: true,
				Unbound: models.Settings{
					Providers: []string{"cloudflare"},
				},
				BlockMalicious:    true,
				BlockAds:          true,
				BlockSurveillance: true,
				UpdatePeriod:      time.Hour,
			},
			lines: []string{
				" |--DNS:",
				"     |--Keep nameserver (disabled blocking): yes",
				"     |--DNS over TLS:",
				"         |--Unbound:",
				"             |--DNS over TLS providers:",
				"                 |--cloudflare",
				"             |--Listening port: 0",
				"             |--Access control:",
				"                 |--Allowed:",
				"             |--Caching: disabled",
				"             |--IPv4 resolution: disabled",
				"             |--IPv6 resolution: disabled",
				"             |--Verbosity level: 0/5",
				"             |--Verbosity details level: 0/4",
				"             |--Validation log level: 0/2",
				"             |--Blocked hostnames:",
				"             |--Blocked IP addresses:",
				"             |--Allowed hostnames:",
				"         |--Block malicious: enabled",
				"         |--Block ads: enabled",
				"         |--Block surveillance: enabled",
				"         |--Update: every 1h0m0s",
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
