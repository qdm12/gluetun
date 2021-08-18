package custom

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/models"
)

var (
	ErrReadCustomConfig  = errors.New("cannot read custom configuration file")
	ErrExtractConnection = errors.New("cannot extract connection from custom configuration file")
)

func BuildConfig(settings configuration.OpenVPN) (
	lines []string, connection models.OpenVPNConnection, err error) {
	lines, err = readCustomConfigLines(settings.Config)
	if err != nil {
		return nil, connection, fmt.Errorf("%w: %s", ErrReadCustomConfig, err)
	}

	connection, err = extractConnectionFromLines(lines)
	if err != nil {
		return nil, connection, fmt.Errorf("%w: %s", ErrExtractConnection, err)
	}

	lines = modifyCustomConfig(lines, settings, connection)

	return lines, connection, nil
}
