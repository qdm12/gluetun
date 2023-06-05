package env

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readShadowsocks() (shadowsocks settings.Shadowsocks, err error) {
	shadowsocks.Enabled, err = s.env.BoolPtr("SHADOWSOCKS")
	if err != nil {
		return shadowsocks, err
	}

	shadowsocks.Address, err = s.readShadowsocksAddress()
	if err != nil {
		return shadowsocks, err
	}
	shadowsocks.LogAddresses, err = s.env.BoolPtr("SHADOWSOCKS_LOG")
	if err != nil {
		return shadowsocks, err
	}
	shadowsocks.CipherName = s.env.String("SHADOWSOCKS_CIPHER",
		env.RetroKeys("SHADOWSOCKS_METHOD"))
	shadowsocks.Password = s.env.Get("SHADOWSOCKS_PASSWORD", env.ForceLowercase(false))

	return shadowsocks, nil
}

func (s *Source) readShadowsocksAddress() (address *string, err error) {
	const currentKey = "SHADOWSOCKS_LISTENING_ADDRESS"
	port, err := s.env.Uint16Ptr("SHADOWSOCKS_PORT") // retro-compatibility
	if err != nil {
		return nil, err
	} else if port != nil {
		s.handleDeprecatedKey("SHADOWSOCKS_PORT", currentKey)
		return ptrTo(fmt.Sprintf(":%d", *port)), nil
	}

	return s.env.Get(currentKey), nil
}
