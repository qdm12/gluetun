package utils

import (
	"math/rand"
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_PickRandomConnection(t *testing.T) {
	t.Parallel()
	connections := []models.OpenVPNConnection{
		{Port: 1}, {Port: 2}, {Port: 3}, {Port: 4},
	}
	source := rand.NewSource(0)

	connection := PickRandomConnection(connections, source)
	assert.Equal(t, models.OpenVPNConnection{Port: 3}, connection)

	connection = PickRandomConnection(connections, source)
	assert.Equal(t, models.OpenVPNConnection{Port: 3}, connection)

	connection = PickRandomConnection(connections, source)
	assert.Equal(t, models.OpenVPNConnection{Port: 2}, connection)
}
