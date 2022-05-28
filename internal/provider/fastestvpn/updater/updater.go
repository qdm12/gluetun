package fastestvpn

import (
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
)

type Updater struct {
	unzipper  unzip.Unzipper
	presolver resolver.Parallel
	warner    Warner
}

type Warner interface {
	Warn(s string)
}

func New(unzipper unzip.Unzipper, presolver resolver.Parallel,
	warner Warner) *Updater {
	return &Updater{
		unzipper:  unzipper,
		presolver: presolver,
		warner:    warner,
	}
}
