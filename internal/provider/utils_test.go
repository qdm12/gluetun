package provider

import (
	"math/rand"
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_pickRandomConnection(t *testing.T) {
	t.Parallel()
	connections := []models.OpenVPNConnection{
		{Port: 1}, {Port: 2}, {Port: 3}, {Port: 4},
	}
	source := rand.NewSource(0)

	connection := pickRandomConnection(connections, source)
	assert.Equal(t, models.OpenVPNConnection{Port: 3}, connection)

	connection = pickRandomConnection(connections, source)
	assert.Equal(t, models.OpenVPNConnection{Port: 3}, connection)

	connection = pickRandomConnection(connections, source)
	assert.Equal(t, models.OpenVPNConnection{Port: 2}, connection)
}

func Test_filterByPossibilities(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		value         string
		possibilities []string
		filtered      bool
	}{
		"no possibilities": {},
		"value not in possibilities": {
			value:         "c",
			possibilities: []string{"a", "b"},
			filtered:      true,
		},
		"value in possibilities": {
			value:         "c",
			possibilities: []string{"a", "b", "c"},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			filtered := filterByPossibilities(testCase.value, testCase.possibilities)
			assert.Equal(t, testCase.filtered, filtered)
		})
	}
}
