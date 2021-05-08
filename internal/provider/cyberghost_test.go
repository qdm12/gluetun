package provider

import (
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_cyberghost_filterServers(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		servers         []models.CyberghostServer
		regions         []string
		hostnames       []string
		group           string
		filteredServers []models.CyberghostServer
	}{
		"no servers": {},
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
			regions: []string{"a", "c"},
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
			group: "1",
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
			regions: []string{"a", "c"},
			group:   "1",
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
			hostnames: []string{"a", "c"},
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
			c := &cyberghost{servers: testCase.servers}
			filteredServers := c.filterServers(testCase.regions, testCase.hostnames, testCase.group)
			assert.Equal(t, testCase.filteredServers, filteredServers)
		})
	}
}
