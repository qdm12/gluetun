package cyberghost

import "github.com/qdm12/gluetun/internal/updater/resolver"

type Updater struct {
	presolver resolver.Parallel
}

func New(presolver resolver.Parallel) *Updater {
	return &Updater{
		presolver: presolver,
	}
}
