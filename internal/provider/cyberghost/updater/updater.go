package cyberghost

import (
	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	presolver common.ParallelResolver
}

func New() *Updater {
	return &Updater{
		presolver: newParallelResolver(),
	}
}
