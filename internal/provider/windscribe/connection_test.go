package windscribe

import (
	"math/rand"
	"net"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
	"github.com/stretchr/testify/assert"
)

func Test_Windscribe_GetConnection(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers    []models.Server
		selection  settings.ServerSelection
		connection models.Connection
		errWrapped error
		errMessage string
	}{
		"no server available": {
			selection:  settings.ServerSelection{}.WithDefaults(providers.Windscribe),
			errWrapped: utils.ErrNoServer,
			errMessage: "no server",
		},
		"no filter": {
			servers: []models.Server{
				{VPN: vpn.OpenVPN, UDP: true, IPs: []net.IP{net.IPv4(1, 1, 1, 1)}},
				{VPN: vpn.OpenVPN, UDP: true, IPs: []net.IP{net.IPv4(2, 2, 2, 2)}},
				{VPN: vpn.OpenVPN, UDP: true, IPs: []net.IP{net.IPv4(3, 3, 3, 3)}},
			},
			selection: settings.ServerSelection{}.WithDefaults(providers.Windscribe),
			connection: models.Connection{
				Type:     vpn.OpenVPN,
				IP:       net.IPv4(1, 1, 1, 1),
				Port:     1194,
				Protocol: constants.UDP,
			},
		},
		"target IP": {
			selection: settings.ServerSelection{
				TargetIP: net.IPv4(2, 2, 2, 2),
			}.WithDefaults(providers.Windscribe),
			servers: []models.Server{
				{IPs: []net.IP{net.IPv4(1, 1, 1, 1)}, VPN: vpn.OpenVPN, UDP: true},
				{IPs: []net.IP{net.IPv4(2, 2, 2, 2)}, VPN: vpn.OpenVPN, UDP: true},
				{IPs: []net.IP{net.IPv4(3, 3, 3, 3)}, VPN: vpn.OpenVPN, UDP: true},
			},
			connection: models.Connection{
				Type:     vpn.OpenVPN,
				IP:       net.IPv4(2, 2, 2, 2),
				Port:     1194,
				Protocol: constants.UDP,
			},
		},
		"with filter": {
			selection: settings.ServerSelection{
				Hostnames: []string{"b"},
			}.WithDefaults(providers.Windscribe),
			servers: []models.Server{
				{Hostname: "a", IPs: []net.IP{net.IPv4(1, 1, 1, 1)}, VPN: vpn.OpenVPN, UDP: true},
				{Hostname: "b", IPs: []net.IP{net.IPv4(2, 2, 2, 2)}, VPN: vpn.OpenVPN, UDP: true},
				{Hostname: "a", IPs: []net.IP{net.IPv4(3, 3, 3, 3)}, VPN: vpn.OpenVPN, UDP: true},
			},
			connection: models.Connection{
				Type:     vpn.OpenVPN,
				IP:       net.IPv4(2, 2, 2, 2),
				Port:     1194,
				Hostname: "b",
				Protocol: constants.UDP,
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			randSource := rand.NewSource(0)

			provider := New(testCase.servers, randSource)

			connection, err := provider.GetConnection(testCase.selection)

			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
			assert.Equal(t, testCase.connection, connection)
		})
	}
}
