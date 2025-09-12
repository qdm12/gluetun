package cli

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

type Source interface {
	Read() (settings settings.Settings, err error)
	ReadHealth() (health settings.Health, err error)
	String() string
}

type SubCommand interface {
	Run(ctx context.Context) (err error)
	Name() string
	Description() string
}
