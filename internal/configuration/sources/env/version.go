package env

import (
	"fmt"
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/govalid/binary"
)

func readVersion() (version settings.Version, err error) {
	version.Enabled, err = readVersionEnabled()
	if err != nil {
		return version, err
	}

	return version, nil
}

func readVersionEnabled() (enabled *bool, err error) {
	s := os.Getenv("VERSION_INFORMATION")
	if s == "" {
		return nil, nil //nolint:nilnil
	}

	enabled = new(bool)
	*enabled, err = binary.Validate(s)
	if err != nil {
		return nil, fmt.Errorf("environment variable VERSION_INFORMATION: %w", err)
	}

	return enabled, nil

}
