package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readShadowsocks() (shadowsocks settings.Shadowsocks, err error) {
	shadowsocks.Enabled, err = s.env.BoolPtr("SHADOWSOCKS")
	if err != nil {
		return shadowsocks, err
	}

	shadowsocks.Address = s.readShadowsocksAddress()
	shadowsocks.LogAddresses, err = s.env.BoolPtr("SHADOWSOCKS_LOG")
	if err != nil {
		return shadowsocks, err
	}
	shadowsocks.CipherName = s.readShadowsocksCipher()
	shadowsocks.Password = s.env.Get("SHADOWSOCKS_PASSWORD", env.ForceLowercase(false))

	return shadowsocks, nil
}

func (s *Source) readShadowsocksAddress() (address *string) {
	key, value := s.getEnvWithRetro("SHADOWSOCKS_LISTENING_ADDRESS",
		[]string{"SHADOWSOCKS_PORT"})
	if value == nil {
		return nil
	}

	if key == "SHADOWSOCKS_LISTENING_ADDRESS" {
		return value
	}

	// Retro-compatibility
	*value = ":" + *value
	return value
}

func (s *Source) readShadowsocksCipher() (cipher string) {
	envKey, _ := s.getEnvWithRetro("SHADOWSOCKS_CIPHER",
		[]string{"SHADOWSOCKS_METHOD"})
	return s.env.String(envKey)
}
