package cyberghost

import (
	"errors"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Cyberghost_filterServers(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		servers         []models.CyberghostServer
		selection       configuration.ServerSelection
		filteredServers []models.CyberghostServer
		err             error
	}{
		"no server": {
			selection: configuration.ServerSelection{VPN: constants.OpenVPN},
			err:       errors.New("no server found: for VPN openvpn; protocol udp"),
		},
		"servers without filter defaults to UDP": {
			servers: []models.CyberghostServer{
				{Country: "a", TCP: true},
				{Country: "b", TCP: true},
				{Country: "c", UDP: true},
				{Country: "d", UDP: true},
			},
			filteredServers: []models.CyberghostServer{
				{Country: "c", UDP: true},
				{Country: "d", UDP: true},
			},
		},
		"servers with TCP selection": {
			servers: []models.CyberghostServer{
				{Country: "a", TCP: true},
				{Country: "b", TCP: true},
				{Country: "c", UDP: true},
				{Country: "d", UDP: true},
			},
			selection: configuration.ServerSelection{
				OpenVPN: configuration.OpenVPNSelection{
					TCP: true,
				},
			},
			filteredServers: []models.CyberghostServer{
				{Country: "a", TCP: true},
				{Country: "b", TCP: true},
			},
		},
		"servers with regions filter": {
			servers: []models.CyberghostServer{
				{Country: "a", UDP: true},
				{Country: "b", UDP: true},
				{Country: "c", UDP: true},
				{Country: "d", UDP: true},
			},
			selection: configuration.ServerSelection{
				Countries: []string{"a", "c"},
			},
			filteredServers: []models.CyberghostServer{
				{Country: "a", UDP: true},
				{Country: "c", UDP: true},
			},
		},
		"servers with hostnames filter": {
			servers: []models.CyberghostServer{
				{Hostname: "a", UDP: true},
				{Hostname: "b", UDP: true},
				{Hostname: "c", UDP: true},
			},
			selection: configuration.ServerSelection{
				Hostnames: []string{"a", "c"},
			},
			filteredServers: []models.CyberghostServer{
				{Hostname: "a", UDP: true},
				{Hostname: "c", UDP: true},
			},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c := &Cyberghost{servers: testCase.servers}
			filteredServers, err := c.filterServers(testCase.selection)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.filteredServers, filteredServers)
		})
	}
}
