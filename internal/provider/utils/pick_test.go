package utils

import (
	"math/rand"
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_pickRandomConnection(t *testing.T) {
	t.Parallel()
	connections := []models.Connection{
		{Port: 1}, {Port: 2}, {Port: 3}, {Port: 4},
	}
	source := rand.NewSource(0)

	connection := pickRandomConnection(connections, source)
	assert.Equal(t, models.Connection{Port: 3}, connection)

	connection = pickRandomConnection(connections, source)
	assert.Equal(t, models.Connection{Port: 3}, connection)

	connection = pickRandomConnection(connections, source)
	assert.Equal(t, models.Connection{Port: 2}, connection)
}
