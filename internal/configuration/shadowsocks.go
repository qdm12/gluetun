package configuration

import (
	"fmt"
	"strings"

	"github.com/qdm12/golibs/params"
	"github.com/qdm12/ss-server/pkg/tcpudp"
)

// ShadowSocks contains settings to configure the Shadowsocks server.
type ShadowSocks struct {
	Enabled bool
	tcpudp.Settings
}

func (settings *ShadowSocks) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *ShadowSocks) lines() (lines []string) {
	if !settings.Enabled {
		return nil
	}

	lines = append(lines, lastIndent+"Shadowsocks server:")

	lines = append(lines, indent+lastIndent+"Listening address: "+settings.Address)

	lines = append(lines, indent+lastIndent+"Cipher: "+settings.CipherName)

	if settings.LogAddresses {
		lines = append(lines, indent+lastIndent+"Log addresses: enabled")
	}

	return lines
}

func (settings *ShadowSocks) read(r reader) (err error) {
	settings.Enabled, err = r.env.OnOff("SHADOWSOCKS", params.Default("off"))
	if !settings.Enabled {
		return nil
	} else if err != nil {
		return fmt.Errorf("environment variable SHADOWSOCKS: %w", err)
	}

	settings.Password, err = r.getFromEnvOrSecretFile("SHADOWSOCKS_PASSWORD", settings.Enabled, nil)
	if err != nil {
		return err
	}

	settings.LogAddresses, err = r.env.OnOff("SHADOWSOCKS_LOG", params.Default("off"))
	if err != nil {
		return fmt.Errorf("environment variable SHADOWSOCKS_LOG: %w", err)
	}

	settings.CipherName, err = r.env.Get("SHADOWSOCKS_CIPHER", params.Default("chacha20-ietf-poly1305"),
		params.RetroKeys([]string{"SHADOWSOCKS_METHOD"}, r.onRetroActive))
	if err != nil {
		return fmt.Errorf("environment variable SHADOWSOCKS_CIPHER (or SHADOWSOCKS_METHOD): %w", err)
	}

	warning, err := settings.getAddress(r.env)
	if warning != "" {
		r.logger.Warn(warning)
	}
	if err != nil {
		return err
	}

	return nil
}

func (settings *ShadowSocks) getAddress(env params.Env) (
	warning string, err error) {
	address, addrWarning, err := env.ListeningAddress("SHADOWSOCKS_LISTENING_ADDRESS")
	if err != nil {
		return addrWarning, fmt.Errorf("environment variable SHADOWSOCKS_LISTENING_ADDRESS: %w", err)
	}

	if address != "" {
		settings.Address = address
		return addrWarning, nil
	}

	// Retro-compatibility
	const retroWarning = "You are using the old environment variable " +
		"SHADOWSOCKS_PORT, please consider using " +
		"SHADOWSOCKS_LISTENING_ADDRESS instead"
	port, _, err := env.ListeningPort("SHADOWSOCKS_PORT")
	if err != nil {
		return retroWarning, fmt.Errorf("environment variable SHADOWSOCKS_PORT: %w", err)
	} else if port > 0 {
		settings.Address = ":" + fmt.Sprint(port)
		return retroWarning, nil
	}

	// Default value
	settings.Address = ":8388"
	return addrWarning, nil
}
