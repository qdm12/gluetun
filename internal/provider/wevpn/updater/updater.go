package updater

import "github.com/qdm12/gluetun/internal/provider/common"

type Updater struct {
	presolver common.ParallelResolver
	warner    common.Warner
}

func New(warner common.Warner) *Updater {
	return &Updater{
		presolver: newParallelResolver(),
		warner:    warner,
	}
}
