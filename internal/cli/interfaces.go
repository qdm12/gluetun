package cli

import "github.com/qdm12/gluetun/internal/configuration/settings"

type Source interface {
	Read() (settings settings.Settings, err error)
	ReadHealth() (health settings.Health, err error)
	String() string
}
