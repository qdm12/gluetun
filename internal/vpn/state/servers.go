package state

import "github.com/qdm12/gluetun/internal/models"

type ServersGetterSetter interface {
	GetServers() (servers models.AllServers)
	SetServers(servers models.AllServers)
}

func (s *State) GetServers() (servers models.AllServers) {
	s.allServersMu.RLock()
	defer s.allServersMu.RUnlock()
	return s.allServers
}

func (s *State) SetServers(servers models.AllServers) {
	s.allServersMu.Lock()
	defer s.allServersMu.Unlock()
	s.allServers = servers
}
