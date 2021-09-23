package configuration

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/qdm12/golibs/params"
)

// HTTPProxy contains settings to configure the HTTP proxy.
type HTTPProxy struct {
	User     string
	Password string
	Port     uint16
	Enabled  bool
	Stealth  bool
	Log      bool
}

func (settings *HTTPProxy) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *HTTPProxy) lines() (lines []string) {
	if !settings.Enabled {
		return nil
	}

	lines = append(lines, lastIndent+"HTTP proxy:")

	lines = append(lines, indent+lastIndent+"Port: "+strconv.Itoa(int(settings.Port)))

	if settings.User != "" {
		lines = append(lines, indent+lastIndent+"Authentication: enabled")
	}

	if settings.Log {
		lines = append(lines, indent+lastIndent+"Log: enabled")
	}

	if settings.Stealth {
		lines = append(lines, indent+lastIndent+"Stealth: enabled")
	}

	return lines
}

func (settings *HTTPProxy) read(r reader) (err error) {
	settings.Enabled, err = r.env.OnOff("HTTPPROXY", params.Default("off"),
		params.RetroKeys([]string{"TINYPROXY", "PROXY"}, r.onRetroActive))
	if err != nil {
		return fmt.Errorf("environment variable HTTPPROXY (or TINYPROXY, PROXY): %w", err)
	}

	settings.User, err = r.getFromEnvOrSecretFile("HTTPPROXY_USER", false, // compulsory
		[]string{"TINYPROXY_USER", "PROXY_USER"})
	if err != nil {
		return fmt.Errorf("environment variable HTTPPROXY_USER (or TINYPROXY_USER, PROXY_USER): %w", err)
	}

	settings.Password, err = r.getFromEnvOrSecretFile("HTTPPROXY_PASSWORD", false,
		[]string{"TINYPROXY_PASSWORD", "PROXY_PASSWORD"})
	if err != nil {
		return fmt.Errorf("environment variable HTTPPROXY_PASSWORD (or TINYPROXY_PASSWORD, PROXY_PASSWORD): %w", err)
	}

	settings.Stealth, err = r.env.OnOff("HTTPPROXY_STEALTH", params.Default("off"))
	if err != nil {
		return fmt.Errorf("environment variable HTTPPROXY_STEALTH: %w", err)
	}

	if err := settings.readLog(r); err != nil {
		return err
	}

	var warning string
	settings.Port, warning, err = r.env.ListeningPort("HTTPPROXY_PORT", params.Default("8888"),
		params.RetroKeys([]string{"TINYPROXY_PORT", "PROXY_PORT"}, r.onRetroActive))
	if len(warning) > 0 {
		r.warner.Warn(warning)
	}
	if err != nil {
		return fmt.Errorf("environment variable HTTPPROXY_PORT (or TINYPROXY_PORT, PROXY_PORT): %w", err)
	}

	return nil
}

func (settings *HTTPProxy) readLog(r reader) error {
	s, err := r.env.Get("HTTPPROXY_LOG",
		params.RetroKeys([]string{"PROXY_LOG_LEVEL", "TINYPROXY_LOG"}, r.onRetroActive))
	if err != nil {
		return fmt.Errorf("environment variable HTTPPROXY_LOG (or TINYPROXY_LOG, PROXY_LOG_LEVEL): %w", err)
	}

	switch strings.ToLower(s) {
	case "on":
		settings.Log = true
	// Retro compatibility
	case "info", "connect", "notice":
		settings.Log = true
	}

	return nil
}
