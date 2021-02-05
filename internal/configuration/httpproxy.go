package configuration

import (
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

func (h *HTTPProxy) String() string {
	return strings.Join(h.lines(), "\n")
}

func (h *HTTPProxy) lines() (lines []string) {
	if !h.Enabled {
		return nil
	}

	lines = append(lines, lastIndent+"HTTP proxy:")

	lines = append(lines, indent+lastIndent+"Port: "+strconv.Itoa(int(h.Port)))

	if h.User != "" {
		lines = append(lines, indent+lastIndent+"Authentication: enabled")
	}

	if h.Log {
		lines = append(lines, indent+lastIndent+"Log: enabled")
	}

	if h.Stealth {
		lines = append(lines, indent+lastIndent+"Stealth: enabled")
	}

	return lines
}

func (settings *HTTPProxy) read(r reader) (err error) {
	settings.Enabled, err = r.env.OnOff("HTTPPROXY", params.Default("off"),
		params.RetroKeys([]string{"TINYPROXY", "PROXY"}, r.onRetroActive))
	if err != nil {
		return err
	}

	settings.User, err = r.getFromEnvOrSecretFile("HTTPPROXY_USER", false, // compulsory
		[]string{"TINYPROXY_USER", "PROXY_USER"})
	if err != nil {
		return err
	}

	settings.Password, err = r.getFromEnvOrSecretFile("HTTPPROXY_USER", false,
		[]string{"TINYPROXY_PASSWORD", "PROXY_PASSWORD"})
	if err != nil {
		return err
	}

	settings.Stealth, err = r.env.OnOff("HTTPPROXY_STEALTH", params.Default("off"))
	if err != nil {
		return err
	}

	if err := settings.readLog(r); err != nil {
		return err
	}

	var warning string
	settings.Port, warning, err = r.env.ListeningPort("HTTPPROXY_PORT", params.Default("8888"),
		params.RetroKeys([]string{"TINYPROXY_PORT", "PROXY_PORT"}, r.onRetroActive))
	if len(warning) > 0 {
		r.logger.Warn(warning)
	}
	if err != nil {
		return err
	}

	return nil
}

func (settings *HTTPProxy) readLog(r reader) error {
	s, err := r.env.Get("HTTPPROXY_LOG",
		params.RetroKeys([]string{"PROXY_LOG_LEVEL", "TINYPROXY_LOG"}, r.onRetroActive))
	if err != nil {
		return err
	}

	switch strings.ToLower(s) {
	case "on":
		settings.Enabled = true
	// Retro compatibility
	case "info", "connect", "notice":
		settings.Enabled = true
	}

	return nil
}
