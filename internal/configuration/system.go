package configuration

import (
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

func (s *System) String() string {
	return strings.Join(s.lines(), "\n")
}

func (s *System) lines() (lines []string) {
	lines = append(lines, lastIndent+"System:")
	lines = append(lines, indent+lastIndent+"Process user ID: "+strconv.Itoa(s.PUID))
	lines = append(lines, indent+lastIndent+"Process group ID: "+strconv.Itoa(s.PGID))

	if len(s.Timezone) > 0 {
		lines = append(lines, indent+lastIndent+"Timezone: "+s.Timezone)
	} else {
		lines = append(lines, indent+lastIndent+"Timezone: NOT SET ⚠️ CAN CAUSE TIME ISSUES")
	}
	return lines
}

func (settings *System) read(r reader) (err error) {
	settings.PUID, err = r.env.IntRange("PUID", 0, 65535, params.Default("1000"),
		params.RetroKeys([]string{"UID"}, r.onRetroActive))
	if err != nil {
		return err
	}

	settings.PGID, err = r.env.IntRange("PGID", 0, 65535, params.Default("1000"),
		params.RetroKeys([]string{"GID"}, r.onRetroActive))
	if err != nil {
		return err
	}

	settings.Timezone, err = r.env.Get("TZ")
	if err != nil {
		return err
	}

	return nil
}
