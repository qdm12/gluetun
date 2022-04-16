package expressvpn

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool { return &b }

func Test_Expressvpn_filterServers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers   []models.Server
		selection settings.ServerSelection
		filtered  []models.Server
		err       error
	}{
		"no server available": {
			selection: settings.ServerSelection{}.WithDefaults(providers.Expressvpn),
			err:       errors.New("no server found: for VPN openvpn; protocol udp"),
		},
		"no filter": {
			servers: []models.Server{
				{Hostname: "a", UDP: true},
				{Hostname: "b", UDP: true},
				{Hostname: "c", UDP: true},
			},
			selection: settings.ServerSelection{}.WithDefaults(providers.Expressvpn),
			filtered: []models.Server{
				{Hostname: "a", UDP: true},
				{Hostname: "b", UDP: true},
				{Hostname: "c", UDP: true},
			},
		},
		"filter by country": {
			selection: settings.ServerSelection{
				Countries: []string{"b"},
			}.WithDefaults(providers.Expressvpn),
			servers: []models.Server{
				{Country: "a", UDP: true},
				{Country: "b", UDP: true},
				{Country: "c", UDP: true},
			},
			filtered: []models.Server{
				{Country: "b", UDP: true},
			},
		},
		"filter by city": {
			selection: settings.ServerSelection{
				Cities: []string{"b"},
			}.WithDefaults(providers.Expressvpn),
			servers: []models.Server{
				{City: "a", UDP: true},
				{City: "b", UDP: true},
				{City: "c", UDP: true},
			},
			filtered: []models.Server{
				{City: "b", UDP: true},
			},
		},
		"filter by hostname": {
			selection: settings.ServerSelection{
				Hostnames: []string{"b"},
			}.WithDefaults(providers.Expressvpn),
			servers: []models.Server{
				{Hostname: "a", UDP: true},
				{Hostname: "b", UDP: true},
				{Hostname: "c", UDP: true},
			},
			filtered: []models.Server{
				{Hostname: "b", UDP: true},
			},
		},
		"filter by protocol": {
			selection: settings.ServerSelection{
				OpenVPN: settings.OpenVPNSelection{
					TCP: boolPtr(true),
				},
			}.WithDefaults(providers.Expressvpn),
			servers: []models.Server{
				{Hostname: "a", UDP: true},
				{Hostname: "b", UDP: true, TCP: true},
				{Hostname: "c", UDP: true},
			},
			filtered: []models.Server{
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
