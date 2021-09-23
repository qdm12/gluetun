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
				{Country: "a", Group: "Premium TCP Asia"},
				{Country: "b", Group: "Premium TCP Europe"},
				{Country: "c", Group: "Premium UDP Asia"},
				{Country: "d", Group: "Premium UDP Europe"},
			},
			filteredServers: []models.CyberghostServer{
				{Country: "c", Group: "Premium UDP Asia"},
				{Country: "d", Group: "Premium UDP Europe"},
			},
		},
		"servers with TCP selection": {
			servers: []models.CyberghostServer{
				{Country: "a", Group: "Premium TCP Asia"},
				{Country: "b", Group: "Premium TCP Europe"},
				{Country: "c", Group: "Premium UDP Asia"},
				{Country: "d", Group: "Premium UDP Europe"},
			},
			selection: configuration.ServerSelection{
				OpenVPN: configuration.OpenVPNSelection{
					TCP: true,
				},
			},
			filteredServers: []models.CyberghostServer{
				{Country: "a", Group: "Premium TCP Asia"},
				{Country: "b", Group: "Premium TCP Europe"},
			},
		},
		"servers with regions filter": {
			servers: []models.CyberghostServer{
				{Country: "a", Group: "Premium UDP Asia"},
				{Country: "b", Group: "Premium UDP Asia"},
				{Country: "c", Group: "Premium UDP Asia"},
				{Country: "d", Group: "Premium UDP Asia"},
			},
			selection: configuration.ServerSelection{
				Countries: []string{"a", "c"},
			},
			filteredServers: []models.CyberghostServer{
				{Country: "a", Group: "Premium UDP Asia"},
				{Country: "c", Group: "Premium UDP Asia"},
			},
		},
		"servers with group filter": {
			servers: []models.CyberghostServer{
				{Country: "a", Group: "Premium UDP Europe"},
				{Country: "b", Group: "Premium UDP Europe"},
				{Country: "c", Group: "Premium TCP Europe"},
				{Country: "d", Group: "Premium TCP Europe"},
			},
			selection: configuration.ServerSelection{
				Groups: []string{"Premium UDP Europe"},
			},
			filteredServers: []models.CyberghostServer{
				{Country: "a", Group: "Premium UDP Europe"},
				{Country: "b", Group: "Premium UDP Europe"},
			},
		},
		"servers with bad group filter": {
			servers: []models.CyberghostServer{
				{Country: "a", Group: "Premium TCP Europe"},
				{Country: "b", Group: "Premium TCP Europe"},
				{Country: "c", Group: "Premium UDP Europe"},
				{Country: "d", Group: "Premium UDP Europe"},
			},
			selection: configuration.ServerSelection{
				Groups: []string{"Premium TCP Europe"},
			},
			err: errors.New("server group does not match protocol: group Premium TCP Europe for protocol UDP"),
		},
		"servers with regions and group filter": {
			servers: []models.CyberghostServer{
				{Country: "a", Group: "Premium UDP Europe"},
				{Country: "b", Group: "Premium TCP Europe"},
				{Country: "c", Group: "Premium UDP Asia"},
				{Country: "d", Group: "Premium TCP Asia"},
			},
			selection: configuration.ServerSelection{
				Countries: []string{"a", "c"},
				Groups:    []string{"Premium UDP Europe"},
			},
			filteredServers: []models.CyberghostServer{
				{Country: "a", Group: "Premium UDP Europe"},
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

	servers := []models.CyberghostServer{
		{Group: "Premium TCP Asia"},
		{Group: "Premium TCP Europe"},
		{Group: "Premium TCP USA"},
		{Group: "Premium UDP Asia"},
		{Group: "Premium UDP Europe"},
		{Group: "Premium UDP USA"},
	}
	expected := []string{
		"Premium TCP Asia", "Premium TCP Europe", "Premium TCP USA",
	}
	choices := tcpGroupChoices(servers)

	assert.Equal(t, expected, choices)
}

func Test_udpGroupChoices(t *testing.T) {
	t.Parallel()

	servers := []models.CyberghostServer{
		{Group: "Premium TCP Asia"},
		{Group: "Premium TCP Europe"},
		{Group: "Premium TCP USA"},
		{Group: "Premium UDP Asia"},
		{Group: "Premium UDP Europe"},
		{Group: "Premium UDP USA"},
	}
	expected := []string{
		"Premium UDP Asia", "Premium UDP Europe", "Premium UDP USA",
	}
	choices := udpGroupChoices(servers)

	assert.Equal(t, expected, choices)
}
