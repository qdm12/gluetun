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
	settings settings.HTTPProxy) *State {
	return &State{
		statusApplier: statusApplier,
		settings:      settings,
	}
}

type State struct {
	statusApplier loopstate.Applier
	settings      settings.HTTPProxy
	settingsMu    sync.RWMutex
}
