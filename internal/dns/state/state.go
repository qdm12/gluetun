package state

import (
	"sync"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/loopstate"
)

var _ Manager = (*State)(nil)

type Manager interface {
	SettingsGetSetter
}

func New(statusApplier loopstate.Applier,
	settings settings.DNS,
	updateTicker chan<- struct{}) *State {
	return &State{
		statusApplier: statusApplier,
		settings:      settings,
		updateTicker:  updateTicker,
	}
}

type State struct {
	statusApplier loopstate.Applier

	settings   settings.DNS
	settingsMu sync.RWMutex

	updateTicker chan<- struct{}
}
