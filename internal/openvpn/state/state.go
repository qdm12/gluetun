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
	GetSettingsAndServers() (vpn configuration.VPN,
		provider configuration.Provider, allServers models.AllServers)
}

func New(statusApplier loopstate.Applier,
	vpn configuration.VPN, provider configuration.Provider,
	allServers models.AllServers) *State {
	return &State{
		statusApplier: statusApplier,
		vpn:           vpn,
		provider:      provider,
		allServers:    allServers,
	}
}

type State struct {
	statusApplier loopstate.Applier

	vpn        configuration.VPN
	provider   configuration.Provider
	settingsMu sync.RWMutex

	allServers   models.AllServers
	allServersMu sync.RWMutex
}

func (s *State) GetSettingsAndServers() (vpn configuration.VPN,
	provider configuration.Provider, allServers models.AllServers) {
	s.settingsMu.RLock()
	s.allServersMu.RLock()
	vpn = s.vpn
	provider = s.provider
	allServers = s.allServers
	s.settingsMu.RUnlock()
	s.allServersMu.RUnlock()
	return vpn, provider, allServers
}
