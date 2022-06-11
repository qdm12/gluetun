package updater

import (
	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	unzipper  common.Unzipper
	presolver common.ParallelResolver
	warner    common.Warner
}

func New(unzipper common.Unzipper, warner common.Warner) *Updater {
	return &Updater{
		unzipper:  unzipper,
		presolver: newParallelResolver(),
		warner:    warner,
	}
}
