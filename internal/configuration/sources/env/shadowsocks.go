package env

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (r *Reader) readShadowsocks() (shadowsocks settings.Shadowsocks, err error) {
	shadowsocks.Enabled, err = envToBoolPtr("SHADOWSOCKS")
	if err != nil {
		return shadowsocks, fmt.Errorf("environment variable SHADOWSOCKS: %w", err)
	}

	shadowsocks.Address = r.readShadowsocksAddress()
	shadowsocks.LogAddresses, err = envToBoolPtr("SHADOWSOCKS_LOG")
	if err != nil {
		return shadowsocks, fmt.Errorf("environment variable SHADOWSOCKS_LOG: %w", err)
	}
	shadowsocks.CipherName = r.readShadowsocksCipher()
	shadowsocks.Password = envToStringPtr("SHADOWSOCKS_PASSWORD")

	return shadowsocks, nil
}

func (r *Reader) readShadowsocksAddress() (address string) {
	key, value := r.getEnvWithRetro("SHADOWSOCKS_LISTENING_ADDRESS", "SHADOWSOCKS_PORT")
	if value == "" {
		return ""
	}

	if key == "SHADOWSOCKS_LISTENING_ADDRESS" {
		return value
	}

	// Retro-compatibility
	return ":" + value
}

func (r *Reader) readShadowsocksCipher() (cipher string) {
	_, cipher = r.getEnvWithRetro("SHADOWSOCKS_CIPHER", "SHADOWSOCKS_METHOD")
	return cipher
}
