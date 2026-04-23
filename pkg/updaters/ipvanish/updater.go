package ipvanish

import (
	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	unzipper         common.Unzipper
	warner           common.Warner
	parallelResolver common.ParallelResolver
}

func New(unzipper common.Unzipper, warner common.Warner,
	parallelResolver common.ParallelResolver,
) *Updater {
	return &Updater{
		unzipper:         unzipper,
		warner:           warner,
		parallelResolver: parallelResolver,
	}
}

func (u *Updater) Version() uint16 {
	return 2 //nolint:mnd
}
