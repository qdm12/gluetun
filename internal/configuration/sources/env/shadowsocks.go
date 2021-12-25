package env

import (
	"fmt"
	"os"

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
	address = os.Getenv("SHADOWSOCKS_LISTENING_ADDRESS")
	if address != "" {
		return address
	}

	// Retro-compatibility
	portString := os.Getenv("SHADOWSOCKS_PORT")
	if portString != "" {
		r.onRetroActive("SHADOWSOCKS_PORT", "SHADOWSOCKS_LISTENING_ADDRESS")
		return ":" + portString
	}

	return ""
}

func (r *Reader) readShadowsocksCipher() (cipher string) {
	cipher = os.Getenv("SHADOWSOCKS_CIPHER")
	if cipher != "" {
		return cipher
	}
	// Retro-compatibility
	cipher = os.Getenv("SHADOWSOCKS_METHOD")
	if cipher != "" {
		r.onRetroActive("SHADOWSOCKS_METHOD", "SHADOWSOCKS_CIPHER")
	}
	return cipher
}
