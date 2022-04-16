package wevpn

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

func Test_Wevpn_filterServers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers   []models.Server
		selection settings.ServerSelection
		filtered  []models.Server
		err       error
	}{
		"no server available": {
			selection: settings.ServerSelection{}.WithDefaults(providers.Wevpn),
			err:       errors.New("no server found: for VPN openvpn; protocol udp"),
		},
		"no filter": {
			servers: []models.Server{
				{Hostname: "a", UDP: true},
				{Hostname: "b", UDP: true},
				{Hostname: "c", UDP: true},
			},
			selection: settings.ServerSelection{}.WithDefaults(providers.Wevpn),
			filtered: []models.Server{
				{Hostname: "a", UDP: true},
				{Hostname: "b", UDP: true},
				{Hostname: "c", UDP: true},
			},
		},
		"filter by protocol": {
			selection: settings.ServerSelection{
				OpenVPN: settings.OpenVPNSelection{TCP: boolPtr(true)},
			}.WithDefaults(providers.Wevpn),
			servers: []models.Server{
				{Hostname: "a", UDP: true},
				{Hostname: "b", TCP: true},
				{Hostname: "c", UDP: true},
			},
			filtered: []models.Server{
				{Hostname: "b", TCP: true},
			},
		},
		"filter by city": {
			selection: settings.ServerSelection{
				Cities: []string{"b"},
			}.WithDefaults(providers.Wevpn),
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
			}.WithDefaults(providers.Wevpn),
			servers: []models.Server{
				{Hostname: "a", UDP: true},
				{Hostname: "b", UDP: true},
				{Hostname: "c", UDP: true},
			},
			filtered: []models.Server{
				{Hostname: "b", UDP: true},
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
