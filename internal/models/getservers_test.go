package models

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AllServers_GetCopy(t *testing.T) {
	allServers := AllServers{
		Cyberghost: Servers{
			Version: 2,
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Expressvpn: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Fastestvpn: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		HideMyAss: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Ipvanish: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Ivpn: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Mullvad: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Nordvpn: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Perfectprivacy: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Privado: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Pia: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Privatevpn: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Protonvpn: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Purevpn: Servers{
			Version: 1,
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Surfshark: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Torguard: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		VPNUnlimited: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Vyprvpn: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Windscribe: Servers{
			Servers: []Server{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
	}

	servers := allServers.GetCopy()

	assert.Equal(t, allServers, servers)
}

func Test_AllServers_GetVyprvpn(t *testing.T) {
	allServers := AllServers{
		Vyprvpn: Servers{
			Servers: []Server{
				{Hostname: "a", IPs: []net.IP{{1, 1, 1, 1}}},
				{Hostname: "b", IPs: []net.IP{{2, 2, 2, 2}}},
			},
		},
	}

	servers := allServers.GetVyprvpn()

	expectedServers := []Server{
		{Hostname: "a", IPs: []net.IP{{1, 1, 1, 1}}},
		{Hostname: "b", IPs: []net.IP{{2, 2, 2, 2}}},
	}
	assert.Equal(t, expectedServers, servers)

	allServers.Vyprvpn.Servers[0].IPs[0][0] = 9
	assert.NotEqual(t, 9, servers[0].IPs[0][0])

	allServers.Vyprvpn.Servers[0].IPs[0][0] = 1
	servers[0].IPs[0][0] = 9
	assert.NotEqual(t, 9, allServers.Vyprvpn.Servers[0].IPs[0][0])
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
