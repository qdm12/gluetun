package utils

import (
	"errors"
	"math/rand"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/stretchr/testify/assert"
)

func Test_GetConnection(t *testing.T) {
	t.Parallel()

	errTest := errors.New("test error")

	testCases := map[string]struct {
		provider        string
		filteredServers []models.Server
		filterError     error
		serverSelection settings.ServerSelection
		defaults        ConnectionDefaults
		ipv6Supported   bool
		randSource      rand.Source
		connection      models.Connection
		errWrapped      error
		errMessage      string
	}{
		"storage filter error": {
			filterError: errTest,
			errWrapped:  errTest,
			errMessage:  "cannot filter servers: test error",
		},
		"server without IPs": {
			filteredServers: []models.Server{
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
			filteredServers: []models.Server{
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
			filteredServers: []models.Server{
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
			filteredServers: []models.Server{
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
		"server with IPv4 and IPv6 and ipv6 supported": {
			filteredServers: []models.Server{
				{
					VPN: vpn.OpenVPN,
					UDP: true,
					IPs: []net.IP{
						net.IPv6zero,
						net.IPv4(1, 1, 1, 1),
					},
				},
			},
			serverSelection: settings.ServerSelection{}.
				WithDefaults(providers.Mullvad),
			defaults:      NewConnectionDefaults(443, 1194, 58820),
			ipv6Supported: true,
			randSource:    rand.NewSource(0),
			connection: models.Connection{
				Type:     vpn.OpenVPN,
				IP:       net.IPv6zero,
				Protocol: constants.UDP,
				Port:     1194,
			},
		},
		"mixed servers": {
			filteredServers: []models.Server{
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
			ctrl := gomock.NewController(t)

			storage := common.NewMockStorage(ctrl)
			storage.EXPECT().
				FilterServers(testCase.provider, testCase.serverSelection).
				Return(testCase.filteredServers, testCase.filterError)

			connection, err := GetConnection(testCase.provider, storage,
				testCase.serverSelection, testCase.defaults, testCase.ipv6Supported,
				testCase.randSource)

			assert.Equal(t, testCase.connection, connection)
			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
		})
	}
}
