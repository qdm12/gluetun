package wevpn

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

func Test_Wevpn_filterServers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers   []models.WevpnServer
		selection configuration.ServerSelection
		filtered  []models.WevpnServer
		err       error
	}{
		"no server available": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
			},
			err: errors.New("no server found: for VPN openvpn; protocol udp"),
		},
		"no filter": {
			servers: []models.WevpnServer{
				{Hostname: "a"},
				{Hostname: "b"},
				{Hostname: "c"},
			},
			filtered: []models.WevpnServer{
				{Hostname: "a"},
				{Hostname: "b"},
				{Hostname: "c"},
			},
		},
		"filter by city": {
			selection: configuration.ServerSelection{
				Cities: []string{"b"},
			},
			servers: []models.WevpnServer{
				{City: "a"},
				{City: "b"},
				{City: "c"},
			},
			filtered: []models.WevpnServer{
				{City: "b"},
			},
		},
		"filter by hostname": {
			selection: configuration.ServerSelection{
				Hostnames: []string{"b"},
			},
			servers: []models.WevpnServer{
				{Hostname: "a"},
				{Hostname: "b"},
				{Hostname: "c"},
			},
			filtered: []models.WevpnServer{
				{Hostname: "b"},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			randSource := rand.NewSource(0)

			w := New(testCase.servers, randSource)

			servers, err := w.filterServers(testCase.selection)

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
