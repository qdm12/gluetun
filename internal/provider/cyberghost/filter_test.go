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
		"servers without filter defaults to UDP": {
			servers: []models.CyberghostServer{
				{Region: "a", Group: "Premium TCP Asia"},
				{Region: "b", Group: "Premium TCP Europe"},
				{Region: "c", Group: "Premium UDP Asia"},
				{Region: "d", Group: "Premium UDP Europe"},
			},
			filteredServers: []models.CyberghostServer{
				{Region: "c", Group: "Premium UDP Asia"},
				{Region: "d", Group: "Premium UDP Europe"},
			},
		},
		"servers with TCP selection": {
			servers: []models.CyberghostServer{
				{Region: "a", Group: "Premium TCP Asia"},
				{Region: "b", Group: "Premium TCP Europe"},
				{Region: "c", Group: "Premium UDP Asia"},
				{Region: "d", Group: "Premium UDP Europe"},
			},
			selection: configuration.ServerSelection{
				TCP: true,
			},
			filteredServers: []models.CyberghostServer{
				{Region: "a", Group: "Premium TCP Asia"},
				{Region: "b", Group: "Premium TCP Europe"},
			},
		},
		"servers with regions filter": {
			servers: []models.CyberghostServer{
				{Region: "a", Group: "Premium UDP Asia"},
				{Region: "b", Group: "Premium UDP Asia"},
				{Region: "c", Group: "Premium UDP Asia"},
				{Region: "d", Group: "Premium UDP Asia"},
			},
			selection: configuration.ServerSelection{
				Regions: []string{"a", "c"},
			},
			filteredServers: []models.CyberghostServer{
				{Region: "a", Group: "Premium UDP Asia"},
				{Region: "c", Group: "Premium UDP Asia"},
			},
		},
		"servers with group filter": {
			servers: []models.CyberghostServer{
				{Region: "a", Group: "Premium UDP Europe"},
				{Region: "b", Group: "Premium UDP Europe"},
				{Region: "c", Group: "Premium TCP Europe"},
				{Region: "d", Group: "Premium TCP Europe"},
			},
			selection: configuration.ServerSelection{
				Groups: []string{"Premium UDP Europe"},
			},
			filteredServers: []models.CyberghostServer{
				{Region: "a", Group: "Premium UDP Europe"},
				{Region: "b", Group: "Premium UDP Europe"},
			},
		},
		"servers with bad group filter": {
			servers: []models.CyberghostServer{
				{Region: "a", Group: "Premium TCP Europe"},
				{Region: "b", Group: "Premium TCP Europe"},
				{Region: "c", Group: "Premium UDP Europe"},
				{Region: "d", Group: "Premium UDP Europe"},
			},
			selection: configuration.ServerSelection{
				Groups: []string{"Premium TCP Europe"},
			},
			err: errors.New("server group does not match protocol: group Premium TCP Europe for protocol UDP"),
		},
		"servers with regions and group filter": {
			servers: []models.CyberghostServer{
				{Region: "a", Group: "Premium UDP Europe"},
				{Region: "b", Group: "Premium TCP Europe"},
				{Region: "c", Group: "Premium UDP Asia"},
				{Region: "d", Group: "Premium TCP Asia"},
			},
			selection: configuration.ServerSelection{
				Regions: []string{"a", "c"},
				Groups:  []string{"Premium UDP Europe"},
			},
			filteredServers: []models.CyberghostServer{
				{Region: "a", Group: "Premium UDP Europe"},
			},
		},
		"servers with hostnames filter": {
			servers: []models.CyberghostServer{
				{Hostname: "a", Group: "Premium UDP Asia"},
				{Hostname: "b", Group: "Premium UDP Asia"},
				{Hostname: "c", Group: "Premium UDP Asia"},
			},
			selection: configuration.ServerSelection{
				Hostnames: []string{"a", "c"},
			},
			filteredServers: []models.CyberghostServer{
				{Hostname: "a", Group: "Premium UDP Asia"},
				{Hostname: "c", Group: "Premium UDP Asia"},
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

func Test_tcpGroupChoices(t *testing.T) {
	t.Parallel()

	expected := []string{
		"Premium TCP Asia", "Premium TCP Europe", "Premium TCP USA",
	}
	choices := tcpGroupChoices()

	assert.Equal(t, expected, choices)
}

func Test_udpGroupChoices(t *testing.T) {
	t.Parallel()

	expected := []string{
		"Premium UDP Asia", "Premium UDP Europe", "Premium UDP USA",
	}
	choices := udpGroupChoices()

	assert.Equal(t, expected, choices)
}
