package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

var (
	ErrSystemPUIDNotValid     = errors.New("PUID is not valid")
	ErrSystemPGIDNotValid     = errors.New("PGID is not valid")
	ErrSystemTimezoneNotValid = errors.New("timezone is not valid")
)

func (r *Reader) readSystem() (system settings.System, err error) {
	system.PUID, err = r.readID("PUID", "UID")
	if err != nil {
		return system, err
	}

	system.PGID, err = r.readID("PGID", "GID")
	if err != nil {
		return system, err
	}

	system.Timezone = os.Getenv("TZ")

	return system, nil
}

var ErrSystemIDNotValid = errors.New("system ID is not valid")

func (r *Reader) readID(key, retroKey string) (
	id *uint16, err error) {
	idEnvKey, idString := r.getEnvWithRetro(key, retroKey)
	if idString == "" {
		return nil, nil //nolint:nilnil
	}

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w: %s: %s",
			idEnvKey, ErrSystemIDNotValid, idString, err)
	} else if idInt < 0 || idInt > 65535 {
		return nil, fmt.Errorf("environment variable %s: %w: %d: must be between 0 and 65535",
			idEnvKey, ErrSystemIDNotValid, idInt)
	}

	return uint16Ptr(uint16(idInt)), nil
}
