package mullvad

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

func Test_Mullvad_filterServers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers   []models.MullvadServer
		selection configuration.ServerSelection
		filtered  []models.MullvadServer
		err       error
	}{
		"no server available": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
			},
			err: errors.New("no server found: for VPN openvpn; protocol udp"),
		},
		"no filter": {
			servers: []models.MullvadServer{
				{Hostname: "a"},
				{Hostname: "b"},
				{Hostname: "c"},
			},
			filtered: []models.MullvadServer{
				{Hostname: "a"},
				{Hostname: "b"},
				{Hostname: "c"},
			},
		},
		"filter OpenVPN out": {
			selection: configuration.ServerSelection{
				VPN: constants.Wireguard,
			},
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
			selection: configuration.ServerSelection{
				Countries: []string{"b"},
			},
			servers: []models.MullvadServer{
				{Country: "a"},
				{Country: "b"},
				{Country: "c"},
			},
			filtered: []models.MullvadServer{
				{Country: "b"},
			},
		},
		"filter by city": {
			selection: configuration.ServerSelection{
				Cities: []string{"b"},
			},
			servers: []models.MullvadServer{
				{City: "a"},
				{City: "b"},
				{City: "c"},
			},
			filtered: []models.MullvadServer{
				{City: "b"},
			},
		},
		"filter by ISP": {
			selection: configuration.ServerSelection{
				ISPs: []string{"b"},
			},
			servers: []models.MullvadServer{
				{ISP: "a"},
				{ISP: "b"},
				{ISP: "c"},
			},
			filtered: []models.MullvadServer{
				{ISP: "b"},
			},
		},
		"filter by hostname": {
			selection: configuration.ServerSelection{
				Hostnames: []string{"b"},
			},
			servers: []models.MullvadServer{
				{Hostname: "a"},
				{Hostname: "b"},
				{Hostname: "c"},
			},
			filtered: []models.MullvadServer{
				{Hostname: "b"},
			},
		},
		"filter by owned": {
			selection: configuration.ServerSelection{
				Owned: true,
			},
			servers: []models.MullvadServer{
				{Hostname: "a"},
				{Hostname: "b", Owned: true},
				{Hostname: "c"},
			},
			filtered: []models.MullvadServer{
				{Hostname: "b", Owned: true},
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
