package utils

import (
	"math/rand"
	"net"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_GetConnection(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers         []models.Server
		serverSelection settings.ServerSelection
		defaults        ConnectionDefaults
		randSource      rand.Source
		connection      models.Connection
		errWrapped      error
		errMessage      string
	}{
		"no server": {
			serverSelection: settings.ServerSelection{}.
				WithDefaults(providers.Mullvad),
			errWrapped: ErrNoServer,
			errMessage: "no server",
		},
		"all servers filtered": {
			servers: []models.Server{
				{VPN: vpn.Wireguard},
				{VPN: vpn.Wireguard},
			},
			serverSelection: settings.ServerSelection{
				VPN: vpn.OpenVPN,
			}.WithDefaults(providers.Mullvad),
			errWrapped: ErrNoServerFound,
			errMessage: "no server found: for VPN openvpn; protocol udp",
		},
		"server without IPs": {
			servers: []models.Server{
				{VPN: vpn.OpenVPN, UDP: true},
				{VPN: vpn.OpenVPN, UDP: true},
			},
			serverSelection: settings.ServerSelection{}.
				WithDefaults(providers.Mullvad),
			defaults: ConnectionDefaults{
				OpenVPNTCPPort: 1,
				OpenVPNUDPPort: 1,
				WireguardPort:  1,
			},
			errWrapped: ErrNoConnectionToPickFrom,
			errMessage: "no connection to pick from",
		},
		"OpenVPN server with hostname": {
			servers: []models.Server{
				{
					VPN:      vpn.OpenVPN,
					UDP:      true,
					IPs:      []net.IP{net.IPv4(1, 1, 1, 1)},
					Hostname: "name",
				},
			},
			serverSelection: settings.ServerSelection{}.
				WithDefaults(providers.Mullvad),
			defaults:   NewConnectionDefaults(443, 1194, 58820),
			randSource: rand.NewSource(0),
			connection: models.Connection{
				Type:     vpn.OpenVPN,
				IP:       net.IPv4(1, 1, 1, 1),
				Protocol: constants.UDP,
				Port:     1194,
				Hostname: "name",
			},
		},
		"OpenVPN server with x509": {
			servers: []models.Server{
				{
					VPN:      vpn.OpenVPN,
					UDP:      true,
					IPs:      []net.IP{net.IPv4(1, 1, 1, 1)},
					Hostname: "hostname",
					OvpnX509: "x509",
				},
			},
			serverSelection: settings.ServerSelection{}.
				WithDefaults(providers.Mullvad),
			defaults:   NewConnectionDefaults(443, 1194, 58820),
			randSource: rand.NewSource(0),
			connection: models.Connection{
				Type:     vpn.OpenVPN,
				IP:       net.IPv4(1, 1, 1, 1),
				Protocol: constants.UDP,
				Port:     1194,
				Hostname: "x509",
			},
		},
		"server with IPv4 and IPv6": {
			servers: []models.Server{
				{
					VPN: vpn.OpenVPN,
					UDP: true,
					IPs: []net.IP{
						net.IPv4(1, 1, 1, 1),
						// All IPv6 is ignored
						net.IPv6zero,
						net.IPv6zero,
						net.IPv6zero,
						net.IPv6zero,
						net.IPv6zero,
					},
				},
			},
			serverSelection: settings.ServerSelection{}.
				WithDefaults(providers.Mullvad),
			defaults:   NewConnectionDefaults(443, 1194, 58820),
			randSource: rand.NewSource(0),
			connection: models.Connection{
				Type:     vpn.OpenVPN,
				IP:       net.IPv4(1, 1, 1, 1),
				Protocol: constants.UDP,
				Port:     1194,
			},
		},
		"mixed servers": {
			servers: []models.Server{
				{
					VPN:      vpn.OpenVPN,
					UDP:      true,
					IPs:      []net.IP{net.IPv4(1, 1, 1, 1)},
					OvpnX509: "ovpnx509",
				},
				{
					VPN:      vpn.Wireguard,
					UDP:      true,
					IPs:      []net.IP{net.IPv4(2, 2, 2, 2)},
					OvpnX509: "ovpnx509",
				},
				{
					VPN: vpn.OpenVPN,
					UDP: true,
					IPs: []net.IP{
						net.IPv4(3, 3, 3, 3),
						{1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, // ipv6 ignored
					},
					Hostname: "hostname",
				},
			},
			serverSelection: settings.ServerSelection{}.
				WithDefaults(providers.Mullvad),
			defaults:   NewConnectionDefaults(443, 1194, 58820),
			randSource: rand.NewSource(0),
			connection: models.Connection{
				Type:     vpn.OpenVPN,
				IP:       net.IPv4(1, 1, 1, 1),
				Protocol: constants.UDP,
				Port:     1194,
				Hostname: "ovpnx509",
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			connection, err := GetConnection(testCase.servers,
				testCase.serverSelection, testCase.defaults,
				testCase.randSource)

			assert.Equal(t, testCase.connection, connection)
			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
		})
	}
}
