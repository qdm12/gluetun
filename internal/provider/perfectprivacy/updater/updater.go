package updater

import "github.com/qdm12/gluetun/internal/updater/unzip"

type Updater struct {
	unzipper unzip.Unzipper
	warner   Warner
}

type Warner interface {
	Warn(s string)
}

func New(unzipper unzip.Unzipper, warner Warner) *Updater {
	return &Updater{
		unzipper: unzipper,
		warner:   warner,
	}
}
