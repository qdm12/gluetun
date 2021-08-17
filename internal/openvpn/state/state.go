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
	GetSettingsAndServers() (openvpn configuration.OpenVPN,
		provider configuration.Provider, allServers models.AllServers)
}

func New(statusApplier loopstate.Applier,
	openvpn configuration.OpenVPN, provider configuration.Provider,
	allServers models.AllServers) *State {
	return &State{
		statusApplier: statusApplier,
		openvpn:       openvpn,
		provider:      provider,
		allServers:    allServers,
	}
}

type State struct {
	statusApplier loopstate.Applier

	openvpn    configuration.OpenVPN
	provider   configuration.Provider
	settingsMu sync.RWMutex

	allServers   models.AllServers
	allServersMu sync.RWMutex
}

func (s *State) GetSettingsAndServers() (openvpn configuration.OpenVPN,
	provider configuration.Provider, allServers models.AllServers) {
	s.settingsMu.RLock()
	s.allServersMu.RLock()
	openvpn = s.openvpn
	provider = s.provider
	allServers = s.allServers
	s.settingsMu.RUnlock()
	s.allServersMu.RUnlock()
	return openvpn, provider, allServers
}
