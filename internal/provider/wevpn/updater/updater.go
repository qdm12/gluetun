package wevpn

import "github.com/qdm12/gluetun/internal/provider/common"

type Updater struct {
	presolver common.ParallelResolver
	warner    Warner
}

type Warner interface {
	Warn(s string)
}

func New(warner Warner) *Updater {
	return &Updater{
		presolver: newParallelResolver(),
		warner:    warner,
	}
}
