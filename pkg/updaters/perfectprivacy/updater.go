package perfectprivacy

import (
	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	unzipper common.Unzipper
	warner   common.Warner
}

func New(unzipper common.Unzipper, warner common.Warner) *Updater {
	return &Updater{
		unzipper: unzipper,
		warner:   warner,
	}
}

func (u *Updater) Version() uint16 {
	return 1
}
