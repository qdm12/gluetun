package state

import (
	"github.com/qdm12/gluetun/internal/models"
)

func (s *State) GetData() (data models.PublicIP) {
	s.ipDataMu.RLock()
	defer s.ipDataMu.RUnlock()
	return s.ipData.Copy()
}

func (s *State) SetData(data models.PublicIP) {
	s.ipDataMu.Lock()
	defer s.ipDataMu.Unlock()
	s.ipData = data.Copy()
}
