package mullvad

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool { return &b }

func Test_Mullvad_filterServers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers   []models.Server
		selection settings.ServerSelection
		filtered  []models.Server
		err       error
	}{
		"no server available": {
			selection: settings.ServerSelection{}.WithDefaults(providers.Mullvad),
			err:       errors.New("no server found: for VPN openvpn; protocol udp"),
		},
		"no filter": {
			servers: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "a"},
				{VPN: vpn.OpenVPN, Hostname: "b"},
				{VPN: vpn.OpenVPN, Hostname: "c"},
			},
			selection: settings.ServerSelection{}.WithDefaults(providers.Mullvad),
			filtered: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "a"},
				{VPN: vpn.OpenVPN, Hostname: "b"},
				{VPN: vpn.OpenVPN, Hostname: "c"},
			},
		},
		"filter OpenVPN out": {
			selection: settings.ServerSelection{
				VPN: vpn.Wireguard,
			}.WithDefaults(providers.Mullvad),
			servers: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "a"},
				{VPN: vpn.Wireguard, Hostname: "b"},
				{VPN: vpn.OpenVPN, Hostname: "c"},
			},
			filtered: []models.Server{
				{VPN: vpn.Wireguard, Hostname: "b"},
			},
		},
		"filter by country": {
			selection: settings.ServerSelection{
				Countries: []string{"b"},
			}.WithDefaults(providers.Mullvad),
			servers: []models.Server{
				{VPN: vpn.OpenVPN, Country: "a"},
				{VPN: vpn.OpenVPN, Country: "b"},
				{VPN: vpn.OpenVPN, Country: "c"},
			},
			filtered: []models.Server{
				{VPN: vpn.OpenVPN, Country: "b"},
			},
		},
		"filter by city": {
			selection: settings.ServerSelection{
				Cities: []string{"b"},
			}.WithDefaults(providers.Mullvad),
			servers: []models.Server{
				{VPN: vpn.OpenVPN, City: "a"},
				{VPN: vpn.OpenVPN, City: "b"},
				{VPN: vpn.OpenVPN, City: "c"},
			},
			filtered: []models.Server{
				{VPN: vpn.OpenVPN, City: "b"},
			},
		},
		"filter by ISP": {
			selection: settings.ServerSelection{
				ISPs: []string{"b"},
			}.WithDefaults(providers.Mullvad),
			servers: []models.Server{
				{VPN: vpn.OpenVPN, ISP: "a"},
				{VPN: vpn.OpenVPN, ISP: "b"},
				{VPN: vpn.OpenVPN, ISP: "c"},
			},
			filtered: []models.Server{
				{VPN: vpn.OpenVPN, ISP: "b"},
			},
		},
		"filter by hostname": {
			selection: settings.ServerSelection{
				Hostnames: []string{"b"},
			}.WithDefaults(providers.Mullvad),
			servers: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "a"},
				{VPN: vpn.OpenVPN, Hostname: "b"},
				{VPN: vpn.OpenVPN, Hostname: "c"},
			},
			filtered: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "b"},
			},
		},
		"filter by owned": {
			selection: settings.ServerSelection{
				OwnedOnly: boolPtr(true),
			}.WithDefaults(providers.Mullvad),
			servers: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "a"},
				{VPN: vpn.OpenVPN, Hostname: "b", Owned: true},
				{VPN: vpn.OpenVPN, Hostname: "c"},
			},
			filtered: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "b", Owned: true},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			randSource := rand.NewSource(0)

			m := New(testCase.servers, randSource)

			servers, err := m.filterServers(testCase.selection)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.filtered, servers)
		})
	}
}
