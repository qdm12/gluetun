package secrets

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func readShadowsocks() (settings settings.Shadowsocks, err error) {
	settings.Password, err = readSecretFileAsStringPtr(
		"SHADOWSOCKS_PASSWORD_SECRETFILE",
		"/run/secrets/shadowsocks_password",
	)
	if err != nil {
		return settings, fmt.Errorf("cannot read Shadowsocks password secret file: %w", err)
	}

	return settings, nil
}
