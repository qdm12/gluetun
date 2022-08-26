package secrets

import "github.com/qdm12/gluetun/internal/configuration/settings"

func (s *Source) ReadHealth() (settings settings.Health, err error) { return settings, nil }
