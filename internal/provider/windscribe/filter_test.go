package windscribe

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Windscribe_filterServers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers   []models.WindscribeServer
		selection configuration.ServerSelection
		filtered  []models.WindscribeServer
		err       error
	}{
		"no server available": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
			},
			err: errors.New("no server found: for VPN openvpn; protocol udp"),
		},
		"no filter": {
			servers: []models.WindscribeServer{
				{Hostname: "a"},
				{Hostname: "b"},
				{Hostname: "c"},
			},
			filtered: []models.WindscribeServer{
				{Hostname: "a"},
				{Hostname: "b"},
				{Hostname: "c"},
			},
		},
		"filter OpenVPN out": {
			selection: configuration.ServerSelection{
				VPN: constants.Wireguard,
			},
			servers: []models.WindscribeServer{
				{VPN: constants.OpenVPN, Hostname: "a"},
				{VPN: constants.Wireguard, Hostname: "b"},
				{VPN: constants.OpenVPN, Hostname: "c"},
			},
			filtered: []models.WindscribeServer{
				{VPN: constants.Wireguard, Hostname: "b"},
			},
		},
		"filter by region": {
			selection: configuration.ServerSelection{
				Regions: []string{"b"},
			},
			servers: []models.WindscribeServer{
				{Region: "a"},
				{Region: "b"},
				{Region: "c"},
			},
			filtered: []models.WindscribeServer{
				{Region: "b"},
			},
		},
		"filter by city": {
			selection: configuration.ServerSelection{
				Cities: []string{"b"},
			},
			servers: []models.WindscribeServer{
				{City: "a"},
				{City: "b"},
				{City: "c"},
			},
			filtered: []models.WindscribeServer{
				{City: "b"},
			},
		},
		"filter by hostname": {
			selection: configuration.ServerSelection{
				Hostnames: []string{"b"},
			},
			servers: []models.WindscribeServer{
				{Hostname: "a"},
				{Hostname: "b"},
				{Hostname: "c"},
			},
			filtered: []models.WindscribeServer{
				{Hostname: "b"},
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
