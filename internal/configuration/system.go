package configuration

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/qdm12/golibs/params"
)

// System contains settings to configure system related elements.
type System struct {
	PUID     int
	PGID     int
	Timezone string
}

func (settings *System) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *System) lines() (lines []string) {
	lines = append(lines, lastIndent+"System:")
	lines = append(lines, indent+lastIndent+"Process user ID: "+strconv.Itoa(settings.PUID))
	lines = append(lines, indent+lastIndent+"Process group ID: "+strconv.Itoa(settings.PGID))

	if len(settings.Timezone) > 0 {
		lines = append(lines, indent+lastIndent+"Timezone: "+settings.Timezone)
	} else {
		lines = append(lines, indent+lastIndent+"Timezone: NOT SET ⚠️ - it can cause time related issues")
	}
	return lines
}

func (settings *System) read(r reader) (err error) {
	const maxID = 65535
	settings.PUID, err = r.env.IntRange("PUID", 0, maxID, params.Default("1000"),
		params.RetroKeys([]string{"UID"}, r.onRetroActive))
	if err != nil {
		return fmt.Errorf("environment variable PUID (or UID): %w", err)
	}

	settings.PGID, err = r.env.IntRange("PGID", 0, maxID, params.Default("1000"),
		params.RetroKeys([]string{"GID"}, r.onRetroActive))
	if err != nil {
		return fmt.Errorf("environment variable PGID (or GID): %w", err)
	}

	settings.Timezone, err = r.env.Get("TZ")
	if err != nil {
		return fmt.Errorf("environment variable TZ: %w", err)
	}

	return nil
}
