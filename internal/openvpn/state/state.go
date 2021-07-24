package state

import (
	"sync"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/loopstate"
	"github.com/qdm12/gluetun/internal/models"
)

var _ Manager = (*State)(nil)

type Manager interface {
	SettingsGetSetter
	ServersGetterSetter
	PortForwardedGetterSetter
	GetSettingsAndServers() (settings configuration.OpenVPN,
		allServers models.AllServers)
}

func New(statusApplier loopstate.Applier,
	settings configuration.OpenVPN,
	allServers models.AllServers) *State {
	return &State{
		statusApplier: statusApplier,
		settings:      settings,
		allServers:    allServers,
	}
}

type State struct {
	statusApplier loopstate.Applier

	settings   configuration.OpenVPN
	settingsMu sync.RWMutex

	allServers   models.AllServers
	allServersMu sync.RWMutex

	portForwarded   uint16
	portForwardedMu sync.RWMutex
}

func (s *State) GetSettingsAndServers() (settings configuration.OpenVPN,
	allServers models.AllServers) {
	s.settingsMu.RLock()
	s.allServersMu.RLock()
	settings = s.settings
	allServers = s.allServers
	s.settingsMu.RUnlock()
	s.allServersMu.RUnlock()
	return settings, allServers
}
