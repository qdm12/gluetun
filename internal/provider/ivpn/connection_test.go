package ivpn

import (
	"errors"
	"math/rand"
	"net"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Ivpn_GetConnection(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers    []models.IvpnServer
		selection  configuration.ServerSelection
		connection models.Connection
		err        error
	}{
		"no server available": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
			},
			err: errors.New("no server found: for VPN openvpn; protocol udp"),
		},
		"no filter": {
			servers: []models.IvpnServer{
				{IPs: []net.IP{net.IPv4(1, 1, 1, 1)}, UDP: true},
				{IPs: []net.IP{net.IPv4(2, 2, 2, 2)}, UDP: true},
				{IPs: []net.IP{net.IPv4(3, 3, 3, 3)}, UDP: true},
			},
			connection: models.Connection{
				IP:       net.IPv4(1, 1, 1, 1),
				Port:     1194,
				Protocol: constants.UDP,
			},
		},
		"target IP": {
			selection: configuration.ServerSelection{
				TargetIP: net.IPv4(2, 2, 2, 2),
			},
			servers: []models.IvpnServer{
				{IPs: []net.IP{net.IPv4(1, 1, 1, 1)}, UDP: true},
				{IPs: []net.IP{net.IPv4(2, 2, 2, 2)}, UDP: true},
				{IPs: []net.IP{net.IPv4(3, 3, 3, 3)}, UDP: true},
			},
			connection: models.Connection{
				IP:       net.IPv4(2, 2, 2, 2),
				Port:     1194,
				Protocol: constants.UDP,
			},
		},
		"with filter": {
			selection: configuration.ServerSelection{
				Hostnames: []string{"b"},
			},
			servers: []models.IvpnServer{
				{Hostname: "a", IPs: []net.IP{net.IPv4(1, 1, 1, 1)}, UDP: true},
				{Hostname: "b", IPs: []net.IP{net.IPv4(2, 2, 2, 2)}, UDP: true},
				{Hostname: "a", IPs: []net.IP{net.IPv4(3, 3, 3, 3)}, UDP: true},
			},
			connection: models.Connection{
				IP:       net.IPv4(2, 2, 2, 2),
				Port:     1194,
				Protocol: constants.UDP,
				Hostname: "b",
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			randSource := rand.NewSource(0)

			m := New(testCase.servers, randSource)

			connection, err := m.GetConnection(testCase.selection)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.connection, connection)
		})
	}
}

func Test_getPort(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		selection configuration.ServerSelection
		port      uint16
	}{
		"default": {
			port: 1194,
		},
		"OpenVPN UDP": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
			},
			port: 1194,
		},
		"OpenVPN TCP": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
				OpenVPN: configuration.OpenVPNSelection{
					TCP: true,
				},
			},
			port: 443,
		},
		"OpenVPN custom port": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
				OpenVPN: configuration.OpenVPNSelection{
					CustomPort: 1234,
				},
			},
			port: 1234,
		},
		"Wireguard": {
			selection: configuration.ServerSelection{
				VPN: constants.Wireguard,
			},
			port: 58237,
		},
		"Wireguard custom port": {
			selection: configuration.ServerSelection{
				VPN: constants.Wireguard,
				Wireguard: configuration.WireguardSelection{
					CustomPort: 1234,
				},
			},
			port: 1234,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			port := getPort(testCase.selection)

			assert.Equal(t, testCase.port, port)
		})
	}
}

func Test_getProtocol(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		selection configuration.ServerSelection
		protocol  string
	}{
		"OpenVPN UDP": {
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

			protocol := getProtocol(testCase.selection)

			assert.Equal(t, testCase.protocol, protocol)
		})
	}
}
