package models

import (
	"testing"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/stretchr/testify/assert"
)

func Test_Servers_ToMarkdown(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		provider         string
		servers          Servers
		expectedMarkdown string
	}{
		providers.Cyberghost: {
			provider: providers.Cyberghost,
			servers: Servers{
				Servers: []Server{
					{Country: "a", UDP: true, Hostname: "xa"},
					{Country: "b", TCP: true, Hostname: "xb"},
				},
			},
			expectedMarkdown: "| Country | Hostname | TCP | UDP |\n" +
				"| --- | --- | --- | --- |\n" +
				"| a | `xa` | ❌ | ✅ |\n" +
				"| b | `xb` | ✅ | ❌ |\n",
		},
		providers.Fastestvpn: {
			provider: providers.Fastestvpn,
			servers: Servers{
				Servers: []Server{
					{Country: "a", Hostname: "xa", VPN: vpn.OpenVPN, TCP: true},
					{Country: "b", Hostname: "xb", VPN: vpn.OpenVPN, UDP: true},
				},
			},
			expectedMarkdown: "| Country | Hostname | VPN | TCP | UDP |\n" +
				"| --- | --- | --- | --- | --- |\n" +
				"| a | `xa` | openvpn | ✅ | ❌ |\n" +
				"| b | `xb` | openvpn | ❌ | ✅ |\n",
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			markdown := testCase.servers.ToMarkdown(testCase.provider)

			assert.Equal(t, testCase.expectedMarkdown, markdown)
		})
	}
}
