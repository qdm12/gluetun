package expressvpn

import (
	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	unzipper         common.Unzipper
	parallelResolver common.ParallelResolver
	warner           common.Warner
}

func New(unzipper common.Unzipper, warner common.Warner,
	parallelResolver common.ParallelResolver,
) *Updater {
	return &Updater{
		unzipper:         unzipper,
		parallelResolver: parallelResolver,
		warner:           warner,
	}
}

func (u *Updater) Version() uint16 {
	return 2 //nolint:mnd
}
