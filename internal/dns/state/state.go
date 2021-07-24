package state

import (
	"sync"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/loopstate"
)

var _ Manager = (*State)(nil)

type Manager interface {
	SettingsGetterSetter
}

func New(statusApplier loopstate.Applier,
	settings configuration.DNS,
	updateTicker chan<- struct{}) *State {
	return &State{
		statusApplier: statusApplier,
		settings:      settings,
		updateTicker:  updateTicker,
	}
}

type State struct {
	statusApplier loopstate.Applier

	settings   configuration.DNS
	settingsMu sync.RWMutex

	updateTicker chan<- struct{}
}
