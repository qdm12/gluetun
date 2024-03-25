package secrets

import (
	"os"
	"path/filepath"

	"github.com/qdm12/gluetun/internal/configuration/sources/files"
)

func (s *Source) lazyLoadWireguardConf() files.WireguardConfig {
	if s.cached.wireguardLoaded {
		return s.cached.wireguardConf
	}

	path := os.Getenv("WIREGUARD_CONF_SECRETFILE")
	if path == "" {
		path = filepath.Join(s.rootDirectory, "wg0.conf")
	}

	s.cached.wireguardLoaded = true
	var err error
	s.cached.wireguardConf, err = files.ParseWireguardConf(path)
	if err != nil {
		s.warner.Warnf("skipping Wireguard config: %s", err)
	}
	return s.cached.wireguardConf
}
