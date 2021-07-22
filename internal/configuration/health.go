package configuration

import (
	"strings"

	"github.com/qdm12/golibs/params"
)

// Health contains settings for the healthcheck and health server.
type Health struct {
	OpenVPN HealthyWait
}

func (settings *Health) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *Health) lines() (lines []string) {
	lines = append(lines, lastIndent+"Health:")

	lines = append(lines, indent+lastIndent+"OpenVPN:")
	for _, line := range settings.OpenVPN.lines() {
		lines = append(lines, indent+indent+line)
	}

	return lines
}

func (settings *Health) read(r reader) (err error) {
	settings.OpenVPN.Initial, err = r.env.Duration("HEALTH_OPENVPN_DURATION_INITIAL", params.Default("6s"))
	if err != nil {
		return err
	}

	settings.OpenVPN.Addition, err = r.env.Duration("HEALTH_OPENVPN_DURATION_ADDITION", params.Default("5s"))
	if err != nil {
		return err
	}

	return nil
}
