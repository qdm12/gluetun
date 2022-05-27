package wevpn

import (
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_sortServers(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		initialServers []models.Server
		sortedServers  []models.Server
	}{
		"no server": {},
		"sorted servers": {
			initialServers: []models.Server{
				{City: "A", Hostname: "A"},
				{City: "A", Hostname: "B"},
				{City: "A", Hostname: "A"},
				{City: "B", Hostname: "A"},
			},
			sortedServers: []models.Server{
				{City: "A", Hostname: "A"},
				{City: "A", Hostname: "A"},
				{City: "A", Hostname: "B"},
				{City: "B", Hostname: "A"},
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
