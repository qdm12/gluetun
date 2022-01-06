package state

import (
	"sync"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/loopstate"
	"github.com/qdm12/gluetun/internal/models"
)

var _ Manager = (*State)(nil)

type Manager interface {
	SettingsGetSetter
	ServersGetterSetter
	GetSettingsAndServers() (vpn settings.VPN, allServers models.AllServers)
}

func New(statusApplier loopstate.Applier,
	vpn settings.VPN, allServers models.AllServers) *State {
	return &State{
		statusApplier: statusApplier,
		vpn:           vpn,
		allServers:    allServers,
	}
}

type State struct {
	statusApplier loopstate.Applier

	vpn        settings.VPN
	settingsMu sync.RWMutex

	allServers   models.AllServers
	allServersMu sync.RWMutex
}

func (s *State) GetSettingsAndServers() (vpn settings.VPN,
	allServers models.AllServers) {
	s.settingsMu.RLock()
	s.allServersMu.RLock()
	vpn = s.vpn
	allServers = s.allServers
	s.settingsMu.RUnlock()
	s.allServersMu.RUnlock()
	return vpn, allServers
}
