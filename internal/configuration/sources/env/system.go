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
	id *uint32, err error) {
	idEnvKey, idString := r.getEnvWithRetro(key, retroKey)
	if idString == "" {
		return nil, nil //nolint:nilnil
	}

	const base = 10
	const bitSize = 64
	const max = uint64(^uint32(0))
	idUint64, err := strconv.ParseUint(idString, base, bitSize)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w: %s",
			idEnvKey, ErrSystemIDNotValid, err)
	} else if idUint64 > max {
		return nil, fmt.Errorf("environment variable %s: %w: %d: must be between 0 and %d",
			idEnvKey, ErrSystemIDNotValid, idUint64, max)
	}

	return uint32Ptr(uint32(idUint64)), nil
}
