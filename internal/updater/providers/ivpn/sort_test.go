package ivpn

import (
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_sortServers(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		initialServers []models.IvpnServer
		sortedServers  []models.IvpnServer
	}{
		"no server": {},
		"sorted servers": {
			initialServers: []models.IvpnServer{
				{Country: "B", City: "A", Hostname: "A"},
				{Country: "A", City: "A", Hostname: "B"},
				{Country: "A", City: "A", Hostname: "A"},
				{Country: "A", City: "B", Hostname: "A"},
			},
			sortedServers: []models.IvpnServer{
				{Country: "A", City: "A", Hostname: "A"},
				{Country: "A", City: "A", Hostname: "B"},
				{Country: "A", City: "B", Hostname: "A"},
				{Country: "B", City: "A", Hostname: "A"},
			},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			sortServers(testCase.initialServers)
			assert.Equal(t, testCase.sortedServers, testCase.initialServers)
		})
	}
}
