package fastestvpn

import (
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/updater/unzip"
)

type Updater struct {
	unzipper  unzip.Unzipper
	presolver common.ParallelResolver
	warner    Warner
}

type Warner interface {
	Warn(s string)
}

func New(unzipper unzip.Unzipper, warner Warner) *Updater {
	return &Updater{
		unzipper:  unzipper,
		presolver: newParallelResolver(),
		warner:    warner,
	}
}
