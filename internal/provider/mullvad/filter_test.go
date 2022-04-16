package mullvad

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool { return &b }

func Test_Mullvad_filterServers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers   []models.MullvadServer
		selection settings.ServerSelection
		filtered  []models.MullvadServer
		err       error
	}{
		"no server available": {
			selection: settings.ServerSelection{}.WithDefaults(providers.Mullvad),
			err:       errors.New("no server found: for VPN openvpn; protocol udp"),
		},
		"no filter": {
			servers: []models.MullvadServer{
				{VPN: constants.OpenVPN, Hostname: "a"},
				{VPN: constants.OpenVPN, Hostname: "b"},
				{VPN: constants.OpenVPN, Hostname: "c"},
			},
			selection: settings.ServerSelection{}.WithDefaults(providers.Mullvad),
			filtered: []models.MullvadServer{
				{VPN: constants.OpenVPN, Hostname: "a"},
				{VPN: constants.OpenVPN, Hostname: "b"},
				{VPN: constants.OpenVPN, Hostname: "c"},
			},
		},
		"filter OpenVPN out": {
			selection: settings.ServerSelection{
				VPN: constants.Wireguard,
			}.WithDefaults(providers.Mullvad),
			servers: []models.MullvadServer{
				{VPN: constants.OpenVPN, Hostname: "a"},
				{VPN: constants.Wireguard, Hostname: "b"},
				{VPN: constants.OpenVPN, Hostname: "c"},
			},
			filtered: []models.MullvadServer{
				{VPN: constants.Wireguard, Hostname: "b"},
			},
		},
		"filter by country": {
			selection: settings.ServerSelection{
				Countries: []string{"b"},
			}.WithDefaults(providers.Mullvad),
			servers: []models.MullvadServer{
				{VPN: constants.OpenVPN, Country: "a"},
				{VPN: constants.OpenVPN, Country: "b"},
				{VPN: constants.OpenVPN, Country: "c"},
			},
			filtered: []models.MullvadServer{
				{VPN: constants.OpenVPN, Country: "b"},
			},
		},
		"filter by city": {
			selection: settings.ServerSelection{
				Cities: []string{"b"},
			}.WithDefaults(providers.Mullvad),
			servers: []models.MullvadServer{
				{VPN: constants.OpenVPN, City: "a"},
				{VPN: constants.OpenVPN, City: "b"},
				{VPN: constants.OpenVPN, City: "c"},
			},
			filtered: []models.MullvadServer{
				{VPN: constants.OpenVPN, City: "b"},
			},
		},
		"filter by ISP": {
			selection: settings.ServerSelection{
				ISPs: []string{"b"},
			}.WithDefaults(providers.Mullvad),
			servers: []models.MullvadServer{
				{VPN: constants.OpenVPN, ISP: "a"},
				{VPN: constants.OpenVPN, ISP: "b"},
				{VPN: constants.OpenVPN, ISP: "c"},
			},
			filtered: []models.MullvadServer{
				{VPN: constants.OpenVPN, ISP: "b"},
			},
		},
		"filter by hostname": {
			selection: settings.ServerSelection{
				Hostnames: []string{"b"},
			}.WithDefaults(providers.Mullvad),
			servers: []models.MullvadServer{
				{VPN: constants.OpenVPN, Hostname: "a"},
				{VPN: constants.OpenVPN, Hostname: "b"},
				{VPN: constants.OpenVPN, Hostname: "c"},
			},
			filtered: []models.MullvadServer{
				{VPN: constants.OpenVPN, Hostname: "b"},
			},
		},
		"filter by owned": {
			selection: settings.ServerSelection{
				OwnedOnly: boolPtr(true),
			}.WithDefaults(providers.Mullvad),
			servers: []models.MullvadServer{
				{VPN: constants.OpenVPN, Hostname: "a"},
				{VPN: constants.OpenVPN, Hostname: "b", Owned: true},
				{VPN: constants.OpenVPN, Hostname: "c"},
			},
			filtered: []models.MullvadServer{
				{VPN: constants.OpenVPN, Hostname: "b", Owned: true},
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
