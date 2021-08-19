package utils

import (
	"errors"
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

var ErrTargetIPNotFound = errors.New("target IP address not found")

func GetTargetIPConnection(connections []models.Connection,
	targetIP net.IP) (connection models.Connection, err error) {
	for _, connection := range connections {
		if targetIP.Equal(connection.IP) {
			return connection, nil
		}
	}
	return connection, fmt.Errorf("%w: in %d filtered connections",
		ErrTargetIPNotFound, len(connections))
}
