package updater

import "github.com/qdm12/gluetun/internal/provider/common"

type Updater struct {
	parallelResolver common.ParallelResolver
	warner           common.Warner
}

func New(warner common.Warner, parallelResolver common.ParallelResolver) *Updater {
	return &Updater{
		parallelResolver: parallelResolver,
		warner:           warner,
	}
}
