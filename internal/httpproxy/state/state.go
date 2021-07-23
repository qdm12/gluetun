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
	settings configuration.HTTPProxy) *State {
	return &State{
		statusApplier: statusApplier,
		settings:      settings,
	}
}

type State struct {
	statusApplier loopstate.Applier
	settings      configuration.HTTPProxy
	settingsMu    sync.RWMutex
}
