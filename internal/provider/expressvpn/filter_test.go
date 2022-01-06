package expressvpn

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool { return &b }

func Test_Expressvpn_filterServers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers   []models.ExpressvpnServer
		selection settings.ServerSelection
		filtered  []models.ExpressvpnServer
		err       error
	}{
		"no server available": {
			selection: settings.ServerSelection{}.WithDefaults(constants.Expressvpn),
			err:       errors.New("no server found: for VPN openvpn; protocol udp"),
		},
		"no filter": {
			servers: []models.ExpressvpnServer{
				{Hostname: "a", UDP: true},
				{Hostname: "b", UDP: true},
				{Hostname: "c", UDP: true},
			},
			selection: settings.ServerSelection{}.WithDefaults(constants.Expressvpn),
			filtered: []models.ExpressvpnServer{
				{Hostname: "a", UDP: true},
				{Hostname: "b", UDP: true},
				{Hostname: "c", UDP: true},
			},
		},
		"filter by country": {
			selection: settings.ServerSelection{
				Countries: []string{"b"},
			}.WithDefaults(constants.Expressvpn),
			servers: []models.ExpressvpnServer{
				{Country: "a", UDP: true},
				{Country: "b", UDP: true},
				{Country: "c", UDP: true},
			},
			filtered: []models.ExpressvpnServer{
				{Country: "b", UDP: true},
			},
		},
		"filter by city": {
			selection: settings.ServerSelection{
				Cities: []string{"b"},
			}.WithDefaults(constants.Expressvpn),
			servers: []models.ExpressvpnServer{
				{City: "a", UDP: true},
				{City: "b", UDP: true},
				{City: "c", UDP: true},
			},
			filtered: []models.ExpressvpnServer{
				{City: "b", UDP: true},
			},
		},
		"filter by hostname": {
			selection: settings.ServerSelection{
				Hostnames: []string{"b"},
			}.WithDefaults(constants.Expressvpn),
			servers: []models.ExpressvpnServer{
				{Hostname: "a", UDP: true},
				{Hostname: "b", UDP: true},
				{Hostname: "c", UDP: true},
			},
			filtered: []models.ExpressvpnServer{
				{Hostname: "b", UDP: true},
			},
		},
		"filter by protocol": {
			selection: settings.ServerSelection{
				OpenVPN: settings.OpenVPNSelection{
					TCP: boolPtr(true),
				},
			}.WithDefaults(constants.Expressvpn),
			servers: []models.ExpressvpnServer{
				{Hostname: "a", UDP: true},
				{Hostname: "b", UDP: true, TCP: true},
				{Hostname: "c", UDP: true},
			},
			filtered: []models.ExpressvpnServer{
				{Hostname: "b", UDP: true, TCP: true},
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
