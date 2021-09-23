package configuration

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/params"
)

// Health contains settings for the healthcheck and health server.
type Health struct {
	ServerAddress string
	AddressToPing string
	VPN           HealthyWait
}

func (settings *Health) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *Health) lines() (lines []string) {
	lines = append(lines, lastIndent+"Health:")

	lines = append(lines, indent+lastIndent+"Server address: "+settings.ServerAddress)

	lines = append(lines, indent+lastIndent+"Address to ping: "+settings.AddressToPing)

	lines = append(lines, indent+lastIndent+"VPN:")
	for _, line := range settings.VPN.lines() {
		lines = append(lines, indent+indent+line)
	}

	return lines
}

// Read is to be used for the healthcheck query mode.
func (settings *Health) Read(env params.Interface, warner Warner) (err error) {
	reader := newReader(env, models.AllServers{}, warner) // note: no need for servers data
	return settings.read(reader)
}

func (settings *Health) read(r reader) (err error) {
	var warning string
	settings.ServerAddress, warning, err = r.env.ListeningAddress(
		"HEALTH_SERVER_ADDRESS", params.Default("127.0.0.1:9999"))
	if warning != "" {
		r.warner.Warn("environment variable HEALTH_SERVER_ADDRESS: " + warning)
	}
	if err != nil {
		return fmt.Errorf("environment variable HEALTH_SERVER_ADDRESS: %w", err)
	}

	settings.AddressToPing, err = r.env.Get("HEALTH_ADDRESS_TO_PING", params.Default("github.com"))
	if err != nil {
		return fmt.Errorf("environment variable HEALTH_ADDRESS_TO_PING: %w", err)
	}

	retroKeyOption := params.RetroKeys([]string{"HEALTH_OPENVPN_DURATION_INITIAL"}, r.onRetroActive)
	settings.VPN.Initial, err = r.env.Duration("HEALTH_VPN_DURATION_INITIAL", params.Default("6s"), retroKeyOption)
	if err != nil {
		return fmt.Errorf("environment variable HEALTH_VPN_DURATION_INITIAL: %w", err)
	}

	retroKeyOption = params.RetroKeys([]string{"HEALTH_OPENVPN_DURATION_ADDITION"}, r.onRetroActive)
	settings.VPN.Addition, err = r.env.Duration("HEALTH_VPN_DURATION_ADDITION", params.Default("5s"), retroKeyOption)
	if err != nil {
		return fmt.Errorf("environment variable HEALTH_VPN_DURATION_ADDITION: %w", err)
	}

	return nil
}
