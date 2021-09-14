package models

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AllServers_GetCopy(t *testing.T) {
	allServers := AllServers{
		Cyberghost: CyberghostServers{
			Version: 2,
			Servers: []CyberghostServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Expressvpn: ExpressvpnServers{
			Servers: []ExpressvpnServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Fastestvpn: FastestvpnServers{
			Servers: []FastestvpnServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		HideMyAss: HideMyAssServers{
			Servers: []HideMyAssServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Ipvanish: IpvanishServers{
			Servers: []IpvanishServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Ivpn: IvpnServers{
			Servers: []IvpnServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Mullvad: MullvadServers{
			Servers: []MullvadServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Nordvpn: NordvpnServers{
			Servers: []NordvpnServer{{
				IP: net.IP{1, 2, 3, 4},
			}},
		},
		Privado: PrivadoServers{
			Servers: []PrivadoServer{{
				IP: net.IP{1, 2, 3, 4},
			}},
		},
		Pia: PiaServers{
			Servers: []PIAServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Privatevpn: PrivatevpnServers{
			Servers: []PrivatevpnServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Protonvpn: ProtonvpnServers{
			Servers: []ProtonvpnServer{{
				EntryIP: net.IP{1, 2, 3, 4},
				ExitIP:  net.IP{1, 2, 3, 4},
			}},
		},
		Purevpn: PurevpnServers{
			Version: 1,
			Servers: []PurevpnServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Surfshark: SurfsharkServers{
			Servers: []SurfsharkServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Torguard: TorguardServers{
			Servers: []TorguardServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		VPNUnlimited: VPNUnlimitedServers{
			Servers: []VPNUnlimitedServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Vyprvpn: VyprvpnServers{
			Servers: []VyprvpnServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
		Windscribe: WindscribeServers{
			Servers: []WindscribeServer{{
				IPs: []net.IP{{1, 2, 3, 4}},
			}},
		},
	}

	servers := allServers.GetCopy()

	assert.Equal(t, allServers, servers)
}

func Test_AllServers_GetVyprvpn(t *testing.T) {
	allServers := AllServers{
		Vyprvpn: VyprvpnServers{
			Servers: []VyprvpnServer{
				{Hostname: "a", IPs: []net.IP{{1, 1, 1, 1}}},
				{Hostname: "b", IPs: []net.IP{{2, 2, 2, 2}}},
			},
		},
	}

	servers := allServers.GetVyprvpn()

	expectedServers := []VyprvpnServer{
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
