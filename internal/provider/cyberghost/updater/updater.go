package updater

import (
	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	parallelResolver common.ParallelResolver
}

func New(parallelResolver common.ParallelResolver) *Updater {
	return &Updater{
		parallelResolver: parallelResolver,
	}
}
