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

func New(statusApplier loopstate.Applier, vpn settings.VPN) *State {
	return &State{
		statusApplier: statusApplier,
		vpn:           vpn,
	}
}

type State struct {
	statusApplier loopstate.Applier

	vpn        settings.VPN
	settingsMu sync.RWMutex
}
