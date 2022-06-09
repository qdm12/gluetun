package ipvanish

import (
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
)

type Updater struct {
	unzipper  unzip.Unzipper
	warner    Warner
	presolver resolver.Parallel
}

type Warner interface {
	Warn(s string)
}

func New(unzipper unzip.Unzipper, warner Warner) *Updater {
	return &Updater{
		unzipper:  unzipper,
		warner:    warner,
		presolver: newParallelResolver(),
	}
}
