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
	lines []string, connection models.Connection, intf string, err error) {
	lines, err = readCustomConfigLines(settings.Config)
	if err != nil {
		return nil, connection, "", fmt.Errorf("%w: %s", ErrReadCustomConfig, err)
	}

	connection, intf, err = extractDataFromLines(lines)
	if err != nil {
		return nil, connection, "", fmt.Errorf("%w: %s", ErrExtractConnection, err)
	}

	if intf == "" {
		intf = settings.Interface
	}

	lines = modifyCustomConfig(lines, settings, connection, intf)

	return lines, connection, intf, nil
}
