package utils

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/models"
)

func PickRandomConnection(connections []models.Connection,
	source rand.Source) models.Connection {
	return connections[rand.New(source).Intn(len(connections))] //nolint:gosec
}
