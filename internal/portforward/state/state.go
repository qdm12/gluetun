package state

import (
	"sync"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/loopstate"
)

var _ Manager = (*State)(nil)

type Manager interface {
	SettingsGetSetter
	PortForwardedGetterSetter
	StartDataGetterSetter
}

func New(statusApplier loopstate.Applier,
	settings settings.PortForwarding) *State {
	return &State{
		statusApplier: statusApplier,
		settings:      settings,
	}
}

type State struct {
	statusApplier loopstate.Applier

	settings   settings.PortForwarding
	settingsMu sync.RWMutex

	portForwarded   uint16
	portForwardedMu sync.RWMutex

	startData   StartData
	startDataMu sync.RWMutex
}
