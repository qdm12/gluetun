package configuration

import (
	"strings"
	"time"
)

type HealthyWait struct {
	// Initial is the initial duration to wait for the program
	// to be healthy before taking action.
	Initial time.Duration
	// Addition is the duration to add to the Initial duration
	// after Initial has expired to wait longer for the program
	// to be healthy.
	Addition time.Duration
}

func (settings *HealthyWait) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *HealthyWait) lines() (lines []string) {
	lines = append(lines, lastIndent+"Initial duration: "+settings.Initial.String())

	if settings.Addition > 0 {
		lines = append(lines, lastIndent+"Addition duration: "+settings.Addition.String())
	}

	return lines
}
