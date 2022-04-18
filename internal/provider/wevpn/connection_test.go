package wevpn

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

func Test_Wevpn_GetConnection(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers    []models.Server
		selection  settings.ServerSelection
		connection models.Connection
		errWrapped error
		errMessage string
	}{
		"no server available": {
			selection: settings.ServerSelection{
				VPN: vpn.OpenVPN,
			}.WithDefaults(providers.Wevpn),
			errWrapped: utils.ErrNoServerFound,
			errMessage: "no server found: for VPN openvpn; protocol udp",
		},
		"no filter": {
			servers: []models.Server{
				{IPs: []net.IP{net.IPv4(1, 1, 1, 1)}, VPN: vpn.OpenVPN, UDP: true},
				{IPs: []net.IP{net.IPv4(2, 2, 2, 2)}, VPN: vpn.OpenVPN, UDP: true},
				{IPs: []net.IP{net.IPv4(3, 3, 3, 3)}, VPN: vpn.OpenVPN, UDP: true},
			},
			selection: settings.ServerSelection{}.WithDefaults(providers.Wevpn),
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
			}.WithDefaults(providers.Wevpn),
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
			}.WithDefaults(providers.Wevpn),
			servers: []models.Server{
				{Hostname: "a", IPs: []net.IP{net.IPv4(1, 1, 1, 1)}, VPN: vpn.OpenVPN, UDP: true},
				{Hostname: "b", IPs: []net.IP{net.IPv4(2, 2, 2, 2)}, VPN: vpn.OpenVPN, UDP: true},
				{Hostname: "a", IPs: []net.IP{net.IPv4(3, 3, 3, 3)}, VPN: vpn.OpenVPN, UDP: true},
			},
			connection: models.Connection{
				Type:     vpn.OpenVPN,
				IP:       net.IPv4(2, 2, 2, 2),
				Port:     1194,
				Protocol: constants.UDP,
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

			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}

			assert.Equal(t, testCase.connection, connection)
		})
	}
}
