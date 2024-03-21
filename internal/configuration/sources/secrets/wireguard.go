package secrets

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/configuration/sources/files"
)

func (s *Source) readWireguard() (settings settings.Wireguard, err error) {
	wireguardConf, err := s.readSecretFileAsStringPtr(
		"WIREGUARD_CONF_SECRETFILE",
		"/run/secrets/wg0.conf",
	)
	if err != nil {
		return settings, fmt.Errorf("reading Wireguard conf secret file: %w", err)
	} else if wireguardConf != nil {
		return files.ParseWireguardConf([]byte(*wireguardConf))
	}
	return settings, nil
}
