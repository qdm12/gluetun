package updater

import (
	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	parallelResolver common.ParallelResolver
	warner           common.Warner
}

func New(parallelResolver common.ParallelResolver, warner common.Warner) *Updater {
	return &Updater{
		parallelResolver: parallelResolver,
		warner:           warner,
	}
}
