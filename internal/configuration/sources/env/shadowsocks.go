package env

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readShadowsocks() (shadowsocks settings.Shadowsocks, err error) {
	shadowsocks.Enabled, err = envToBoolPtr("SHADOWSOCKS")
	if err != nil {
		return shadowsocks, fmt.Errorf("environment variable SHADOWSOCKS: %w", err)
	}

	shadowsocks.Address = s.readShadowsocksAddress()
	shadowsocks.LogAddresses, err = envToBoolPtr("SHADOWSOCKS_LOG")
	if err != nil {
		return shadowsocks, fmt.Errorf("environment variable SHADOWSOCKS_LOG: %w", err)
	}
	shadowsocks.CipherName = s.readShadowsocksCipher()
	shadowsocks.Password = envToStringPtr("SHADOWSOCKS_PASSWORD")

	return shadowsocks, nil
}

func (s *Source) readShadowsocksAddress() (address string) {
	key, value := s.getEnvWithRetro("SHADOWSOCKS_LISTENING_ADDRESS", "SHADOWSOCKS_PORT")
	if value == "" {
		return ""
	}

	if key == "SHADOWSOCKS_LISTENING_ADDRESS" {
		return value
	}

	// Retro-compatibility
	return ":" + value
}

func (s *Source) readShadowsocksCipher() (cipher string) {
	_, cipher = s.getEnvWithRetro("SHADOWSOCKS_CIPHER", "SHADOWSOCKS_METHOD")
	return strings.ToLower(cipher)
}
