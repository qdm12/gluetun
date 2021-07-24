package state

import (
	"net"
	"sync"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/loopstate"
)

var _ Manager = (*State)(nil)

type Manager interface {
	SettingsGetSetter
	PublicIPGetSetter
}

func New(statusApplier loopstate.Applier,
	settings configuration.PublicIP,
	updateTicker chan<- struct{}) *State {
	return &State{
		statusApplier: statusApplier,
		settings:      settings,
		updateTicker:  updateTicker,
	}
}

type State struct {
	statusApplier loopstate.Applier

	settings   configuration.PublicIP
	settingsMu sync.RWMutex

	publicIP   net.IP
	publicIPMu sync.RWMutex

	updateTicker chan<- struct{}
}
