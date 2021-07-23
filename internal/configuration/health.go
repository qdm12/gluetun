package configuration

import (
	"fmt"
	"strings"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/params"
)

// Health contains settings for the healthcheck and health server.
type Health struct {
	ServerAddress string
	OpenVPN       HealthyWait
}

func (settings *Health) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *Health) lines() (lines []string) {
	lines = append(lines, lastIndent+"Health:")

	lines = append(lines, indent+lastIndent+"Server address: "+settings.ServerAddress)

	lines = append(lines, indent+lastIndent+"OpenVPN:")
	for _, line := range settings.OpenVPN.lines() {
		lines = append(lines, indent+indent+line)
	}

	return lines
}

// Read is to be used for the healthcheck query mode.
func (settings *Health) Read(env params.Env, logger logging.Logger) (err error) {
	reader := newReader(env, logger)
	return settings.read(reader)
}

func (settings *Health) read(r reader) (err error) {
	var warning string
	settings.ServerAddress, warning, err = r.env.ListeningAddress(
		"HEALTH_SERVER_ADDRESS", params.Default("127.0.0.1:9999"))
	if warning != "" {
		r.logger.Warn("environment variable HEALTH_SERVER_ADDRESS: " + warning)
	}
	if err != nil {
		return fmt.Errorf("environment variable HEALTH_SERVER_ADDRESS: %w", err)
	}

	settings.OpenVPN.Initial, err = r.env.Duration("HEALTH_OPENVPN_DURATION_INITIAL", params.Default("6s"))
	if err != nil {
		return fmt.Errorf("environment variable HEALTH_OPENVPN_DURATION_INITIAL: %w", err)
	}

	settings.OpenVPN.Addition, err = r.env.Duration("HEALTH_OPENVPN_DURATION_ADDITION", params.Default("5s"))
	if err != nil {
		return fmt.Errorf("environment variable HEALTH_OPENVPN_DURATION_ADDITION: %w", err)
	}

	return nil
}
