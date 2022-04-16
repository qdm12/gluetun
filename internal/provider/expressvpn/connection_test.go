package expressvpn

import (
	"errors"
	"math/rand"
	"net"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Provider_GetConnection(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers    []models.Server
		selection  settings.ServerSelection
		connection models.Connection
		err        error
	}{
		"no server available": {
			selection: settings.ServerSelection{}.WithDefaults(providers.Expressvpn),
			err:       errors.New("no server found: for VPN openvpn; protocol udp"),
		},
		"no filter": {
			servers: []models.Server{
				{IPs: []net.IP{net.IPv4(1, 1, 1, 1)}, UDP: true},
				{IPs: []net.IP{net.IPv4(2, 2, 2, 2)}, UDP: true},
				{IPs: []net.IP{net.IPv4(3, 3, 3, 3)}, UDP: true},
			},
			selection: settings.ServerSelection{}.WithDefaults(providers.Expressvpn),
			connection: models.Connection{
				Type:     constants.OpenVPN,
				IP:       net.IPv4(1, 1, 1, 1),
				Port:     1195,
				Protocol: constants.UDP,
			},
		},
		"target IP": {
			selection: settings.ServerSelection{
				TargetIP: net.IPv4(2, 2, 2, 2),
			}.WithDefaults(providers.Expressvpn),
			servers: []models.Server{
				{IPs: []net.IP{net.IPv4(1, 1, 1, 1)}, UDP: true},
				{IPs: []net.IP{net.IPv4(2, 2, 2, 2)}, UDP: true},
				{IPs: []net.IP{net.IPv4(3, 3, 3, 3)}, UDP: true},
			},
			connection: models.Connection{
				Type:     constants.OpenVPN,
				IP:       net.IPv4(2, 2, 2, 2),
				Port:     1195,
				Protocol: constants.UDP,
			},
		},
		"with filter": {
			selection: settings.ServerSelection{
				Hostnames: []string{"b"},
			}.WithDefaults(providers.Expressvpn),
			servers: []models.Server{
				{Hostname: "a", IPs: []net.IP{net.IPv4(1, 1, 1, 1)}, UDP: true},
				{Hostname: "b", IPs: []net.IP{net.IPv4(2, 2, 2, 2)}, UDP: true},
				{Hostname: "a", IPs: []net.IP{net.IPv4(3, 3, 3, 3)}, UDP: true},
			},
			connection: models.Connection{
				Type:     constants.OpenVPN,
				IP:       net.IPv4(2, 2, 2, 2),
				Port:     1195,
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
