package wevpn

import (
	"github.com/qdm12/gluetun/internal/updater/resolver"
)

type Updater struct {
	presolver resolver.Parallel
	warner    Warner
}

type Warner interface {
	Warn(s string)
}

func New(presolver resolver.Parallel, warner Warner) *Updater {
	return &Updater{
		presolver: presolver,
		warner:    warner,
	}
}
