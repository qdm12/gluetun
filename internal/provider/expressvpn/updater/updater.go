package updater

import (
	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	unzipper  common.Unzipper
	presolver common.ParallelResolver
	warner    common.Warner
}

func New(unzipper common.Unzipper, warner common.Warner,
	parallelResolver common.ParallelResolver) *Updater {
	return &Updater{
		unzipper:  unzipper,
		presolver: parallelResolver,
		warner:    warner,
	}
}
