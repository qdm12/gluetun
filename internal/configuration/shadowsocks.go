package configuration

import (
	"strconv"
	"strings"

	"github.com/qdm12/golibs/params"
)

// ShadowSocks contains settings to configure the Shadowsocks server.
type ShadowSocks struct {
	Method   string
	Password string
	Port     uint16
	Enabled  bool
	Log      bool
}

func (settings *ShadowSocks) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *ShadowSocks) lines() (lines []string) {
	if !settings.Enabled {
		return nil
	}

	lines = append(lines, lastIndent+"Shadowsocks server:")

	lines = append(lines, indent+lastIndent+"Listening port: "+strconv.Itoa(int(settings.Port)))

	lines = append(lines, indent+lastIndent+"Method: "+settings.Method)

	if settings.Log {
		lines = append(lines, indent+lastIndent+"Logging: enabled")
	}

	return lines
}

func (settings *ShadowSocks) read(r reader) (err error) {
	settings.Enabled, err = r.env.OnOff("SHADOWSOCKS", params.Default("off"))
	if err != nil || !settings.Enabled {
		return err
	}

	settings.Password, err = r.getFromEnvOrSecretFile("SHADOWSOCKS_PASSWORD", false, nil)
	if err != nil {
		return err
	}

	settings.Log, err = r.env.OnOff("SHADOWSOCKS_LOG", params.Default("off"))
	if err != nil {
		return err
	}

	settings.Method, err = r.env.Get("SHADOWSOCKS_METHOD", params.Default("chacha20-ietf-poly1305"))
	if err != nil {
		return err
	}

	var warning string
	settings.Port, warning, err = r.env.ListeningPort("SHADOWSOCKS_PORT", params.Default("8388"))
	if len(warning) > 0 {
		r.logger.Warn(warning)
	}
	if err != nil {
		return err
	}

	return nil
}
