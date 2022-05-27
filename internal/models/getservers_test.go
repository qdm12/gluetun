package models

import (
	"net"
	"testing"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AllServers_GetCopy(t *testing.T) {
	allServers := AllServers{
		Version: 1,
		ProviderToServers: map[string]Servers{
			providers.Cyberghost: {
				Version: 2,
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Expressvpn: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Fastestvpn: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.HideMyAss: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Ipvanish: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Ivpn: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Mullvad: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Nordvpn: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Perfectprivacy: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Privado: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.PrivateInternetAccess: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Privatevpn: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Protonvpn: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Purevpn: {
				Version: 1,
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Surfshark: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Torguard: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.VPNUnlimited: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Vyprvpn: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Wevpn: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
			providers.Windscribe: {
				Servers: []Server{{
					IPs: []net.IP{{1, 2, 3, 4}},
				}},
			},
		},
	}

	servers := allServers.GetCopy()

	assert.Equal(t, allServers, servers)
}

func Test_copyIPs(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		toCopy []net.IP
		copied []net.IP
	}{
		"nil": {},
		"empty": {
			toCopy: []net.IP{},
			copied: []net.IP{},
		},
		"single IP": {
			toCopy: []net.IP{{1, 1, 1, 1}},
			copied: []net.IP{{1, 1, 1, 1}},
		},
		"two IPs": {
			toCopy: []net.IP{{1, 1, 1, 1}, {2, 2, 2, 2}},
			copied: []net.IP{{1, 1, 1, 1}, {2, 2, 2, 2}},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Reserver leading 9 for copy modifications below
			for _, ipToCopy := range testCase.toCopy {
				require.NotEqual(t, 9, ipToCopy[0])
			}

			copied := copyIPs(testCase.toCopy)

			assert.Equal(t, testCase.copied, copied)

			if len(copied) > 0 {
				original := testCase.toCopy[0][0]
				testCase.toCopy[0][0] = 9
				assert.NotEqual(t, 9, copied[0][0])
				testCase.toCopy[0][0] = original

				copied[0][0] = 9
				assert.NotEqual(t, 9, testCase.toCopy[0][0])
			}
		})
	}
}
