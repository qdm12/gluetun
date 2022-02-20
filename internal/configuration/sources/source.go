package sources

import "github.com/qdm12/gluetun/internal/configuration/settings"

type Source interface {
	Read() (settings settings.Settings, err error)
	ReadHealth() (settings settings.Health, err error)
	String() string
}
