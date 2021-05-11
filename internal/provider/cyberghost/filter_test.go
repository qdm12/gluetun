package cyberghost

import (
	"errors"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration"
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
		"no servers": {
			err: errors.New("no server found: for protocol udp"),
		},
		"servers without filter": {
			servers: []models.CyberghostServer{
				{Region: "a", Group: "1"},
				{Region: "b", Group: "1"},
				{Region: "c", Group: "2"},
				{Region: "d", Group: "2"},
			},
			filteredServers: []models.CyberghostServer{
				{Region: "a", Group: "1"},
				{Region: "b", Group: "1"},
				{Region: "c", Group: "2"},
				{Region: "d", Group: "2"},
			},
		},
		"servers with regions filter": {
			servers: []models.CyberghostServer{
				{Region: "a", Group: "1"},
				{Region: "b", Group: "1"},
				{Region: "c", Group: "2"},
				{Region: "d", Group: "2"},
			},
			selection: configuration.ServerSelection{
				Regions: []string{"a", "c"},
			},
			filteredServers: []models.CyberghostServer{
				{Region: "a", Group: "1"},
				{Region: "c", Group: "2"},
			},
		},
		"servers with group filter": {
			servers: []models.CyberghostServer{
				{Region: "a", Group: "1"},
				{Region: "b", Group: "1"},
				{Region: "c", Group: "2"},
				{Region: "d", Group: "2"},
			},
			selection: configuration.ServerSelection{
				Group: "1",
			},
			filteredServers: []models.CyberghostServer{
				{Region: "a", Group: "1"},
				{Region: "b", Group: "1"},
			},
		},
		"servers with regions and group filter": {
			servers: []models.CyberghostServer{
				{Region: "a", Group: "1"},
				{Region: "b", Group: "1"},
				{Region: "c", Group: "2"},
				{Region: "d", Group: "2"},
			},
			selection: configuration.ServerSelection{
				Regions: []string{"a", "c"},
				Group:   "1",
			},
			filteredServers: []models.CyberghostServer{
				{Region: "a", Group: "1"},
			},
		},
		"servers with hostnames filter": {
			servers: []models.CyberghostServer{
				{Hostname: "a"},
				{Hostname: "b"},
				{Hostname: "c"},
			},
			selection: configuration.ServerSelection{
				Hostnames: []string{"a", "c"},
			},
			filteredServers: []models.CyberghostServer{
				{Hostname: "a"},
				{Hostname: "c"},
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
