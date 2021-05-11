package utils

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/models"
)

func PickRandomConnection(connections []models.OpenVPNConnection,
	source rand.Source) models.OpenVPNConnection {
	return connections[rand.New(source).Intn(len(connections))] //nolint:gosec
}
