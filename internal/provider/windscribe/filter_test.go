package windscribe

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

func Test_Windscribe_filterServers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers   []models.Server
		selection settings.ServerSelection
		filtered  []models.Server
		err       error
	}{
		"no server available": {
			selection: settings.ServerSelection{}.WithDefaults(providers.Windscribe),
			err:       errors.New("no server found: for VPN openvpn; protocol udp"),
		},
		"no filter": {
			servers: []models.Server{
				{VPN: constants.OpenVPN, Hostname: "a"},
				{VPN: constants.OpenVPN, Hostname: "b"},
				{VPN: constants.OpenVPN, Hostname: "c"},
			},
			selection: settings.ServerSelection{}.WithDefaults(providers.Windscribe),
			filtered: []models.Server{
				{VPN: constants.OpenVPN, Hostname: "a"},
				{VPN: constants.OpenVPN, Hostname: "b"},
				{VPN: constants.OpenVPN, Hostname: "c"},
			},
		},
		"filter OpenVPN out": {
			selection: settings.ServerSelection{
				VPN: constants.Wireguard,
			}.WithDefaults(providers.Windscribe),
			servers: []models.Server{
				{VPN: constants.OpenVPN, Hostname: "a"},
				{VPN: constants.Wireguard, Hostname: "b"},
				{VPN: constants.OpenVPN, Hostname: "c"},
			},
			filtered: []models.Server{
				{VPN: constants.Wireguard, Hostname: "b"},
			},
		},
		"filter by region": {
			selection: settings.ServerSelection{
				Regions: []string{"b"},
			}.WithDefaults(providers.Windscribe),
			servers: []models.Server{
				{VPN: constants.OpenVPN, Region: "a"},
				{VPN: constants.OpenVPN, Region: "b"},
				{VPN: constants.OpenVPN, Region: "c"},
			},
			filtered: []models.Server{
				{VPN: constants.OpenVPN, Region: "b"},
			},
		},
		"filter by city": {
			selection: settings.ServerSelection{
				Cities: []string{"b"},
			}.WithDefaults(providers.Windscribe),
			servers: []models.Server{
				{VPN: constants.OpenVPN, City: "a"},
				{VPN: constants.OpenVPN, City: "b"},
				{VPN: constants.OpenVPN, City: "c"},
			},
			filtered: []models.Server{
				{VPN: constants.OpenVPN, City: "b"},
			},
		},
		"filter by hostname": {
			selection: settings.ServerSelection{
				Hostnames: []string{"b"},
			}.WithDefaults(providers.Windscribe),
			servers: []models.Server{
				{VPN: constants.OpenVPN, Hostname: "a"},
				{VPN: constants.OpenVPN, Hostname: "b"},
				{VPN: constants.OpenVPN, Hostname: "c"},
			},
			filtered: []models.Server{
				{VPN: constants.OpenVPN, Hostname: "b"},
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
